package data_test

import (
	"path/filepath"
	"testing"

	"github.com/taimats/pgit/data"
)

func TestDiffFiles(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			desc     string
			fromPath string
			toPath   string
			want     string
		}{
			{
				desc:     "01_all set",
				fromPath: "./test/diff/from",
				toPath:   "./test/diff/to",
				want:     "-The only difference is from.\n+The only difference is to.\n",
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got, err := data.DiffFiles(tt.fromPath, tt.toPath)

				if err != nil {
					t.Errorf("error should be nil:\nerror:%s\n", err)
				}
				CmpStructs(t, got, tt.want)
			})
		}
	})
}

func TestDiffTrees(t *testing.T) {
	srcDir := filepath.Join("./test", "difftrees")
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			desc   string
			srcDir string
			from   data.Tree
			to     data.Tree
			want   []*data.Diff
		}{
			{
				desc:   "01_all set",
				srcDir: srcDir,
				from: data.Tree{
					"file_01": &data.TreeElem{
						ObjType: data.ObjTypeBlob,
						Oid:     "testoid_01",
						Name:    "file_01",
						Child:   nil,
					},
				},
				to: data.Tree{
					"file_01": &data.TreeElem{
						ObjType: data.ObjTypeBlob,
						Oid:     "testoid_02",
						Name:    "file_01",
						Child:   nil,
					},
				},
				want: []*data.Diff{
					{
						Filename: "file_01",
						Diff:     "-This is a test message.\n+This is a text message.\n",
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got, err := data.DiffTrees(tt.from, tt.to, tt.srcDir)

				if err != nil {
					t.Errorf("error should be nil\n{ error: %s }\n", err)
				}
				CmpStructs(t, got, tt.want)
			})
		}
	})
}
