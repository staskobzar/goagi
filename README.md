# goagi: Golang library to build agi/fastagi applications

Simple library that helps to build AGI sctipts or FastAGI servers with Go.
```go
import "github.com/staskobzar/goagi"
```

## Usage AGI
Main method for AGI scripts is ```goagi.NewAGI()```.

Example:
```go
    import (
    	"github.com/staskobzar/goagi"
    	"log"
    )

    // callback function
    // net accepted connection will be closed with the callback function returns
    func myAgiProc(agi *AGI) {
	agi.Verbose("New AGI session.")
	agi.Answer()
	clid, err := agi.GetVariable("CALLERID")
	if err != nil {
		log.Fatalln(err)
	}
	agi.Verbose("Call from " + clid)
    }

    func main() {
	// listen and serve
	err := NewFastAGI(":8000", myAgiProc)
	if err != nil {
	    log.Fatalln(err)
	}
    }
```

## Usage FastAGI
Main method to build FastAGI scripts is ```goagi.NewFastAGI(addr, callback)```.

Example:
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

See more examples in ```/examples``` folder.

## goagi API

### Errors
```go
var (
	// EInvalResp error returns when AGI response does not match pattern
	EInvalResp = errorNew("Invalid AGI response")
	// EHangUp error when HANGUP signal received
	EHangUp = errorNew("HANGUP")
	// EInvalEnv error returned when AGI environment header is not valid
	EInvalEnv = errorNew("Invalid AGI env variable")
)
```

### Functions
```go
func NewFastAGI(listenAddr string, callback callbackFunc) error
```
> NewFastAGI starts listening and serve AGI network calls

```go
func NewAGI() (*AGI, error)
```
> NewAGI creates and returns AGI object. Parses AGI arguments and set ready
> for communication.


## AGI struct API

```go
func (agi *AGI) Answer() (bool, error)
```
    Answer executes AGI command "ANSWER" Answers channel if not already in
    answer state.

```go
func (agi *AGI) AsyncAGIBreak() (bool, error)
```
    AsyncAGIBreak Interrupts Async AGI

    Interrupts expected flow of Async AGI commands and returns control

    to previous source (typically, the PBX dialplan).

```go
func (agi *AGI) ChannelStatus(channel string) (int, error)
```
    ChannelStatus returns status of the connected channel.

    If no channel name is given (empty line) then returns the status of the
    current channel.

    Return values:

    0 - Channel is down and available.
    1 - Channel is down, but reserved.
    2 - Channel is off hook.
    3 - Digits (or equivalent) have been dialed.
    4 - Line is ringing.
    5 - Remote end is ringing.
    6 - Line is up.
    7 - Line is busy.

```go
func (agi *AGI) Command(cmd string) (code int, result int, respStr string, err error)
```
    Command sends command as string to the AGI and returns response valus with
    text response

```go
func (agi *AGI) ControlStreamFile(filename, digits string, args ...interface{}) (int32, error)
```
    ControlStreamFile sends audio file on channel and allows the listener to
    control the stream.

    Send the given file, allowing playback to be controlled by the given digits, if any.

    Use double quotes for the digits if you wish none to be permitted. If
    offsetms is provided then the audio will seek to offsetms before play
    starts.

    Returns 0 if playback completes without a digit being pressed, or the ASCII numerical

    value of the digit if one was pressed, or -1 on error or if the channel was
    disconnected.

    Returns the position where playback was terminated as endpos.
    Example:
    agi.ControlStreamFile("prompt_en", "19", "3000", "#", "0", "#", "1600")
    agi.ControlStreamFile("prompt_en", "")
    agi.ControlStreamFile("prompt_en", "19", "", "", "", "#", "1600")

```go
func (agi *AGI) DatabaseDel(family, key string) (bool, error)
```
    DatabaseDel deletes an entry in the Asterisk database for a given family and
    key.

    Returns status and error if fails.

```go
func (agi *AGI) DatabaseDelTree(family, keytree string) (bool, error)
```
    DatabaseDelTree deletes a family or specific keytree within a family in the
    Asterisk database.

```go
func (agi *AGI) DatabaseGet(family, key string) (string, error)
```
    DatabaseGet Retrieves an entry in the Asterisk database for a given family
    and key.

    Returns value as string or error if failed or value not set

```go
func (agi *AGI) DatabasePut(family, key, val string) (bool, error)
```
    DatabasePut adds or updates an entry in the Asterisk database for a given
    family, key, and value.

```go
func (agi *AGI) Env(key string) string
```
    Env returns AGI environment variable by key

```go
func (agi *AGI) EnvArgs() []string
```
    EnvArgs returns list of environment arguments

```go
func (agi *AGI) Exec(app, opts string) (int, error)
```
    Exec executes application with given options.

```go
func (agi *AGI) GetData(file string, args ...interface{}) (digit string, timeout bool, err error)
```
    GetData Stream the given file, and receive DTMF data.

```go
func (agi *AGI) GetFullVariable(name string, channel ...string) (string, error)
```
    GetFullVariable evaluates a channel expression

```go
func (agi *AGI) GetOption(filename, digits string, timeout int32) (int, int32, error)
```
    GetOption Stream file, prompt for DTMF, with timeout.

    Behaves similar to STREAM FILE but used with a timeout option.
    Returns digit pressed, offset and error

```go
func (agi *AGI) GetVariable(name string) (string, error)
```
    GetVariable Gets a channel variable.

```go
func (agi *AGI) Hangup(channel ...string) (bool, error)
```
    Hangup a channel.

    Hangs up the specified channel. If no channel name is given, hangs up the current channel

```go
func (agi *AGI) Noop() error
```
    Noop Does nothing.

```go
func (agi *AGI) ReceiveChar(timeout int32) (int, error)
```
    ReceiveChar Receives one character from channels supporting it.

    Most channels do not support the reception of text. Returns the decimal value of

    the character if one is received, or 0 if the channel does not support text
    reception.

    timeout - The maximum time to wait for input in milliseconds, or 0 for infinite. Most channels
    Returns -1 only on error/hangup.

```go
func (agi *AGI) ReceiveText(timeout int32) (string, error)
```
    ReceiveText Receives text from channels supporting it.

    timeout - The timeout to be the maximum time to wait for input in milliseconds, or 0 for infinite.

```go
func (agi *AGI) RecordFile(file, format, escDigits string,
```
	timeout, offset int, beep bool, silence int) error
    RecordFile Record to a file until a given dtmf digit in the sequence is
    received. The format will specify what kind of file will be recorded. The
    timeout is the maximum record time in milliseconds, or -1 for no timeout.

    offset samples is optional, and, if provided, will seek to the offset without

    exceeding the end of the file.

    beep causes Asterisk to play a beep to the channel that is about to be recorded.
    silence is the number of seconds of silence allowed before the function returns

    despite the lack of dtmf digits or reaching timeout.

    silence is the number of seconds of silence that are permitted before the

    recording is terminated, regardless of the escape_digits or timeout
    arguments

```go
func (agi *AGI) SayAlpha(number, escDigits string) error
```
    SayAlpha says a given character string, returning early if any of the given
    DTMF digits are received on the channel.

```go
func (agi *AGI) SayDate(date, escDigits string) error
```
    SayDate say a given date, returning early if any of the given DTMF digits
    are received on the channel

```go
func (agi *AGI) SayDatetime(time, escDigits, format, timezone string) error
```
    SayDatetime say a given time, returning early if any of the given DTMF
    digits are received on the channel

```go
func (agi *AGI) SayDigits(number, escDigits string) error
```
    SayDigits say a given digit string, returning early if any of the given DTMF
    digits are received on the channel

func (agi *AGI) SayNumber(number, escDigits string) error
    SayNumber say a given digit string, returning early if any of the given DTMF
    digits are received on the channel

func (agi *AGI) SayPhonetic(str, escDigits string) error
    SayPhonetic say a given character string with phonetics, returning early if
    any of the given DTMF digits are received on the channel

func (agi *AGI) SayTime(time, escDigits string) error
    SayTime say a given time, returning early if any of the given DTMF digits
    are received on the channel

func (agi *AGI) SendImage(image string) error
    SendImage Sends the given image on a channel. Most channels do not support
    the transmission of images.

func (agi *AGI) SendText(text string) error
    SendText Sends the given text on a channel. Most channels do not support the
    transmission of text.

func (agi *AGI) SetAutoHangup(seconds int) error
    SetAutoHangup Cause the channel to automatically hangup at time seconds in
    the future. Setting to 0 will cause the autohangup feature to be disabled on
    this channel.

func (agi *AGI) SetCallerid(clid string) error
    SetCallerid Changes the callerid of the current channel.

func (agi *AGI) SetContext(ctx string) error
    SetContext Sets the context for continuation upon exiting the application.

func (agi *AGI) SetExtension(ext string) error
    SetExtension Changes the extension for continuation upon exiting the
    application.

func (agi *AGI) SetMusic(opt string, class ...string) error
    SetMusic Enables/Disables the music on hold generator. If class is not
    specified, then the default music on hold class will be used.

    Parameters: opt is "on" or "off", and music class as string

func (agi *AGI) SetPriority(priority string) error
    SetPriority Changes the priority for continuation upon exiting the
    application. The priority must be a valid priority or label.

func (agi *AGI) SetVariable(name, value string) error
    SetVariable Sets a variable to the current channel.

func (agi *AGI) StreamFile(file, escDigits string, offset int) (int, error)
    StreamFile Send the given file, allowing playback to be interrupted by the
    given digits, if any.

func (agi *AGI) TDDMode(mode string) error
    TDDMode Enable/Disable TDD transmission/reception on a channel.

    Modes: on, off, mate, tdd

func (agi *AGI) Verbose(msg string, level ...int) error
    Verbose Sends message to the console via verbose message system. level is
    the verbose level (1-4)

func (agi *AGI) WaitForDigit(timeout int) (string, error)
    WaitForDigit Waits up to timeout *milliseconds* for channel to receive a
    DTMF digit. Use -1 for the timeout value if you desire the call to block
    indefinitely.

    Return digit pressed as string or error

