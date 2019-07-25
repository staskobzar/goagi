// Example of usage NewAGI for Asterisk.
//	Dialplan example:
//		exten => _X.,1,NoOp(Test goagi)
//		 same => n,Answer()
//		 same => n,AGI(/path/to/application)
// Reproduces Asterisk agi-test.agi script
package main

import (
	"fmt"
	"github.com/staskobzar/goagi"
	"log"
)

var (
	tests, fail, pass int
)

func checkResult(err error) {
	tests++
	if err != nil {
		fail++
	} else {
		pass++
	}
}

func main() {
	agi, err := goagi.NewAGI()
	if err != nil {
		log.Fatalln(err)
	}

	verb := func(msg string, args ...interface{}) {
		if err := agi.Verbose(fmt.Sprintf(msg, args...)); err != nil {
			log.Fatalln(err)
		}
	}

	verb("1.  Testing 'sendfile'...")
	_, err = agi.StreamFile("beep", "", 0)
	checkResult(err)

	verb("2.  Testing 'sendtext'...")
	err = agi.SendText("hello world")
	checkResult(err)

	verb("3.  Testing 'sendimage'...")
	err = agi.SendImage("hello world")
	checkResult(err)

	verb("4.  Testing 'saynumber'...")
	err = agi.SayNumber("192837465", "")
	checkResult(err)

	verb("5.  Testing 'waitdtmf'...")
	dig, err := agi.WaitForDigit(1000)
	checkResult(err)
	verb("Digit received: '%s'", dig)

	verb("6.  Testing 'record'...")
	err = agi.RecordFile("testagi", "gsm", "1234", 3000, 0, false, 0)
	checkResult(err)

	verb("6a.  Testing 'record' playback...")
	_, err = agi.StreamFile("testagi", "", 0)
	checkResult(err)

	verb("================== Complete ======================")
	verb("%d tests completed, %d passed, %d failed", tests, pass, fail)
	verb("==================================================")
}
