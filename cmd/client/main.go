package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

const (
	limit       = 60 * 1024
	defaultAddr = "0.0.0.0:8081"
)

func main() {
	var (
		addr, text string
		count      int
	)
	flag.StringVar(&addr, "c", defaultAddr, "server listening host/port")
	flag.IntVar(&count, "l", 1, "number to repeat")
	flag.Parse()

	fmt.Fscan(os.Stdin, &text)
	if len(text) > limit {
		fmt.Println(fmt.Sprintf(
			"string is too big. Max allowed length is %d, current: %d", limit, len(text),
		))
		return
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go send(addr, &text)
	}
	wg.Wait()
}

func send(addr string, text *string) {
	defer wg.Done()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("could not connect to server: ", err)
		return
	}
	fmt.Fprintf(conn, *text+"\n")
	fmt.Println("send string with length: ", len(*text))
}
