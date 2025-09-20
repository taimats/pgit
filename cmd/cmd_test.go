package cmd_test

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/cmd"
)

// creates all necessary directories and returns a fixed path(= "...rootDir/.pgit/objects")
func initPgitForTest(t *testing.T) (objDirPath string) {
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

	objdir, err := cmd.AbsObjDirPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Dir(objdir)); err != nil {
		t.Log(err)
	}
}

func newObjID(data []byte) string {
	s := sha1.Sum(data)
	return hex.EncodeToString(s[:])
}

type testCase struct {
	desc string
	args []string
	out  wantOutput
}

// the main test for examining the common behavior of each command
func execCmd(t *testing.T, cmd *cobra.Command, args []string) (stdout string, err error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	os.Stdout = w

	cmd.SetArgs(args)
	err = cmd.RunE(cmd, args)

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	return buf.String(), err
}

type output struct {
	fileType string //"file" or "dir"
	path     string
}

type wantOutput struct {
	stdout string
	outs   []output
}

func newWantOutput(stdout string, outs []output) wantOutput {
	return wantOutput{
		stdout: stdout,
		outs:   outs,
	}
}

func assertOutput(t *testing.T, stdout string, output wantOutput) {
	if output.stdout != "" {
		if stdout != output.stdout {
			t.Errorf("Stdout should be equal: (got=%s, want=%s)", stdout, output.stdout)
		}
	}
	if len(output.outs) != 0 {
		for _, o := range output.outs {
			fi, err := os.Stat(o.path)
			if err != nil {
				t.Errorf("file should be generated: (error: %s, path: %s)", err, o.path)
				return
			}
			switch o.fileType {
			case "dir":
				if !fi.IsDir() {
					t.Errorf("file should be dir: (got: %s)", fi.Name())
				}
			case "file":
				if fi.IsDir() {
					t.Errorf("file should not be dir: (got: %s)", fi.Name())
				}
			}
		}
	}
}

func TestInitCMD(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		objdir, err := cmd.AbsObjDirPath()
		if err != nil {
			t.Fatal(err)
		}
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{},
				out: newWantOutput("", []output{
					{"dir", filepath.Dir(objdir)},
					{"dir", objdir},
				}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				t.Cleanup(func() {
					removeTmpPgitDir(t)
				})

				stdout, err := execCmd(t, cmd.InitCmd, tt.args)

				if err != nil {
					t.Errorf("error should be nil: {error: %s}", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestHashObjectCmd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		objdir, err := cmd.AbsObjDirPath()
		if err != nil {
			t.Fatal(err)
		}
		content := `This is a message for test.`
		oid := newObjID([]byte(content))

		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{"test"},
				out: newWantOutput("", []output{
					{"file", filepath.Join(objdir, oid)},
				}),
			},
			{
				desc: "02_with type flag",
				args: []string{"test", "--type", "blob"},
				out: newWantOutput("", []output{
					{"file", filepath.Join(objdir, oid)},
				}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				t.Cleanup(func() { removeTmpPgitDir(t) })

				pgitdirpath := filepath.Base(initPgitForTest(t))
				rootdir := filepath.Dir(pgitdirpath)

				f, err := os.Create(filepath.Join(rootdir, "test"))
				if err != nil {
					t.Fatal(err)
				}
				f.WriteString(content)
				f.Close()

				t.Cleanup(func() {
					if err := os.RemoveAll(f.Name()); err != nil {
						t.Log(err)
					}
				})

				stdout, err := execCmd(t, cmd.HashObjectCmd, tt.args)

				if err != nil {
					t.Errorf("error should be empty: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestCatFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		content := `This is a message for test.`
		oid := newObjID([]byte(content))

		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{oid},
				out: newWantOutput(
					"This is a message for test.\n",
					[]output{},
				),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				objdir := initPgitForTest(t)

				f, err := os.Create(filepath.Join(objdir, oid))
				if err != nil {
					t.Fatal(err)
				}
				f.WriteString(content)
				f.Close()

				t.Cleanup(func() {
					removeTmpPgitDir(t)
					if err := os.RemoveAll(f.Name()); err != nil {
						t.Log(err)
					}
				})

				stdout, err := execCmd(t, cmd.CatFileCmd, tt.args)

				if err != nil {
					t.Errorf("error should be empty: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}
