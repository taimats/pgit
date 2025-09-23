package cmd

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	//identifier for object type
	IdentBlob = []byte{'b', 0x00}
)

const (
	ObjTypeBlob = "blob"
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
func (o *Object) Indent() []byte {
	return o.ident
}
func (o *Object) Data() []byte {
	return o.data
}

var hashObjCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "save a hashed-object",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := CheckPgitInit()
		if err != nil {
			return err
		}
		filename := filepath.Clean(args[0])
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		oid, err := SaveHashObj(content)
		if err != nil {
			return err
		}
		log.Printf("saved a hashed-object!!\noid: %s\n", oid)
		return nil
	},
}

// oid is an so-called object ID.
func IssueObjID(data []byte) (oid string) {
	s := sha1.Sum(data)
	oid = hex.EncodeToString(s[:])
	return oid
}

// covert bytes to object and save it under object storage(="current dir/.pgit/objects/")
func SaveHashObj(content []byte) (oid string, err error) {
	obj := NewObject("blob", IdentBlob, content)
	oid = IssueObjID(obj.Encode())
	f, err := os.Create(filepath.Join(PgitDir, ObjDir, oid))
	if err != nil {
		return "", fmt.Errorf("failed to create a file: (error: %w)", err)
	}
	f.Write(content)
	f.Close()

	return oid, nil
}

func init() {
	rootCmd.AddCommand(hashObjCmd)
}
