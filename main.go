package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/gtils/sys"
)

func main() {
	InitLog()
	s := sys.GetOsInfo(Log)
	fmt.Println("s: ", s)
}

var Log *zerolog.Logger

func InitLog() {
	fmt.Println("initializing logger...")
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	var c base.Config
	c.EncodeLogsAsJson = true
	c.ConsoleLoggingEnabled = true
	c.FileLoggingEnabled = false
	c.LocalTime = true
	Log = base.Configure(c)
	Log.Info().Msg("init log done")
}
