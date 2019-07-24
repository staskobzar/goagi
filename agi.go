package goagi

import (
	"bufio"
	"net"
	"os"
)

/*
NewAGI creates and returns AGI object.
Parses AGI arguments and set ready for communication.

Usage example:

	import (
		"github.com/staskobzar/goagi"
		"log"
		"os"
	)

	int main() {
		agi, err := goagi.NewAgi()
		if err != nil {
			log.Fatalln(err)
		}
		agi.Verbose("New AGI session.")
		if err := agi.SetMusic("on", "jazz"); err != nil {
			log.Fatalln(err)
		}

		clid, err := agi.GetVariable("CALLERID")
		if err != nil {
			log.Fatalln(err)
		}
		agi.Verbose("Call from " + clid)
		if err := agi.SetMusic("off"); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}

*/
func NewAGI() (*AGI, error) {
	in := bufio.NewWriter(os.Stdin)
	out := bufio.NewReader(os.Stdout)

	agi, err := newInterface(bufio.NewReadWriter(out, in))
	if err != nil {
		return nil, err
	}
	return agi, nil
}

type callbackFunc func(agi *AGI)

/*
NewFastAGI starts listening and serve AGI network calls

Usage example:

	import (
		"github.com/staskobzar/goagi"
		"log"
	)

	// listen and serve
	err := NewFastAGI(":8000", myAgiProc)
	if err != nil {
		log.Fatalln(err)
	}

	// callback function
	func myAgiProc(agi *AGI) {
		agi.Verbose("New AGI session.")
		agi.Answer()
		clid agi.GetVariable("CALLERID")
		clid, err := agi.GetVariable("CALLERID")
		if err != nil {
			log.Fatalln(err)
		}
		agi.Verbose("Call from " + clid)
	}
*/
func NewFastAGI(listenAddr string, callback callbackFunc) error {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		in := bufio.NewWriter(conn)
		out := bufio.NewReader(conn)
		agi, err := newInterface(bufio.NewReadWriter(out, in))
		if err != nil {
			return err
		}
		go callback(agi)
	}
}
