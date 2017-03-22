package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Duration     string `long:"duration" default:"0" description:"Recording duration (0 == infinite)"`
	UseLocalTime bool   `long:"local" description:"Use local time to track duration instead of supplied metadata"`
	UserAgent    string `long:"ua" default:"hlsdownloader" description:"User-Agent for HTTP client"`
	M3U8         string `long:"m3u8" default:"" description:"M3U8 playlist URI to be downloaded"`
	Output       string `long:"output" default:"" description:"Output path"`
	LogLevel     string `long:"log_level" default:"info" description:"log level"`
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func init() {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)

	_, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
}

func init() {
	if level, err := log.ParseLevel(strings.ToLower(opts.LogLevel)); err != nil {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
