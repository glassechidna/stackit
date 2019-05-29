package zipper

import (
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func Zip(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrap(err, "determining absolute path")
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", errors.Errorf("no file exists at '%s'", path)
	}

	if isZip(path) {
		return path, nil
	}

	fw, err := ioutil.TempFile("", fmt.Sprintf("%s*.zip", filepath.Base(path)))
	if err != nil {
		return "", errors.Wrap(err, "creating temporary file to write to")
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)

	if fi.IsDir() {
		err = filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			return addFileToZip(zw, path, subpath)
		})
	} else {
		err = addFileToZip(zw, path, path)
	}
	if err != nil {
		return "", err
	}

	err = zw.Close()
	if err != nil {
		return "", errors.Wrap(err, "finialising zip file")
	}

	return fw.Name(), nil
}

func addFileToZip(zw *zip.Writer, topPath, path string) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	fh, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(topPath, path)
	if err != nil {
		return err
	}
	if relPath == "." {
		relPath = filepath.Base(path)
	}

	fh.Name = relPath
	fh.Method = zip.Deflate

	zf, err := zw.CreateHeader(fh)
	if err != nil {
		return err
	}

	_, err = io.Copy(zf, inputFile)
	return err
}

func isZip(path string) bool {
	if fi, err := os.Stat(path); fi != nil && fi.IsDir() {
		return false
	} else if os.IsNotExist(err) {
		return false
	} else {
		f, err := os.Open(path)
		if err != nil {
			return false
		}
		defer f.Close()

		buf := make([]byte, 4)
		_, err = io.ReadAtLeast(f, buf, 4)
		if err != nil {
			return false
		}

		contentType := http.DetectContentType(buf)
		return contentType == "application/zip"
	}
}
