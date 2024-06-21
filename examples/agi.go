// Example of usage NewAGI for Asterisk.
//
//	Dialplan example:
//		exten => _X.,1,NoOp(Test goagi)
//		 same => n,Answer()
//		 same => n,AGI(/path/to/application)
//
// Reproduces Asterisk agi-test.agi script
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/staskobzar/goagi"
)

var tests, fail, pass int

func checkResult(err error, resp goagi.Response) {
	tests++
	if err != nil {
		fail++
	} else {
		pass++
	}
	log.Printf("Response code:   %d", resp.Code())
	log.Printf("Response value:  %s", resp.Value())
	log.Printf("Response data:   %s", resp.Data())
	log.Printf("Response endpos: %s", resp.EndPos())
}

func main() { //nolint:typecheck
	agi, err := goagi.New(os.Stdin, os.Stdout, nil)
	if err != nil {
		log.Fatalln(err)
	}

	verb := func(msg string, args ...interface{}) {
		if _, err := agi.Verbose(fmt.Sprintf(msg, args...)); err != nil {
			log.Fatalln(err)
		}
	}

	verb("1.  Testing 'sendfile'...")
	resp, err := agi.StreamFile("beep", "", 0)
	checkResult(err, resp)

	verb("2.  Testing 'sendtext'...")
	resp, err = agi.SendText("hello world")
	checkResult(err, resp)

	verb("3.  Testing 'sendimage'...")
	resp, err = agi.SendImage("hello world")
	checkResult(err, resp)

	verb("4.  Testing 'saynumber'...")
	resp, err = agi.SayNumber("192837465", "")
	checkResult(err, resp)

	verb("5.  Testing 'waitdtmf'...")
	resp, err = agi.WaitForDigit(1000)
	checkResult(err, resp)
	verb("Digit received: '%s'", resp.Value())

	verb("6.  Testing 'record'...")
	resp, err = agi.RecordFile("testagi", "gsm", "1234", 3000, 0, false, 0)
	checkResult(err, resp)

	verb("6a.  Testing 'record' playback...")
	resp, err = agi.StreamFile("testagi", "", 0)
	checkResult(err, resp)

	verb("================== Complete ======================")
	verb("%d tests completed, %d passed, %d failed", tests, pass, fail)
	verb("==================================================")
}
