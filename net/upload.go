package net

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/viper"
)

// MUploadToRemote issues a POST to surl to upload file specified by filePath.
// category, bucket and env are used to make up upload parameters.
// It returns DfsId and an error.
func MUploadToRemote(surl, category, bucket, filePath, id string, header map[string]string, callback func(id string, size, chunkSize, partNum int64)) (string, error) {

	// DONE change to concurrency solid
	// startTime := time.Now().UTC()
	BufferSize := viper.GetInt("savePath.limit.uploadMinChunkSize")
	if BufferSize < 5242880 {
		BufferSize = 5 * 1024 * 1024
	}
	file, err := os.Open(filePath)
	if err != nil {
		// fmt.Println("file open err: ", err)
		return "", err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		// fmt.Println("get file info err: ", err)
		return "", err
	}

	filesize := int(fileinfo.Size())

	hasher := md5.New()
	hasher.Write([]byte(category + bucket + filePath + fmt.Sprintf("%d", filesize) + fmt.Sprintf("%v", fileinfo.ModTime())))
	identifier := hex.EncodeToString(hasher.Sum(nil))

	fileName := fileinfo.Name()
	// Number of go routines we need to spawn.
	uLimit := viper.GetInt("savePath.limit.upload")
	if uLimit == 0 {
		uLimit = 10
	}
	concurrency := filesize / BufferSize
	if concurrency > uLimit {
		BufferSize = filesize / uLimit
		concurrency = filesize / BufferSize
	}
	// fmt.Println("task: ", fileName, " size: ", filesize, " uLimit: ", uLimit, " concurrency: ", concurrency, " BufferSize: ", BufferSize)
	if filesize <= BufferSize {
		BufferSize = filesize
		concurrency = 1
	}

	// buffer sizes that each of the go routine below should use. ReadAt
	// returns an error if the buffer size is larger than the bytes returned
	// from the file.
	chunksizes := make([]chunk, concurrency)

	// All buffer sizes are the same in the normal case. Offsets depend on the
	// index. Second go routine should start at 100, for example, given our
	// buffer size of 100.
	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

	// check for any left over bytes. Add the residual number of bytes as the
	// the last chunk size.
	if remainder := filesize % BufferSize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
		concurrency++
		chunksizes = append(chunksizes, c)
	}
	// fmt.Println("task: ", fileName, " size: ", filesize, " uLimit: ", uLimit, " concurrency: ", concurrency, " BufferSize: ", BufferSize)

	dfsID := ""

	var wg sync.WaitGroup
	// wg.Add(concurrency)
	// fmt.Println(fileName, " start with ", concurrency, " uploads from ", startTime)
	// okPart := 0
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		// go doChunk(chunksizes, i, filesize, concurrency, BufferSize, file, &wg, &dfsID, category, bucket, fileName, filePath, identifier, surl, id, callback, header)
		go func(chunksizes []chunk, i, filesize int) {

			chunk := chunksizes[i]
			buffer := make([]byte, chunk.bufsize)
			bytesread, err := file.ReadAt(buffer, chunk.offset)

			if err != nil {
				// fmt.Println("file read err: ", err)
				return
			}
			if bytesread == chunk.bufsize {
				currentChunk := i + 1
				var pm chunkObj
				pm.Category = category
				pm.Bucket = bucket
				pm.Filename = fileName
				pm.TotalSize = int64(filesize)
				pm.RelativePath = filePath
				pm.ChunkSize = int64(BufferSize)
				pm.ChunkNumber = currentChunk
				pm.TotalChunks = concurrency
				pm.CurrentChunkSize = int64(chunk.bufsize)
				pm.Identifier = identifier
				pmv, err := base.Values(pm)
				if err != nil {
					// fmt.Println("Values in MUploadToRemote err: ", err)
				}
				mUpGetURL := fmt.Sprintf("%s?%s", surl, pmv.Encode())
				if strings.Index(surl, "?") >= 0 {
					mUpGetURL = fmt.Sprintf("%s&%s", surl, pmv.Encode())
				}
				resCode, rmu := mUpGet(mUpGetURL, header)
				if base.FindInInt64Slice([]int64{400, 404}, resCode) {
					rssd := mUpPost(surl, pm, buffer, header)
					// DONE 需要返回上传进度
					if rssd.Result.DfsID != "" && rssd.Result.DfsID != "OK" {
						dfsID = rssd.Result.DfsID
						// okPart = okPart + 1
						// fmt.Println("rssd.Result.DfsID: ", rssd.Result.DfsID, " dfsID: ", *dfsID, " with pm: ", pm, "rs: ", rssd, " on ", time.Now().UTC())
						if callback != nil {
							callback(id, pm.TotalSize, pm.CurrentChunkSize, int64(pm.ChunkNumber))
						}
					}
					if rssd.Result.DfsID == "OK" {
						// okPart = okPart + 1
						// fmt.Println("finish ", pm, " rs: ", rssd, " on ", time.Now().UTC())
						if callback != nil {
							callback(id, pm.TotalSize, pm.CurrentChunkSize, int64(pm.ChunkNumber))
						}
					}
				} else if resCode == 200 {
					// DONE 需要返回上传进度
					if rmu.Result.DfsID != "" && rmu.Result.DfsID != "OK" {
						dfsID = rmu.Result.DfsID
						// okPart = okPart + 1
						// fmt.Println("rssd.Result.DfsID: ", rssd.Result.DfsID, " dfsID: ", *dfsID, " with pm: ", pm, "rs: ", rssd, " on ", time.Now().UTC())
						if callback != nil {
							callback(id, pm.TotalSize, pm.CurrentChunkSize, int64(pm.ChunkNumber))
						}
					}
					if rmu.Result.DfsID == "OK" {
						// okPart = okPart + 1
						// fmt.Println("finish ", pm, " rs: ", rssd, " on ", time.Now().UTC())
						if callback != nil {
							callback(id, pm.TotalSize, pm.CurrentChunkSize, int64(pm.ChunkNumber))
						}
					}
				}
			}
			wg.Done()
		}(chunksizes, i, filesize)
	}

	wg.Wait()
	// fmt.Println(filePath, " upload to: ", dfsID, " upload to: ", dfsID, " done from ", startTime, " to ", time.Now().UTC())
	if dfsID == "" {
		// fmt.Println(filePath, " upload to: ", dfsID, " done from ", startTime, " to ", time.Now().UTC())
		// } else {
		// fmt.Println("need to reupload for file: ", filePath, " when concurrency: ", concurrency, " with dfsID: ", dfsID)
		// return MUploadToRemote(surl, category, bucket, filePath, id, header, callback)
	}
	return dfsID, nil
}

func mUpPost(url string, p chunkObj, chunkData []byte, header map[string]string) resMUp {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	for k, v := range base.StructToFormMap(p) {
		bodyWriter.WriteField(k, v.(string))
	}
	fileWriter, err := bodyWriter.CreateFormFile("file", p.Filename)
	if err != nil {
		// fmt.Println("body CreateFormFile err: ", err)
		return mUpPost(url, p, chunkData, header)
	}
	io.Copy(fileWriter, bytes.NewReader(chunkData))
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	var req *http.Request
	var res *http.Response
	// var err error
	// client := &http.Client{}

	client := &http.Client{Transport: &http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}}
	// client.Timeout = time.Second * 150
	client.Timeout = time.Second * 150000000
	// fmt.Println(formBody)
	req, err = http.NewRequest("POST", url, bodyBuffer)
	req = setHeader(req, header)
	req.Header.Set("Content-Type", contentType)
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Keep-Alive")

	if err != nil {
		// fmt.Println("req err: ", err)
		return mUpPost(url, p, chunkData, header)
	}
	res, err = client.Do(req)
	if err != nil {
		// fmt.Println("req do err: ", err)
		return mUpPost(url, p, chunkData, header)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// fmt.Println("read res body not 200: ", res.StatusCode)
		return mUpPost(url, p, chunkData, header)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// fmt.Println("read res body err: ", err)
		return mUpPost(url, p, chunkData, header)
	}
	// fmt.Println(p.Filename, " put data upload res data: ", string(data))
	var rmu resMUp
	if string(data) == "OK" {
		rmu.Result.DfsID = "OK"
		return rmu
	}
	if err := json.Unmarshal(data, &rmu); err == nil {
		return rmu
	}
	if string(data) == "NotExist" {
		// fmt.Println("data not exist nned to upload ")
		return mUpPost(url, p, chunkData, header)
	}
	return mUpPost(url, p, chunkData, header)
}

func mUpGet(url string, header map[string]string) (int64, resMUp) {
	var req *http.Request
	var res *http.Response
	var err error
	client := &http.Client{Transport: &http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}}
	client.Timeout = time.Second * 150
	// fmt.Println(formBody)
	req, err = http.NewRequest("GET", url, nil)
	req = setHeader(req, header)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Keep-Alive")
	var rmu resMUp

	if err != nil {
		return 500, rmu
	}
	res, err = client.Do(req)
	if err != nil {
		return 500, rmu
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return 460, rmu
		}
		if string(data) == "OK" {
			rmu.Result.DfsID = "OK"
			return 200, rmu
		}
		// fmt.Println(url, ": Get data upload res data: ", string(data))
		// fmt.Println("res data: ", string(data))
		if err := json.Unmarshal(data, &rmu); err == nil {
			return 200, rmu
		}
	}
	return int64(res.StatusCode), rmu
}

type chunk struct {
	bufsize int
	offset  int64
}

type resURL struct {
	Result `json:"result"`
	Status string `json:"status"`
}

// Result file and url info from server
type Result struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
}

type chunkObj struct {
	Category         string `form:"category" json:"category"`
	SubCate          string `form:"subcate" json:"subcate"`
	Bucket           string `form:"bucket" json:"bucket"`
	ChunkNumber      int    `form:"chunkNumber" on:"chunkNumber"`
	Identifier       string `form:"identifier" json:"identifier"`
	Filename         string `form:"filename" json:"filename"`
	RelativePath     string `form:"relativePath" json:"relativePath"`
	CurrentChunkSize int64  `form:"currentChunkSize" json:"currentChunkSize"`
	ChunkSize        int64  `form:"chunkSize" json:"chunkSize"`
	TotalSize        int64  `form:"totalSize" json:"totalSize"`
	TotalChunks      int    `form:"totalChunks" json:"totalChunks"`
	DownValidTo      int64  `form:"downValidTo" json:"downValidTo"`
	// FileUp           *multipart.FileHeader `form:"file" json:"file"`
}

type resMUp struct {
	Result attachment `json:"result"`
	Status string     `json:"status"`
}

type attachment struct {
	ID       string    `json:"id"`
	CreateOn time.Time `json:"create_on"`

	Name        string `json:"name"`
	Extension   string `json:"extension"`
	FileLength  int64  `json:"file_length"`
	DfsID       string `json:"dfs_id"`
	DownURL     string `json:"downURL"`
	Catetory    string `json:"catetory"`
	Endpoint    string `json:"endpoint"`
	Bucket      string `json:"bucket"`
	IsPub       string `json:"is_pub"`
	ContentType string `json:"content_type"`
}
