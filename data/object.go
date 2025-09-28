package data

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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

// converts bytes data into sha1 hashed string
func IssueObjID(data []byte) (oid string) {
	s := sha1.Sum(data)
	oid = hex.EncodeToString(s[:])
	return oid
}

// converts content in byte into a blob object and
// save it as a file with an oid in the target directory
// e.g. { targetDir: .pgit/objects, savedfile: .pgit/objects/{oid} }
func SaveBlobObj(targetDir string, content []byte) (oid string, err error) {
	obj := NewObject(ObjTypeBlob, IdentBlob, content)
	oid = IssueObjID(obj.Encode())
	if err := WriteFile(targetDir, content); err != nil {
		return "", fmt.Errorf("SaveHashObj: %w", err)
	}
	return oid, nil
}
