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

func joinTestDir(t *testing.T, name string) (rootPath string) {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	rootPath = filepath.Join(cwd, name)
	if err := os.Mkdir(rootPath, os.ModeDir); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(rootPath); err != nil {
		t.Fatal(err)
	}
	return rootPath
}

func leaveTestDir(t *testing.T, rootPath string) {
	t.Helper()

	if err := os.Chdir(filepath.Dir(rootPath)); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(rootPath); err != nil {
		t.Fatal(err)
	}
}

// creates all necessary directories
func initPgitForTest(t *testing.T) {
	t.Helper()

	if err := cmd.InitCmd.RunE(cmd.InitCmd, []string{}); err != nil {
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

	if err := cmd.ParseFlags(args); err != nil {
		t.Fatal(err)
	}
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
					{"dir", filepath.Join(cmd.PgitDir, cmd.RefDir)},
					{"dir", filepath.Join(cmd.PgitDir, cmd.RefDir, cmd.TagDir)},
					{"file", filepath.Join(cmd.PgitDir, cmd.RefDir, cmd.TagDir, "HEAD")},
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

func TestCommit(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{"-m", "test message"},
				out:  newWantOutput("", []output{}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				rootPath := joinTestDir(t, "commit")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})
				_, err := loadAndSetFiles(cmdDir, "*.go", rootPath)
				if err != nil {
					t.Fatal(err)
				}
				_, err = cmd.WriteTree(rootPath)
				if err != nil {
					t.Fatal(err)
				}
				ents, err := os.ReadDir(filepath.Join(rootPath, cmd.PgitDir, cmd.ObjDir))
				if err != nil {
					t.Fatal(err)
				}
				currentFileNum := len(ents)

				stdout, err := execCmd(t, cmd.CommitCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				afterEnts, err := os.ReadDir(filepath.Join(rootPath, cmd.PgitDir, cmd.ObjDir))
				if err != nil {
					t.Fatal(err)
				}
				gotFileNum := len(afterEnts)
				if gotFileNum != (currentFileNum + 1) {
					t.Errorf("file num should be equal: (gotNum: %d, wantNum: %d)", gotFileNum, currentFileNum+1)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

// This is a helper function for test. The feature is to load files from srdDir corresponding to pattern and save them in targetDir.
// The returned paths are filepaths saved in the targetDir.
func loadAndSetFiles(srcDir string, pattern string, targetDir string) (paths []string, err error) {
	fps, err := filepath.Glob(filepath.Join(srcDir, pattern))
	if err != nil {
		return nil, err
	}
	paths = make([]string, 0, len(fps))
	for _, fp := range fps {
		f, err := os.Open(fp)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		f.Close()
		path := filepath.Join(targetDir, filepath.Base(fp))
		testFile, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		testFile.Write(b)
		testFile.Close()

		paths = append(paths, path)
	}
	return paths, nil
}

func TestLog(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{},
				out:  newWantOutput("", []output{}),
			},
			{
				desc: "02_with @ ailias",
				args: []string{"@"},
				out:  newWantOutput("", []output{}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				rootPath := joinTestDir(t, "log")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})
				_, err := loadAndSetFiles(cmdDir, "*.go", rootPath)
				if err != nil {
					t.Fatal(err)
				}
				oid, err := cmd.NewCommit("test message")
				if err != nil {
					t.Fatal(err)
				}
				tt.out = newWantOutput(fmt.Sprintf("%s\n", oid), []output{})

				stdout, err := execCmd(t, cmd.LogCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestCheckout(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{},
				out:  newWantOutput("", []output{}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				rootPath := joinTestDir(t, "checkout")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})

				paths, err := loadAndSetFiles(cmdDir, "*.go", ".")
				if err != nil {
					t.Fatal(err)
				}
				outs := make([]output, 0, len(paths))
				for _, path := range paths {
					outs = append(outs, output{
						fileType: "file",
						path:     path,
					})
				}
				tt.out = newWantOutput("", outs)

				_, err = cmd.WriteTree(".")
				if err != nil {
					t.Fatal(err)
				}
				oid, err := cmd.NewCommit("test message")
				if err != nil {
					t.Fatal(err)
				}
				tt.args = []string{oid}
				if err := cmd.SweepDir("."); err != nil {
					t.Fatal(err)
				}

				stdout, err := execCmd(t, cmd.CheckoutCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestTag(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{},
				out:  newWantOutput("", []output{}),
			},
			{
				desc: "02_no oid",
				args: []string{"test"},
				out:  newWantOutput("", []output{}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				rootPath := joinTestDir(t, "tag")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})
				_, err := loadAndSetFiles(cmdDir, "*.go", rootPath)
				if err != nil {
					t.Fatal(err)
				}
				oid, err := cmd.NewCommit("test message")
				if err != nil {
					t.Fatal(err)
				}
				if len(tt.args) == 0 {
					tt.args = []string{"test", oid}
				}
				tt.out = newWantOutput("", []output{
					{
						fileType: "file",
						path:     filepath.Join(cmd.PgitDir, cmd.RefDir, cmd.TagDir, "test"),
					},
				})

				stdout, err := execCmd(t, cmd.TagCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestK(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_printing all refs with no arg",
				args: []string{},
				out:  newWantOutput("", []output{}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				rootPath := joinTestDir(t, "tag")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})
				_, err := loadAndSetFiles(cmdDir, "*.go", rootPath)
				if err != nil {
					t.Fatal(err)
				}
				oid, err := cmd.NewCommit("test message")
				if err != nil {
					t.Fatal(err)
				}
				if err := cmd.TagCmd.RunE(cmd.TagCmd, []string{"test", oid}); err != nil {
					t.Fatal(err)
				}

				stdout, err := execCmd(t, cmd.KCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				if stdout == "" {
					t.Errorf("stdout should not be empty")
				}
			})
		}
	})
}
