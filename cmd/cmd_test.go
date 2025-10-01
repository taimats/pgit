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
	"github.com/taimats/pgit/data"
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

func leaveTestDir(t *testing.T, currentPath string) {
	t.Helper()

	if err := os.Chdir(filepath.Dir(currentPath)); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(currentPath); err != nil {
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
	return filepath.Join(cmd.ObjDir, oid), oid
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
					{"dir", cmd.ObjDir},
					{"dir", cmd.RefDir},
					{"file", data.RefHEADPath},
					{"dir", cmd.TagDir},
					{"dir", cmd.HeadDir},
					{"file", filepath.Join(cmd.HeadDir, "master")},
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
					{"file", filepath.Join(cmd.ObjDir, oid)},
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
					{"file", filepath.Join(cmd.ObjDir, oid)},
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
				rootPath := joinTestDir(t, "writeTree")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})
				paths, err := loadAndSetFiles(cmdDir, "*.go", rootPath)
				if err != nil {
					t.Fatal(err)
				}

				stdout, err := execCmd(t, cmd.WriteTreeCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				ents, err := os.ReadDir(filepath.Join(cmd.ObjDir))
				if err != nil {
					t.Fatal(err)
				}
				if len(ents) != len(paths)+1 {
					t.Errorf("file num Not equal: (gotNum: %d, wantNum: %d)", len(ents), len(paths)+1)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func TestReadTree(t *testing.T) {
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
				rootPath := joinTestDir(t, "readTree")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})
				_, err := loadAndSetFiles(cmdDir, "*.go", rootPath)
				if err != nil {
					t.Fatal(err)
				}
				oid, err := data.WriteTree(rootPath, cmd.ObjDir)
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
				currentNum := len(allFileNames(t, cmd.ObjDir))

				stdout, err := execCmd(t, cmd.CommitCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				afterNum := len(allFileNames(t, cmd.ObjDir))
				//Two files (tree and commit) should be added to the object storage.
				if afterNum != currentNum+2 {
					t.Errorf("fileNum should be equal:\n{ gotNum: %d, wantNum: %d }", afterNum, currentNum)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}

func allFileNames(t *testing.T, dirPath string) []string {
	t.Helper()

	ents, err := os.ReadDir(dirPath)
	if err != nil {
		t.Fatal(err)
	}
	fns := make([]string, 0, len(ents))
	for _, ent := range ents {
		fns = append(fns, ent.Name())
	}
	return fns
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
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all well done",
				args: []string{"test"},
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
				_, err := cmd.NewCommit("test commit")
				if err != nil {
					t.Fatal(err)
				}
				branch, err := cmd.NewBranch("test")
				if err != nil {
					t.Fatal(err)
				}

				stdout, err := execCmd(t, cmd.CheckoutCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				headRef, err := data.ReadValueFromFile(data.RefHEADPath, []byte("ref:"))
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(headRef, []byte(branch)) {
					t.Errorf("checkout branch should be equal:\n{ got: %s, want: %s }\n", string(headRef), branch)
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
						path:     filepath.Join(cmd.TagDir, "test"),
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

// func TestK(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		tests := []testCase{
// 			{
// 				desc: "01_printing all refs with no arg",
// 				args: []string{"second"},
// 				out:  newWantOutput("", []output{}),
// 			},
// 		}
// 		for _, tt := range tests {
// 			t.Run(tt.desc, func(t *testing.T) {
// 				rootPath := joinTestDir(t, "k")
// 				initPgitForTest(t)
// 				t.Cleanup(func() {
// 					leaveTestDir(t, rootPath)
// 				})
// 				var buf strings.Builder
// 				oid, err := cmd.NewCommit()

// 				stdout, err := execCmd(t, cmd.KCmd, tt.args)

// 				if err != nil {
// 					t.Errorf("error should be emtpy: (error: %s)", err)
// 				}
// 				assertOutput(t, stdout, tt.out)
// 			})
// 		}
// 	})
// }

func TestBranch(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdDir := filepath.Dir(cwd)
	t.Run("success", func(t *testing.T) {
		tests := []testCase{
			{
				desc: "01_all set",
				args: []string{"test"},
				out: newWantOutput("", []output{
					{fileType: "file", path: filepath.Join(cmd.HeadDir, "test")},
				}),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				rootPath := joinTestDir(t, "branch")
				initPgitForTest(t)
				t.Cleanup(func() {
					leaveTestDir(t, rootPath)
				})
				_, err := loadAndSetFiles(cmdDir, "*.go", rootPath)
				if err != nil {
					t.Fatal(err)
				}

				stdout, err := execCmd(t, cmd.BranchCmd, tt.args)

				if err != nil {
					t.Errorf("error should be emtpy: (error: %s)", err)
				}
				assertOutput(t, stdout, tt.out)
			})
		}
	})
}
