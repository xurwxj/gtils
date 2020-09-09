package base

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
