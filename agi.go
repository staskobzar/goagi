package goagi

import (
	"bufio"
	"net"
	"os"
	"sync"
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
	mu  sync.Mutex
	ch  chan *AGI
	ln  net.Listener
	err error
}

// Conn returns AGI instance on every Asterisk connection
func (fagi *FastAGI) Conn() <-chan *AGI {
	return fagi.ch
}

// Err sets error on Fast AGI processing error
func (fagi *FastAGI) Err() error {
	fagi.mu.Lock()
	err := fagi.err
	fagi.mu.Unlock()
	return err
}

// Close terminates Fast AGI
func (fagi *FastAGI) Close() {
	fagi.ln.Close()
	close(fagi.ch)
}

func (fagi *FastAGI) setErr(err error) {
	fagi.mu.Lock()
	fagi.err = err
	fagi.mu.Unlock()
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
		// listen, serve and reconnect on fail
		for {
			fagi, err := goagi.NewFastAGI(":8000")
			if err != nil {
				log.Println("Connection error. Re-try in 3 sec.")
				<-time.After(time.Second * 3)
				continue
			}

			for agi := range fagi.Conn() {
				agi.Verbose("New FastAGI session")
				agi.Answer()
				if clid, err := agi.GetVariable("CALLERID"); err == nil {
					log.Printf("CallerID %s\n", clid)
					ag.Varbose("Call from " + clid)
				}
			}
			if agi.Err() != nil {
				fmt.Printf("Error: %s\n", agi.Err())
			}
		}
	}
```
*/
func NewFastAGI(listenAddr string) (*FastAGI, error) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	ch := make(chan *AGI)
	fagi := &FastAGI{ln: ln, ch: ch}

	go func(fastagi *FastAGI) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fastagi.setErr(err)
				return
			}

			in := bufio.NewWriter(conn)
			out := bufio.NewReader(conn)
			agi, err := newInterface(bufio.NewReadWriter(out, in))
			if err != nil {
				fastagi.setErr(err)
				return
			}
			fastagi.mu.Lock()
			if fastagi.err == nil {
				fastagi.ch <- agi
			}
			fastagi.mu.Unlock()
		}
	}(fagi)
	return fagi, nil
}
