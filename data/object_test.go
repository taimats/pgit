package data_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/taimats/pgit/data"
)

// This is a helper function for test. The feature is to load files from srdDir corresponding to a pattern and
// save them in trgDir. The returned paths are filepaths saved in the targetDir.
func loadAndSetFiles(srcDir string, pattern string, trgDir string) (paths []string, err error) {
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
		path := filepath.Join(trgDir, filepath.Base(fp))
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

func TestSaveBlobObj(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()
		tests := []struct {
			desc    string
			trgPath string
			content []byte
		}{
			{
				desc:    "01_all set",
				trgPath: tmpDir,
				content: []byte("test message"),
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				oid, err := data.SaveBlobObj(tt.trgPath, tt.content)

				if err != nil {
					t.Errorf("should be nil: \nerror: %s", err)
				}
				CmpFileContent(t, filepath.Join(tt.trgPath, oid), tt.content)
			})
		}
	})
}

func TestWriteTree(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmdDir := t.TempDir()
		tests := []struct {
			desc       string
			srcDirPath string
			trgDirPath string
		}{
			{
				desc:       "01_all set",
				srcDirPath: tmdDir,
				trgDirPath: "./test/tree/trg",
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				paths, err := loadAndSetFiles(".", "*.go", tt.srcDirPath)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.MkdirAll(tt.trgDirPath, os.ModeDir); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					for _, p := range paths {
						os.RemoveAll(p)
					}
					os.RemoveAll(tt.trgDirPath)
				})

				_, err = data.WriteTree(tt.srcDirPath, tt.trgDirPath)

				if err != nil {
					t.Errorf("should be nil: \n{ error: %s }", err)
				}
				ents, err := os.ReadDir(tt.trgDirPath)
				if err != nil {
					t.Fatal(err)
				}
				fileNum := len(ents)
				if fileNum != len(paths)+1 {
					t.Errorf("file num should be equal: \n{ gotNum: %d, wantNum: %d }", fileNum, len(paths)+1)
				}

			})
		}
	})
}

func TestReadTree(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()
		paths, err := loadAndSetFiles(".", "*.go", tmpDir)
		if err != nil {
			t.Fatal(err)
		}
		srcDir := "./test/tree"
		if err := os.MkdirAll(srcDir, os.ModeDir); err != nil {
			t.Fatal(err)
		}
		treeOid, err := data.WriteTree(tmpDir, srcDir)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.RemoveAll(srcDir)
		})
		trgFileNum := len(paths)

		tests := []struct {
			desc       string
			treeOid    string
			srcDirPath string
			trgDirPath string
		}{
			{
				desc:       "01_all set",
				treeOid:    treeOid,
				srcDirPath: srcDir,
				trgDirPath: "./test/tree/trg",
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				if err := os.MkdirAll(tt.trgDirPath, os.ModeDir); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					os.RemoveAll(tt.trgDirPath)
				})

				err := data.ReadTree(tt.treeOid, tt.srcDirPath, tt.trgDirPath)

				if err != nil {
					t.Errorf("should be nil: \n{ error: %s }", err)
				}
				ents, err := os.ReadDir(tt.trgDirPath)
				if err != nil {
					t.Fatal(err)
				}
				fileNum := len(ents)
				if fileNum != trgFileNum {
					t.Errorf("file num should be equal: \n{ gotNum: %d, wantNum: %d }", fileNum, trgFileNum)
				}

			})
		}
	})
}

func TestGetCommit(t *testing.T) {
	tests := []struct {
		desc string
		oid  string
		want *data.Commit
	}{
		{
			desc: "",
			oid:  "testoid",
			want: &data.Commit{
				TreeOid: "testTreeOid",
				Parent:  "testParent",
				Msg:     "test message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tmpDir := filepath.Join(data.PgitDirBase, data.ObjDirBase)
			if err := os.MkdirAll(tmpDir, os.ModeDir); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { os.RemoveAll(filepath.Dir(tmpDir)) })

			content := fmt.Sprintf("tree %v\nparent %v\n\n%v\n", tt.want.TreeOid, tt.want.Parent, tt.want.Msg)
			if err := data.WriteFile(filepath.Join(tmpDir, tt.oid), []byte(content)); err != nil {
				t.Fatal(err)
			}

			got, err := data.GetCommit(tt.oid)

			if err != nil {
				t.Errorf("should be nil:\n{ error: %s }", err)
			}
			CmpStructs(t, got, tt.want)
		})
	}
}
