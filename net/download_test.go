package net

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestChunkDownloadCtxEx chunkDownloadCtx
// 利用 cancel 取消当前下载
func TestChunkDownloadCtxEx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	url:="url" // down file
	go ChunkDownloadEx("./","test",url,"",0,false,ctx,CallFunc)
	tic:=time.NewTicker(time.Second*10)
	for {
		select {
		case <-tic.C:
			t.Log("ctx done")
			cancel()
			time.Sleep(2*time.Second)
			return
		}
	}
}

// TestChunkDownloadEx test chunkDownload
func TestChunkDownloadEx(t *testing.T) {
	url:="url" // down file
	ChunkDownloadEx("./","test",url,"",0,false,context.Background(),CallFunc)

}

func CallFunc(id string,size,chunkSize,partNum int64){
	fmt.Println(id,size,chunkSize,partNum)
}