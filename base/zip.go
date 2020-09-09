package base

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
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
	var err error
	zipReader, _ := zip.OpenReader(tFile)
	if zipReader != nil {
		for _, file := range zipReader.Reader.File {
			fHeader := file.FileHeader
			fname := []byte(file.Name)
			if fHeader.NonUTF8 {
				// fmt.Println("file.Name: ", file.Name)
				fname, err = Decode([]byte(file.Name), targetCharset)
				// fmt.Println("file.Name decode: ", string(fname))
				if err != nil {
					return err
				}
			}
			tname := string(fname)

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
				// fmt.Println("File extracted:", tname)

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

var encodings = map[string]encoding.Encoding{
	"ibm866":            charmap.CodePage866,
	"iso8859_2":         charmap.ISO8859_2,
	"iso8859_3":         charmap.ISO8859_3,
	"iso8859_4":         charmap.ISO8859_4,
	"iso8859_5":         charmap.ISO8859_5,
	"iso8859_6":         charmap.ISO8859_6,
	"iso8859_7":         charmap.ISO8859_7,
	"iso8859_8":         charmap.ISO8859_8,
	"iso8859_8I":        charmap.ISO8859_8I,
	"iso8859_10":        charmap.ISO8859_10,
	"iso8859_13":        charmap.ISO8859_13,
	"iso8859_14":        charmap.ISO8859_14,
	"iso8859_15":        charmap.ISO8859_15,
	"iso8859_16":        charmap.ISO8859_16,
	"koi8r":             charmap.KOI8R,
	"koi8u":             charmap.KOI8U,
	"macintosh":         charmap.Macintosh,
	"windows874":        charmap.Windows874,
	"windows1250":       charmap.Windows1250,
	"windows1251":       charmap.Windows1251,
	"windows1252":       charmap.Windows1252,
	"windows1253":       charmap.Windows1253,
	"windows1254":       charmap.Windows1254,
	"windows1255":       charmap.Windows1255,
	"windows1256":       charmap.Windows1256,
	"windows1257":       charmap.Windows1257,
	"windows1258":       charmap.Windows1258,
	"macintoshcyrillic": charmap.MacintoshCyrillic,
	"gbk":               simplifiedchinese.GBK,
	"gb18030":           simplifiedchinese.GB18030,
	"big5":              traditionalchinese.Big5,
	"eucjp":             japanese.EUCJP,
	"iso2022jp":         japanese.ISO2022JP,
	"shiftjis":          japanese.ShiftJIS,
	"euckr":             korean.EUCKR,
	"utf16be":           unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	"utf16le":           unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
}

func GetEncoding(charset string) (encoding.Encoding, bool) {
	charset = strings.ToLower(charset)
	enc, ok := encodings[charset]
	return enc, ok
}

func Decode(in []byte, charset string) ([]byte, error) {
	if enc, ok := GetEncoding(charset); ok {
		return enc.NewDecoder().Bytes(in)
	}
	return nil, errors.New("charset not found!")
}
