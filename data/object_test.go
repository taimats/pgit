package data_test

import (
	"path/filepath"
	"testing"

	"github.com/taimats/pgit/data"
)

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
