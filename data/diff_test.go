package data_test

import (
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
