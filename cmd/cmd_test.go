package cmd_test

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/taimats/pgit/cmd"
)

// Each test of cmd package is highly likely to disturb the current working directory.
// This will negatively affect each test case and produce no valid outcome.
// In light of this, TestMain prepares a temporary directory for test and execute each test there.
// After all tests are done, TestMain is expected to clean up the prep.
func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	newDir := filepath.Join(cwd, "testDir")
	if err := os.Mkdir(newDir, os.ModeDir); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(newDir); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := os.Chdir(cwd); err != nil {
		log.Println(err)
	}
	if err := os.RemoveAll(newDir); err != nil {
		log.Println(err)
	}

	os.Exit(code)
}

// creates all necessary directories and returns a fixed path(= "...rootDir/.pgit/objects")
func initPgitForTest(t *testing.T) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(cmd.PgitDir, cmd.ObjDir), os.ModeDir); err != nil {
		t.Fatal(err)
	}
}

func removePgitDirForTest(t *testing.T) {
	t.Helper()
	os.RemoveAll(filepath.Join(cmd.PgitDir))
}

func newObjID(data []byte) string {
	s := sha1.Sum(data)
	return hex.EncodeToString(s[:])
}

func newBlobObj(t *testing.T, content []byte) (path string, oid string) {
	t.Helper()
	obj := cmd.NewObject("blob", cmd.IdentBlob, content)
	oid = newObjID(obj.Encode())
	return filepath.Join(cmd.PgitDir, cmd.ObjDir, oid), oid
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
	t.Cleanup(func() {
		for _, out := range output.outs {
			os.Remove(out.path)
		}
	})

	if output.stdout != "" {
		if stdout != output.stdout {
			t.Errorf("Stdout should be equal: (got=%s, want=%s)", stdout, output.stdout)
		}
	}
	if len(output.outs) != 0 {
		for _, o := range output.outs {
			fi, err := os.Stat(o.path)
			if err != nil {
				t.Errorf("wantFile should be generated: (path: %s, error: %s)", o.path, err)
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
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{},
				out: newWantOutput("", []output{
					{"dir", cmd.PgitDir},
					{"dir", filepath.Join(cmd.PgitDir, cmd.ObjDir)},
				}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				t.Cleanup(func() {
					removePgitDirForTest(t)
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
		content := `This is a message for hash object test.`
		f, err := os.Create("test")
		if err != nil {
			t.Fatal(err)
		}
		f.WriteString(content)
		f.Close()
		t.Cleanup(func() { os.Remove("test") })

		data := cmd.NewObject("blob", cmd.IdentBlob, []byte(content)).Encode()
		oid := newObjID(data)

		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{"test"},
				out: newWantOutput("", []output{
					{"file", filepath.Join(cmd.PgitDir, cmd.ObjDir, oid)},
				}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				initPgitForTest(t)
				t.Cleanup(func() {
					removePgitDirForTest(t)
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
		content := `This is a message for cat-file test.`
		_, oid := newBlobObj(t, []byte(content))
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{oid},
				out: newWantOutput("", []output{
					{"file", filepath.Join(cmd.PgitDir, cmd.ObjDir, oid)},
				}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				initPgitForTest(t)
				t.Cleanup(func() {
					removePgitDirForTest(t)
				})
				_, err := cmd.SaveHashObj([]byte(content))
				if err != nil {
					t.Fatal(err)
				}

				stdout, err := execCmd(t, cmd.CatFileCmd, tt.args)

				if err != nil {
					t.Errorf("error should be empty: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestWriteTree(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	fps, err := filepath.Glob(filepath.Join(cmdDir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, fp := range fps {
		f, err := os.Open(fp)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		_, oid := newBlobObj(t, b)
		baseName := filepath.Base(fp)
		fmt.Fprintf(&buf, "%s %s %s\n", cmd.ObjTypeBlob, oid, baseName) //e.g. "blob oid hoge.txt"

		testFile, err := os.Create(filepath.Join(cwd, baseName))
		if err != nil {
			t.Fatal(err)
		}
		testFile.Write(b)
		defer testFile.Close()
	}
	treePath, _ := newBlobObj(t, buf.Bytes())
	wantOutput := newWantOutput("", []output{{"file", treePath}})

	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{},
				out:  wantOutput,
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				initPgitForTest(t)
				t.Cleanup(func() {
					removePgitDirForTest(t)
				})

				stdout, err := execCmd(t, cmd.WriteTreeCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestReadTree(t *testing.T) {
	//setting up temporary files in the test directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	fps, err := filepath.Glob(filepath.Join(cmdDir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	outs := make([]output, 0, len(fps))
	for _, fp := range fps {
		f, err := os.Open(fp)
		if err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		baseName := filepath.Base(fp)
		testFile, err := os.Create(filepath.Join(cwd, baseName))
		if err != nil {
			t.Fatal(err)
		}
		testFile.Write(b)
		testFile.Close()

		outs = append(outs, output{
			fileType: "file",
			path:     filepath.Join(cwd, baseName),
		})
	}

	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{},
				out:  newWantOutput("", outs),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				initPgitForTest(t)
				t.Cleanup(func() { removePgitDirForTest(t) })

				oid, err := cmd.WriteTree(".")
				if err != nil {
					t.Fatal(err)
				}
				tt.args = []string{oid}

				stdout, err := execCmd(t, cmd.ReadTreeCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}
