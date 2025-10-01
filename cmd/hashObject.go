package cmd

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/data"
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
		content, err := data.ReadAllFileContent(filename)
		if err != nil {
			return fmt.Errorf("hash-object: internal error: %w", err)
		}
		oid, err := SaveHashObj(content)
		if err != nil {
			return err
		}
		log.Printf("saved a hashed-object!!\noid: %s\n", oid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashObjCmd)
}
