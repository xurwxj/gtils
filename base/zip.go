package base

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/mholt/archiver/v3"
)

// Zip zip d to filename
// d can contain directory or file
func Zip(d []string, filename string) error {
	z := archiver.Zip{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	return z.Archive(d, filename)
}

// Unzip unzip file according by ext
// d zip file path
func Unzip(filename, d, targetCharset string) error {
	// z := archiver.Zip{
	// 	MkdirAll:               true,
	// 	OverwriteExisting:      true,
	// 	ImplicitTopLevelFolder: false,
	// }
	// return z.Unarchive(filename, d)
	return unzip(filename, d, targetCharset)
}

// HasFilesInZip check file exist or not in spec exts
func HasFilesInZip(filename string, exts []string) (bool, error) {
	z := archiver.Zip{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	has := false
	err := z.Walk(filename, func(f archiver.File) error {
		// fmt.Println(f.Name())
		if FindInStringSlice(exts, strings.TrimSpace(strings.ToLower(strings.TrimPrefix(path.Ext(f.Name()), ".")))) && f.Size() > 0 && !f.IsDir() {
			has = true
		}
		return nil
	})
	return has, err
}

func unzip(tFile, targetDir, targetCharset string) error {
	// zip.FileInfoHeader(tFile)
	zipReader, err := zip.OpenReader(tFile)
	if err != nil {
		return err
	}
	if zipReader != nil {
		targetCharset = strings.TrimSpace(strings.ToLower(targetCharset))
		for _, file := range zipReader.Reader.File {
			fHeader := file.FileHeader

			tname := file.Name
			validUTF8, requireUTF8 := detectUTF8(tname)
			if fHeader.NonUTF8 && (!validUTF8 || !requireUTF8) {
				switch targetCharset {
				case "ja", "ja-jp":
					tname = DecodingJPString(tname)
				case "zh", "cn", "zh-cn":
					tname = DecodingGBKString(tname)
				default:
					if HasJP(tFile) || HasJPReg(tFile) {
						tname = DecodingJPString(tname)
					} else if HasGBK(tFile) {
						tname = DecodingGBKString(tname)
					} else {
						tname = DecodingFromString(tname)
					}
				}
			}

			zippedFile, err := file.Open()
			if err != nil {
				return err
			}
			defer zippedFile.Close()

			extractedFilePath := filepath.Join(
				targetDir,
				tname,
			)

			if file.FileInfo().IsDir() {
				// fmt.Println("Directory Created:", extractedFilePath)
				os.MkdirAll(extractedFilePath, file.Mode())
			} else {
				upPath := filepath.Dir(extractedFilePath)
				os.MkdirAll(upPath, os.ModePerm)

				outputFile, err := os.OpenFile(
					extractedFilePath,
					os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
					file.Mode(),
				)
				if err != nil {
					return err
				}
				defer outputFile.Close()

				_, err = io.Copy(outputFile, zippedFile)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("not a valid zip file for %s", tFile)
	}
	return nil
}

// detectUTF8 reports whether s is a valid UTF-8 string, and whether the string
// must be considered UTF-8 encoding (i.e., not compatible with CP-437, ASCII,
// or any other common encoding).
func detectUTF8(s string) (valid, require bool) {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		// Officially, ZIP uses CP-437, but many readers use the system's
		// local character encoding. Most encoding are compatible with a large
		// subset of CP-437, which itself is ASCII-like.
		//
		// Forbid 0x7e and 0x5c since EUC-KR and Shift-JIS replace those
		// characters with localized currency and overline characters.
		if r < 0x20 || r > 0x7d || r == 0x5c {
			if !utf8.ValidRune(r) || (r == utf8.RuneError && size == 1) {
				return false, false
			}
			require = true
		}
	}
	return true, require
}
