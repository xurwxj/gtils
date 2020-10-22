package net

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/viper"
)

// ChunkDownloadEx able to resume downloading.
func ChunkDownloadEx(savePath, fileName, turl, id string, size int64, enableRangeSupport bool, callback func(id string, size, chunkSize, partNum int64)) (string, error) {
	// startTime := time.Now().UTC()
	if savePath != "." && savePath != "./" {
		if stat, err := os.Stat(savePath); os.IsNotExist(err) {
			if err = os.MkdirAll(savePath, os.ModePerm); err != nil {
				return "", err
			}
		} else if !stat.IsDir() {
			return "", err
		}
	}
	tsize, rangeSupport, tfileName, _ := GetSizeNameAndCheckRangeSupport(turl)
	if size == 0 {
		size = tsize
	}
	if fileName == "" {
		fileName = tfileName
	}
	if !rangeSupport && !enableRangeSupport {
		return FileDownload(id, savePath, fileName, turl, callback)
	}
	chunkSize := viper.GetInt64("savePath.limit.downloadMinChunkSize")
	if chunkSize < 1 || chunkSize > 10 {
		chunkSize = 5
	}
	chunkSize = chunkSize * 1024 * 1024
	workerCount := viper.GetInt64("savePath.limit.download")
	if workerCount < 1 {
		workerCount = 10
	}
	partialSize := int64(size / workerCount)
	wc := int64(size / chunkSize)
	if wc < workerCount {
		workerCount = wc
		if workerCount == 0 {
			workerCount = 1
		}
		partialSize = chunkSize
	}
	// fmt.Println("88")
	// fmt.Println("88 size:", size)
	// fmt.Println("88 chunkSize:", chunkSize)
	// fmt.Println("88:", fileName)

	filePath := filepath.Join(savePath, fileName)
	fi := base.CheckFileExistBackInfo(filePath, true)
	if fi != nil && fi.Size() == size {
		return fileName, nil
	}
	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	// handleError(err)
	defer f.Close()
	// fmt.Println("99")

	var start, end int64
	filep := make(map[string]*os.File)
	var fileNames []string
	var workerWG sync.WaitGroup
	var reachedMaxErr chan struct{} = make(chan struct{}, workerCount)
	var trunkFileSize int64
	for num := int64(0); num < workerCount; num++ {
		trunkFileSize = partialSize
		if num == workerCount-1 {
			end = size // last part
			trunkFileSize = size - partialSize*num
		} else {
			end = (num + 1) * partialSize
		}

		singleFile := fmt.Sprintf("%s_%d_%d_downloading", fileName, num, trunkFileSize)
		sfilePath := filepath.Join(savePath, singleFile)
		f, err := os.OpenFile(sfilePath, os.O_CREATE|os.O_RDWR, 0666)
		// f, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return "", err
		}
		filep[sfilePath] = f
		fileNames = append(fileNames, sfilePath)
		var realSize int64
		if fi, err := f.Stat(); err != nil {
			return "", err
		} else {
			realSize = fi.Size()
		}

		start = num*partialSize + realSize

		var worker = Worker{
			URL:           turl,
			ID:            id,
			File:          f,
			TrunkFileSize: trunkFileSize,
			TotalSize:     size,
			SyncWG:        &workerWG,
			Offset:        realSize,
			Parted:        true,
			Callback:      callback,
			ReachedMaxErr: reachedMaxErr,
		}
		workerWG.Add(1)
		go worker.writeRange(num, start, end-1)
	}
	workerWG.Wait()
	// After all workers exit, check if someone has reached max error
	select {
	case <-reachedMaxErr:
		return "", fmt.Errorf("maxWorkerReachedErr:%v", workerCount)
	default:
	}
	if err = mergeFileAndClean(f, filep, fileNames); err != nil {
		return "", err
	}

	return fileName, nil
}

// mergeFileAndClean is called by ChunkDownloadEx to merge files trunks when they are downloaded.
// mainFile is the file to download and files is the map between file trunk names and file handlers.
func mergeFileAndClean(mainFile *os.File, files map[string]*os.File, fileNames []string) error {
	// fmt.Println("merging file and clean")
	for _, Key := range fileNames {
		if v, ok := files[Key]; ok {
			b, err := ioutil.ReadAll(v)
			if err != nil {
				return err
			}
			mainFile.Write(b)
			v.Close()
			os.Remove(Key)
		}
	}
	return nil
}

// Worker struct for download goroutine
type Worker struct {
	URL           string
	ID            string
	File          *os.File
	Count         int64
	SyncWG        *sync.WaitGroup
	TotalSize     int64
	TrunkFileSize int64
	Offset        int64
	Parted        bool
	ErrCount      int64
	Callback      func(id string, size, chunkSize, partNum int64)
	ReachedMaxErr chan struct{}
}

// writeRange calls getRangeBody to get file trunk data and then write them into w.File.
// writeRange can resume from errors.
func (w *Worker) writeRange(partNum int64, start int64, end int64) {
	defer w.SyncWG.Done()
	if w.ErrCount >= 10 {
		w.ReachedMaxErr <- struct{}{}
		return
	}
	var written int64
	body, size, err := w.getRangeBody(start, end)

	if err != nil {
		w.SyncWG.Add(1)
		w.ErrCount++
		go w.writeRange(partNum, start, end)
		return
	}

	if size != end-start+1 {
		w.SyncWG.Add(1)
		w.ErrCount++
		go w.writeRange(partNum, start, end)
		return
	}

	defer body.Close()

	if !w.Parted {
		w.Offset = start
	}
	// make a buffer to keep chunks that are read
	buf := make([]byte, 4*1024)
	for {
		nr, er := body.Read(buf)
		if nr > 0 {
			nw, err := w.File.WriteAt(buf[0:nr], w.Offset)
			if err != nil {
				w.SyncWG.Add(1)
				w.ErrCount++
				go w.writeRange(partNum, start, end)
				break
			}
			if nr != nw {
				w.SyncWG.Add(1)
				w.ErrCount++
				go w.writeRange(partNum, start, end)
				break
			}

			start = int64(nw) + start
			w.Offset += int64(nw)
			if nw > 0 {
				written += int64(nw)
			}
			p := int64(float32(written) / float32(size) * 100)
			if p%20 == 0 {
				// fmt.Println(fmt.Sprintf("Part %d  %d%% write success.", partNum, p), time.Now().UTC())
			}
		}
		// fmt.Println("size: --- ", size)
		// fmt.Println("written: --- ", written)
		if er != nil {
			if er.Error() == "EOF" {
				if size == written {
					// w.SyncWG.Done()
					// return
					// DONE 需要返回下载进度
					if w.Callback != nil {
						go w.Callback(w.ID, w.TotalSize, w.TrunkFileSize, partNum)
					}
				} else {
					w.SyncWG.Add(1)
					w.ErrCount++
					go w.writeRange(partNum, start, end)
					break
				}
				break
			} else {
				w.SyncWG.Add(1)
				w.ErrCount++
				go w.writeRange(partNum, start, end)
				break
			}
		}
	}
}

// getRangeBody issues a GET to w.Url to get file trunk data.
// It returns an io.ReadCloser of data, the data size and an error.
func (w *Worker) getRangeBody(start int64, end int64) (io.ReadCloser, int64, error) {
	// var client http.Client

	client := &http.Client{Transport: &http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}}
	client.Timeout = time.Second * 150
	req, err := http.NewRequest("GET", w.URL, nil)
	// req.Header.Set("cookie", "")
	// log.Printf("Request header: %s\n", req.Header)
	if err != nil {
		return nil, 0, err
	}

	// Set range header
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	size, err := strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	if err != nil {
		return nil, 0, err
	}
	return resp.Body, size, nil
}

// FileDownload issues a GET to url and saves the file to savePath.
// If fileName is empty, FileDownload gets it from url. If still empty, FileDownload will make one by timestamp.
// FileDownload returns the real filename and an error.
func FileDownload(id, savePath, fileName, url string, callback func(id string, size, chunkSize, partNum int64)) (string, error) {
	if savePath != "." && savePath != "./" {
		if stat, err := os.Stat(savePath); os.IsNotExist(err) {
			if err = os.MkdirAll(savePath, os.ModePerm); err != nil {
				return "", err
			}
		} else if !stat.IsDir() {
			return "", err
		}
	}
	var fileSize int64
	if fileName == "" {
		fileSize, _, fileName, _ = GetSizeNameAndCheckRangeSupport(url)
	}
	if fileName == "" {
		fileName = fmt.Sprintf("%v", time.Now().UTC().Unix())
	}
	// fmt.Println("fileSize:", fileSize)
	// fmt.Println("fileName: ", fileName)
	filePath := filepath.Join(savePath, fileName)
	fi := base.CheckFileExistBackInfo(filePath, true)
	if fi != nil && fi.Size() == fileSize {
		return fileName, nil
	}
	client := &http.Client{Transport: &http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// // Write the body to file
	// _, err = io.Copy(out, resp.Body)
	// if err != nil {
	// 	return "", err
	// }
	chunkSize := viper.GetInt64("savePath.limit.downloadMinChunkSize")
	if chunkSize < 1 || chunkSize > 10 {
		chunkSize = 5
	}
	chunkSize = chunkSize * 1024
	buf := make([]byte, chunkSize)
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := out.Write(buf[0:nr])
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
			callback(id, fileSize, int64(nw), 0)
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return fileName, nil
}

// GetSizeNameAndCheckRangeSupport checks whether url supports header "accpet-ranges".
// If so, returns content length.
func GetSizeNameAndCheckRangeSupport(url string) (size int64, rangeSupport bool, fileName string, err error) {
	client := &http.Client{Transport: &http.Transport{
		Dial:              PrintLocalDial,
		DisableKeepAlives: true,
	}}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}
	sheader := make(map[string]string)
	sheader["range"] = "bytes=0-1"
	req = setHeader(req, sheader)
	res, err := client.Do(req)

	if err != nil {
		return
	}
	// fmt.Println(res.Header)
	defer res.Body.Close()
	header := res.Header
	contentRangeHeader := header["Content-Range"]
	if len(contentRangeHeader) > 0 {
		crh := strings.Split(contentRangeHeader[0], "/")[1]
		size, err = strconv.ParseInt(crh, 10, 64)
	}
	if hcd, ok := header["Content-Disposition"]; ok && len(hcd) > 0 {
		hcds := strings.Split(hcd[0], "=")
		if len(hcds) > 1 {
			if filename := hcds[1]; filename != "" {
				fileName = filepath.Base(filename)
			}
		}
	}
	acceptRanges, supported := header["Accept-Ranges"]
	if supported && acceptRanges[0] == "bytes" {
		rangeSupport = true
	}
	return
}
