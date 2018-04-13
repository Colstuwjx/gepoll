package main

import (
	"flag"
	"log"

	"github.com/Colstuwjx/gepoll/epoller"
)

func lineHandler(buf []byte, n int) {
	log.Println("data ", string(buf))
}

func main() {
	deviceFlag := flag.String("device", "/dev/kmsg", "device to use")
	flag.Parse()

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
