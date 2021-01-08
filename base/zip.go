package base

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

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
func Unzip(filename, d string) error {
	// z := archiver.Zip{
	// 	MkdirAll:               true,
	// 	OverwriteExisting:      true,
	// 	ImplicitTopLevelFolder: false,
	// }
	// return z.Unarchive(filename, d)
	return unzip(filename, d, "GBK")
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
		for _, file := range zipReader.Reader.File {
			fHeader := file.FileHeader

			tname := file.Name
			if fHeader.NonUTF8 {
				if HasGBK(tFile) {
					tname = DecodingGBKString(tname)
				} else if HasJP(tFile) {
					tname = DecodingJPString(tname)
				} else {
					tname = DecodingFromString(tname)
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
