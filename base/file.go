package base

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// CheckFileExistBackInfo Get file info if is on disk.
func CheckFileExistBackInfo(file string, ignoreZero bool) os.FileInfo {
	// fi, err := os.Stat(file)
	fi, err := os.Lstat(file)
	if !os.IsNotExist(err) && !fi.IsDir() {
		if fi.Size() <= 0 && ignoreZero {
			return nil
		}
		return fi
	}
	return nil
}

// GetAllFile get all file under pathname and its sub directory
func GetAllFile(pathname string, s []string, ignoreZero bool) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := filepath.Join(pathname, fi.Name())
			s, err = GetAllFile(fullDir, s, ignoreZero)
			if err != nil {
				return s, err
			}
		} else {
			if fi.Size() <= 0 && ignoreZero {
				continue
			}
			fullName := filepath.Join(pathname, fi.Name())
			s = append(s, fullName)
		}
	}
	return s, nil
}

// GetFileNameExt get file name and ext from path
func GetFileNameExt(f string) (string, string) {
	dfPurName := filepath.Base(strings.TrimSpace(f))
	ext := ""
	lastDotIndex := strings.LastIndex(dfPurName, ".")
	if lastDotIndex > -1 {
		ext = strings.ToLower(strings.TrimPrefix(dfPurName[lastDotIndex:], "."))
		dfPurName = dfPurName[:lastDotIndex]
	}
	return dfPurName, ext
}

// GetFileMimeTypeExtByPath get file mimetype from file path
func GetFileMimeTypeExtByPath(f string) (mime *mimetype.MIME) {
	mime, _ = mimetype.DetectFile(f)
	return
}

// GetFileMimeTypeExtByReader get file mimetype from file reader
func GetFileMimeTypeExtByReader(f io.Reader) (mime *mimetype.MIME) {
	mime, _ = mimetype.DetectReader(f)
	return
}

// GetFileMimeTypeExtByBytes get file mimetype from file []byte
func GetFileMimeTypeExtByBytes(f []byte) *mimetype.MIME {
	return mimetype.Detect(f)
}

// CheckFileInMimeTypes check file mimetype in []string, such as []string{"text/plain", "text/html", "text/csv"}
func CheckFileInMimeTypes(mime *mimetype.MIME, mts []string) bool {
	return mimetype.EqualsAny(mime.String(), mts...)
}

// CheckFileIsBinaryByMimeType check file is binary or not by mimetype
func CheckFileIsBinaryByMimeType(m *mimetype.MIME) bool {
	isBinary := true
	for mime := m; mime != nil; mime = mime.Parent() {
		if mime.Is("text/plain") {
			isBinary = false
		}
	}
	return isBinary
}
