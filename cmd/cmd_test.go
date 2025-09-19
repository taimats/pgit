package cmd_test

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/cmd"
)

type outputFile struct {
	fileType string // dir or file
	path     string // "/dir/file.txt" etc
}

func newoutputFile(fileType string, path string) *outputFile {
	return &outputFile{fileType: fileType, path: path}
}

func (of *outputFile) fileInfo() os.FileInfo {
	fi, err := os.Stat(of.path)
	if err != nil {
		return nil
	}
	return fi
}

func (of *outputFile) newFile() (*os.File, error) {
	if of.fileType == "dir" {
		return nil, errors.New("Cannot genertate a file: this filetype is \"directory\"")
	}
	f, err := os.Create(of.path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (of *outputFile) newDir() error {
	if of.fileType == "file" {
		return errors.New("Cannot genertate a directory: this filetype is \"file\"")
	}
	return os.MkdirAll(of.path, 0750)
}

func (of *outputFile) remove() error {
	if err := os.RemoveAll(of.path); err != nil {
		return err
	}
	return nil
}

type testCase struct {
	desc    string //description
	args    []string
	genFile bool
	files   []*outputFile
	hasErr  bool
}

// Testing the main behavior of command in a normal setting
func testCmd(t *testing.T, tc testCase, cmd *cobra.Command) {
	t.Cleanup(
		func() {
			if len(tc.files) > 0 {
				for _, f := range tc.files {
					if err := f.remove(); err != nil {
						t.Log(err)
					}
				}
			}
		},
	)

	err := cmd.RunE(cmd, tc.args)

	if tc.genFile {
		for _, f := range tc.files {
			fi := f.fileInfo()
			if fi == nil {
				t.Error("some file should be generated: (got: nil)")
			}
			switch f.fileType {
			case "dir":
				if !fi.IsDir() {
					t.Errorf("A generated file should be a directory: (got: %s)", fi.Name())
				}
			case "file":
				if fi.IsDir() {
					t.Errorf("A generated file should not be a directory: (got: %s)", fi.Name())
				}
			}
			if fi.Mode().Perm() != 511 {
				t.Errorf(
					"%s dir's permission is not valid: (got: %d, want: %d)",
					fi.Name(),
					fi.Mode().Perm(),
					os.FileMode(0511),
				)
			}
		}
	}
	if tc.hasErr {
		if err == nil {
			t.Errorf("error should not be nil")
		}
	}
}

// creates all necessary directories and returns a fixed path(= "/rootDir/.pgit/objects")
func initPgitForTest(t *testing.T) string {
	t.Helper()

	path, err := cmd.AbsObjDirPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(path, 0750); err != nil {
		t.Fatal(err)
	}
	return path
}

func removeTmpPgitDir(t *testing.T) {
	t.Helper()

	if err := os.RemoveAll(cmd.PgitDir); err != nil {
		t.Log(err)
	}
}

func newObjID(data []byte) string {
	s := sha1.Sum(data)
	return hex.EncodeToString(s[:])
}

func TestInitCMD(t *testing.T) {
	tests := []testCase{
		{
			desc:    "success_01",
			args:    []string{},
			genFile: true,
			files: []*outputFile{
				newoutputFile("dir", cmd.PgitDir),
				newoutputFile("dir", filepath.Join(cmd.PgitDir, cmd.ObjDir)),
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testCmd(t, tt, cmd.InitCmd)
		})
	}
}

func TestHashObjectCmd(t *testing.T) {
	objdir := initPgitForTest(t)
	t.Cleanup(func() { removeTmpPgitDir(t) })

	arg := filepath.Join("./test")
	content := `This is a message for test.`
	oid := newObjID([]byte(content))

	f, err := os.Create(arg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("file name: %s", f.Name())
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() {
		if err := os.RemoveAll(f.Name()); err != nil {
			t.Log(err)
		}
	})

	err = cmd.HashObjectCmd.RunE(cmd.HashObjectCmd, []string{arg})

	if err != nil {
		t.Errorf("should be nil: (error: %s)", err)
	}
	fi, err := os.Stat(filepath.Join(objdir, oid))
	if err != nil {
		t.Errorf("file should be generated")
	}
	if fi.IsDir() {
		t.Errorf("file should not be a directory")
	}
}
