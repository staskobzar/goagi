package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/staskobzar/goagi"
)

func serve(conn net.Conn) {
	defer conn.Close()

	dbg := log.New(os.Stdout, "fastagi example: ", log.Lmicroseconds)
	agi, err := goagi.New(conn, conn, dbg)
	if err != nil {
		panic(err)
	}

	resp, err := agi.Verbose("Hello World!")
	if err != nil {
		dbg.Printf("Failed verbose command: %s", err)
		return
	}
	dbg.Printf("Verbose response code: %d", resp.Code())

	resp, _ = agi.GetVariable("CHANNEL")
	dbg.Printf("Asterisk channel variable: %s", resp.Value())

	resp, err = agi.GetData("welcome", 0, 2)
	if err != nil {
		panic(err)
	}

	dbg.Printf("Get Data response code: %d", resp.Code())
	dbg.Printf("Get Data response result: %d", resp.Result())
	dbg.Printf("Get Data response value: %s", resp.Value())
	dbg.Printf("Get Data response data: %s", resp.Data())

	resp, err = agi.GetVariable("CDR(duration)")
	if err != nil {
		panic(err)
	}

	dbg.Printf("CDR duration value: %s", resp.Value())

	agi.Verbose("Goodbye!")
}

func main() {
	fmt.Println("[x] Starting FastAGI script")
	ln, err := net.Listen("tcp", "127.0.0.1:4575")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go serve(conn)
	}
}
