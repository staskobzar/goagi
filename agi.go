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

```go
	import (
		"github.com/staskobzar/goagi"
		"log"
	)

	int main() {
		agi, err := goagi.NewAGI()
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
	}
```

*/
func NewAGI() (*AGI, error) {
	in := bufio.NewWriter(os.Stdout)
	out := bufio.NewReader(os.Stdin)

	agi, err := newInterface(bufio.NewReadWriter(out, in))
	if err != nil {
		return nil, err
	}
	return agi, nil
}

// FastAGI defines sturcture of fast AGI server
type FastAGI struct {
	agi  *AGI
	conn net.Conn
}

// Conn returns AGI instance on every Asterisk connection
func (fagi *FastAGI) AGI() *AGI {
	return fagi.agi
}

// Close terminates Fast AGI
func (fagi *FastAGI) Close() error {
	return fagi.conn.Close()
}

// RemoteAddr returns remote connected client host and port as string
func (fagi *FastAGI) RemoteAddr() string {
	return fagi.conn.RemoteAddr().String()
}

/*
NewFastAGI starts listening and serve AGI network calls.

Usage example:

```go
	import (
    	"github.com/staskobzar/goagi"
		"time"
    	"log"
    )

	func main() {
		serve := func (fagi *goagi.FastAGI) {
			agi := fagi.AGI()
			agi.Verbose("New FastAGI session")
			agi.Answer()
			if clid, err := agi.GetVariable("CALLERID"); err == nil {
				log.Printf("CallerID %s\n", clid)
				ag.Varbose("Call from " + clid)
			}
			fagi.Close()
		}
		// listen, serve and reconnect on fail
		for {
			ln, err := net.Listen("tcp", "127.0.0.1:4573")
			if err != nil {
				log.Println("Connection error. Re-try in 3 sec.")
				<-time.After(time.Second * 3)
				continue
			}
			chFagi, chErr := goagi.NewFastAGI(ln)

		Loop:
			for {
				select {
				case fagi := <-chFagi:
					go serve(fagi)
				case err :=<-chErr:
					ln.Close()
					log.Println(err)
					break Loop
				}
			}
		}
	}
```
*/
func NewFastAGI(ln net.Listener) (<-chan *FastAGI, <-chan error) {
	chFagi := make(chan *FastAGI)
	chErr := make(chan error)

	go func(chFagi chan *FastAGI, chErr chan error) {
		defer close(chFagi)
		defer close(chErr)
		for {
			conn, err := ln.Accept()
			if err != nil {
				chErr <- err
				return
			}

			in := bufio.NewWriter(conn)
			out := bufio.NewReader(conn)
			agi, err := newInterface(bufio.NewReadWriter(out, in))
			if err != nil {
				// invalid input data. skip
				conn.Close()
				continue
			}
			fagi := &FastAGI{agi: agi, conn: conn}
			chFagi <- fagi
		}
	}(chFagi, chErr)
	return chFagi, chErr
}
