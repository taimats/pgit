package data_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/taimats/pgit/data"
)

func CmpStructs(t *testing.T, got any, want any) {
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("should be equal: %s", diff)
	}
}

func newTestRef(t *testing.T, path string) *data.Ref {
	t.Helper()

	ref, err := data.NewRef(path)
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

func CmpFileContent(t *testing.T, srcPath string, want []byte) {
	t.Helper()

	f, err := os.Open(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	src, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	if diff := cmp.Diff(src, want); diff != "" {
		t.Errorf("File content should be equal: \n%s", diff)
	}
}

func cleanupFileContent(t *testing.T, path []string, updatedContent []byte) {
	for _, p := range path {
		f, err := os.Create(p)
		if err != nil {
			t.Log(err)
		}
		f.Write(updatedContent)
		f.Close()
	}
}

func TestNewRef(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		tests := []struct {
			desc string
			path string
			want *data.Ref
		}{
			{
				desc: "01_creates a symbolic ref",
				path: "./test/NewRef_Success_01",
				want: &data.Ref{
					Path:       "./test/NewRef_Success_01",
					Oid:        "",
					IsSymbolic: true,
					Next:       ".pgit/refs/heads/test",
				},
			},
			{
				desc: "02_creates a direct ref (pointing to an object id directly)",
				path: "./test/NewRef_Success_02",
				want: &data.Ref{
					Path:       "./test/NewRef_Success_02",
					Oid:        "testoid",
					IsSymbolic: false,
					Next:       "",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got, err := data.NewRef(tt.path)
				if err != nil {
					t.Errorf("should be nil: (error: %s)", err)
				}
				CmpStructs(t, got, tt.want)
			})
		}
	})
}

func TestRefResolveSymbolic(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		tests := []struct {
			desc string
			path string
			want *data.Ref
		}{
			{
				desc: "01_resolved by the way of three routes",
				path: "./test/ref/route_01",
				want: newTestRef(t, "test/ref/route_03"),
			},
			{
				desc: "02_resolved directly (without any middle route)",
				path: "test/ref/route_03",
				want: nil,
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				ref := newTestRef(t, tt.path)

				got, err := ref.ResolveSymbolic(ref.Next)

				if err != nil {
					t.Errorf("should be nil: (error: %s)", err)
				}
				CmpStructs(t, got, tt.want)
			})
		}
	})
}

func TestRefUpdate(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		tests := []struct {
			desc           string
			path           string
			oid            string
			outPath        string
			cleanupContent []byte
		}{
			{
				desc:           "01_updated with a symbolic ref",
				path:           "./test/ref/update_symbolic",
				oid:            "updatedoid",
				outPath:        "./test/ref/update_out",
				cleanupContent: []byte{},
			},
			{
				desc:           "02_updated with a direct ref",
				path:           "./test/ref/update_direct",
				oid:            "updatedoid",
				outPath:        "./test/ref/update_direct",
				cleanupContent: []byte{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				t.Cleanup(func() {
					cleanupFileContent(t, []string{tt.outPath}, tt.cleanupContent)
				})
				ref := newTestRef(t, tt.path)

				err := ref.Update(tt.oid)

				if err != nil {
					t.Errorf("should be nil: (error: %s)", err)
				}
				CmpFileContent(t, tt.outPath, []byte(tt.oid))
			})
		}
	})
}

func TestRefUpdateSymbolic(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		tests := []struct {
			desc           string
			sutPath        string
			refPath        string
			cleanupContent []byte
		}{
			{
				desc:           "01_updateSymbolic",
				sutPath:        "./test/ref/updateSymbolic",
				refPath:        "./test/update/symbolic",
				cleanupContent: []byte{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				t.Cleanup(func() {
					cleanupFileContent(t, []string{tt.sutPath}, tt.cleanupContent)
				})
				ref := newTestRef(t, tt.sutPath)

				err := ref.UpdateSymbolic(tt.refPath)

				if err != nil {
					t.Errorf("should be nil: (error: %s)", err)
				}
				CmpFileContent(t, tt.sutPath, []byte(fmt.Sprintf("ref: %s <- HEAD\n", tt.refPath)))
			})
		}
	})
}
