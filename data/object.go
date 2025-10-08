package data

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	//identifier for object type
	IdentBlob = []byte{'b', 0x00}
)

const (
	ObjTypeBlob = "blob"
	ObjTypeTree = "tree"
)

type Object struct {
	objType string
	ident   []byte //identifier
	data    []byte
}

func NewObject(objType string, ident []byte, data []byte) *Object {
	return &Object{
		objType: objType,
		ident:   ident,
		data:    data,
	}
}

func (o *Object) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteString(o.objType)
	buf.Write(o.ident)
	buf.Write(o.data)
	return buf.Bytes()
}

func (o *Object) Type() string {
	return o.objType
}
func (o *Object) Ident() []byte {
	return o.ident
}
func (o *Object) Data() []byte {
	return o.data
}

// converts bytes data into sha1-hashed string
func IssueObjID(data []byte) (oid string) {
	s := sha1.Sum(data)
	oid = hex.EncodeToString(s[:])
	return oid
}

// converts content in byte into a blob object under the hood, and
// save it as a file with an oid in the dirPath
// e.g. { dirPath: .pgit/objects, savedfile: .pgit/objects/{oid} }
func SaveBlobObj(dirPath string, content []byte) (oid string, err error) {
	obj := NewObject(ObjTypeBlob, IdentBlob, content)
	oid = IssueObjID(obj.Encode())
	if err := WriteFile(filepath.Join(dirPath, oid), content); err != nil {
		return "", fmt.Errorf("SaveHashObj: %w", err)
	}
	return oid, nil
}

// "Tree object" represents a directory in the whole package and the real stuff is just a file.
// WriteTree walks through the srcDirPath and do the following things for each file (or directory):
// ・convert each file to a hashed-object, save its oid in the trgDirPath, and record it in a new file (= tree)
// ・if the given file is a directory, then recursively do the same
// ・at the end, save the whole directory (i.e. srcDir) as a hashed-object (= tree object) in the trgDirPath
func WriteTree(srcDirPath string, trgDirPath string) (treeOid string, err error) {
	if !filepath.IsAbs(trgDirPath) {
		trgDirPath, err = filepath.Abs(trgDirPath)
		if err != nil {
			return "", fmt.Errorf("WriteTree: %s", err)
		}
	}
	if err != nil {
		return "", fmt.Errorf("WriteTree: %w", err)
	}
	var buf bytes.Buffer
	err = filepath.WalkDir(srcDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == srcDirPath {
			return nil
		}
		if isExcluded(d.Name()) {
			return filepath.SkipDir
		}
		if d.IsDir() {
			oid, err := WriteTree(path, trgDirPath)
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintf(&buf, "%s %s %s\n", ObjTypeTree, oid, d.Name()); err != nil {
				return err
			}
			return nil
		}
		b, err := ReadAllFileContent(path)
		if err != nil {
			return fmt.Errorf("writeTree: %w", err)
		}
		oid, err := SaveBlobObj(trgDirPath, b)
		if err != nil {
			return fmt.Errorf("writeTree: %w", err)
		}
		fmt.Fprintf(&buf, "%s %s %s\n", ObjTypeBlob, oid, d.Name())
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("writeTree: %w", err)
	}
	treeOid, err = SaveBlobObj(trgDirPath, buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("writeTree: %w", err)
	}
	return treeOid, nil
}

func isExcluded(baseName string) bool {
	return baseName == PgitDirBase
}

// ReadTree reads the content of a file (= srcDirPath/{treeOid}) and
// lays out all the files and directories in the target directory.
func ReadTree(treeOid string, srcDirPath string, trgDirPath string) error {
	treeFilePath := filepath.Join(srcDirPath, treeOid)
	treeContent, err := ReadAllFileContent(treeFilePath)
	if err != nil {
		return fmt.Errorf("ReadTree: %w", err)
	}
	sc := bufio.NewScanner(bytes.NewReader(treeContent))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		line := sc.Bytes()
		sep := bytes.Split(line, []byte{' '}) //separated bytes hold "objType, oid, filename"
		if len(sep) < 3 {
			return fmt.Errorf("ReadTree: invalid data: { object: %s }", sep)
		}
		_, oid, filename := sep[0], sep[1], sep[2]
		fc, err := ReadAllFileContent(filepath.Join(srcDirPath, string(oid)))
		if err != nil {
			return fmt.Errorf("ReadTree: %w", err)
		}
		f, err := os.Create(filepath.Join(trgDirPath, string(filename)))
		if err != nil {
			return err
		}
		f.Write(fc)
		f.Close()
	}
	return nil
}

type Commit struct {
	TreeOid string
	Parent  string
	Msg     string
}

// Read a content of a file (= .pgit/objects/{oid}), and convert it to Commit struct.
func GetCommit(oid string) (*Commit, error) {
	c := &Commit{}
	b, err := ReadAllFileContent(filepath.Join(PgitDirBase, ObjDirBase, oid))
	if err != nil {
		return nil, fmt.Errorf("GetCommit: %w", err)
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		line := sc.Text()
		sep := strings.Split(line, " ")
		switch sep[0] {
		case "tree":
			c.TreeOid = sep[1]
		case "parent":
			c.Parent = sep[1]
		default:
			c.Msg = strings.Join(sep, " ")
		}
	}
	return c, nil
}

type TreeElem struct {
	ObjType string //blob or tree
	Oid     string
	Name    string //filefname
	Child   Tree
}

// { key: filename, value: TreeElem }
type Tree map[string]*TreeElem

// Parse tree files existing in the path specified, and convert them into type Tree
func ParseTreeFile(path string) (Tree, error) {
	c, err := ReadAllFileContent(path)
	if err != nil {
		return nil, fmt.Errorf("ParseTree: %w", err)
	}
	tree := make(Tree)
	sc := bufio.NewScanner(bytes.NewReader(c))
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		line := sc.Text()
		sep := strings.Split(line, " ")
		if len(sep) < 3 {
			continue
		}
		objType, oid, name := sep[0], sep[1], sep[2]
		elm := &TreeElem{
			ObjType: objType,
			Oid:     oid,
			Name:    name,
			Child:   nil,
		}
		if objType == ObjTypeTree {
			elm.Child, err = ParseTreeFile(filepath.Join(filepath.Dir(path), name))
			if err != nil {
				return nil, fmt.Errorf("ParseTree: %w", err)
			}
			tree[name] = elm
			continue
		}
		tree[name] = elm
	}
	return tree, nil
}
