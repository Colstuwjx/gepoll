package main

import (
	"flag"
	"log"

	"gepoll/epoller"
)

var deviceFlag = flag.String("device", "/dev/kmsg", "device to use")

func lineHandler(buf []byte, n int) {
	log.Println("data ", string(buf))
}

func main() {
	ep := epoller.NewEpoller(lineHandler)
	err := ep.Open(*deviceFlag)
	if err != nil {
		panic(err)
	}

	defer ep.Close()
	err = ep.Dispatch()
	if err != nil {
		panic(err)
	}
}
