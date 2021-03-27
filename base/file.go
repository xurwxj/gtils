package base

import (
	"bytes"
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

// GetFileOrDirectoryInfo Get file  or directory info if is on disk.
func GetFileOrDirectoryInfo(targetPath string) os.FileInfo {
	// fi, err := os.Stat(file)
	fi, err := os.Lstat(targetPath)
	if !os.IsNotExist(err) {
		return fi
	}
	return nil
}

// WriteFile writes data to a file named by name under path.
// WrtieFile creates the path and file if not exist, truncating the file if already exists.
func WriteFile(path, name string, data []byte) error {
	CheckPathExistOrCreate(path)
	out, err := os.Create(filepath.Join(path, name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}

// CheckPathExistOrCreate creates tPath if it does not exist.
func CheckPathExistOrCreate(tPath string) {
	if _, err := os.Stat(tPath); os.IsNotExist(err) {
		os.MkdirAll(tPath, os.ModePerm)
	}
}

// CheckPathExistOrCreateWithMode check path exist or not, create when not exist
func CheckPathExistOrCreateWithMode(tPath, mode string) {
	if _, err := os.Stat(tPath); os.IsNotExist(err) {
		os.MkdirAll(tPath, os.ModePerm)
	}
}

// CheckFileExsit Get download info from bbolt and check if is on disk.
func CheckFileExsit(file string) bool {
	_, err := os.Lstat(file)
	return !os.IsNotExist(err)
}

// CheckIsFile Get download info from bbolt and check if is on disk.
func CheckIsFile(file string) bool {
	fi, err := os.Lstat(file)
	if err == nil && !fi.IsDir() {
		return true
	}
	return false
}

// GetAllDirectoryName get all directory name under pathname and its sub directory
func GetAllDirectoryName(pathname string, s []string, recursion bool) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			s = append(s, fi.Name())
			if recursion {
				fullDir := filepath.Join(pathname, fi.Name())
				s, err = GetAllFile(fullDir, s, recursion)
				if err != nil {
					return s, err
				}
			}
		}
	}
	return s, nil
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

// GetAllFileInDirectory get all file under pathname and its sub directory
func GetAllFileInDirectory(pathname string, s, exts, excludeDirs []string, includeEXTs, ignoreZero bool) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			dirName := fi.Name()
			if checkContainDir(excludeDirs, dirName) {
				continue
			}
			fullDir := filepath.Join(pathname, dirName)
			s, err = GetAllFileInDirectory(fullDir, s, exts, excludeDirs, includeEXTs, ignoreZero)
			if err != nil {
				return s, err
			}
		} else {
			if fi.Size() <= 0 && ignoreZero {
				continue
			}
			fileName := fi.Name()
			_, fileExt := GetFileNameExt(fileName)
			if (includeEXTs && FindInStringSlice(exts, fileExt)) || (!includeEXTs && !FindInStringSlice(exts, fileExt)) || len(exts) == 0 {
				fullName := filepath.Join(pathname, fileName)
				s = append(s, fullName)
			}
		}
	}
	return s, nil
}

func checkContainDir(excludeDirs []string, dirName string) bool {
	if len(excludeDirs) > 0 {
		for _, d := range excludeDirs {
			if strings.Contains(dirName, d) {
				return true
			}
		}
	}
	return false
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
