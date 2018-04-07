package gutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Skip()
	err := Copy(filepath.Join(os.TempDir(), "hotfix"), "testdata/hotfix")
	if err != nil {
		panic(err)
	}
}
