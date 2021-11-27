package main

import (
	"flag"
	"github.com/golang/glog"
	"httpServer/connector"
	"os"
)

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
	flag.Set("v", "2")
	// This is wa
	flag.Parse()
}

func main() {
	glog.Infof("Starting Server")
	server := connector.NewSever()
	server.Run("8080")
}
