package net

import (
	"context"
	"fmt"
	"github.com/xurwxj/gtils/base"
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


type TestDown struct {
 	Sing chan Msg

}

type Msg struct {
	Stop bool
	Start bool
	Kill bool
}


func TestChunkDownloadCtxEx1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	url:="" // down file
	var td TestDown
	td.Sing = make(chan Msg,1)
	//go ChunkDownloadExStart("./tmp/order/","test2.dmg",url,"",0,false,ctx,CallFunc)
	//return
	go td.Job()
	//tic:=time.NewTicker(time.Second*20)
	for {
		select {
		case s:=<-td.Sing:
			if s.Stop{
				t.Log("ctx stop")
				cancel()
				time.Sleep(2*time.Second)
			}
			if s.Kill{
				t.Log("ctx kill")
				cancel()
				time.Sleep(2*time.Second)
				return
			}
			if s.Start{
				t.Log("ctx start")
				ctx, cancel = context.WithCancel(context.Background())
				var td TestDown
				td.Sing = make(chan Msg,1)

				go ChunkDownloadEx("./tmp/order/","test2.dmg",url,"",0,false,ctx,CallFunc)
			}
		}
	}
}

func TestCheckFileExistBackInfo(t *testing.T) {
	file:=base.CheckFileExistBackInfo("./tmp/order/test2.dmg",true)
	if file !=nil && file.Size()>0{
		t.Log("ok...")
	}

}

func (t TestDown)Job()  {
	//time.Sleep(time.Second*3)
	//t.Sing<-Msg{Stop: true}
	time.Sleep(time.Second*2)
	t.Sing<-Msg{Start: true}
	//time.Sleep(time.Second*2)
	//t.Sing<-Msg{Kill: true}
}