package zipper

import (
	"archive/zip"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsZip(t *testing.T) {
	assert.False(t, isZip("testdata/notzip.txt"))
	assert.False(t, isZip("testdata/short.txt"))
	assert.False(t, isZip("testdata/dir"))
	assert.True(t, isZip("testdata/realzip.zip"))
	assert.True(t, isZip("./testdata/../testdata/realzip.zip"))
}

func TestZipMaintainsPermissions(t *testing.T) {
	path, err := Zip("testdata")
	assert.NoError(t, err)

	zf, err := zip.OpenReader(path)
	assert.NoError(t, err)

	foundExecutable := false

	for _, f := range zf.File {
		if f.Name == "executable.sh" {
			foundExecutable = true
			assert.Equal(t, f.Mode().Perm(), os.FileMode(0755))
		}
	}

	assert.True(t, foundExecutable)
}
