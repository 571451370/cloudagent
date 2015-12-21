package channel

import (
	"encoding/json"
	"io"

	"github.com/vtolstov/cloudagent/qga"
)

var (
	ErrMessageFormat = &Response{Error: &Error{Code: -1, Desc: "Invalid Message Format"}}
)

// Channel interface provide communication channel with cloudagent
type Channel interface {
	Open() error
	Close() error
	Reset() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Poll() error
}

// Request struct used to parse incoming request
type Request struct {
	Execute string          `json:"execute"`
	RawArgs json.RawMessage `json:"arguments,omitempty"`
	ID      string          `json:"id,omitempty"`
}

// Error struct used to indicate error when processing command
type Error struct {
	Class  string `json:"class,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Bufb64 string `json:"bufb64,omitempty"`
	Code   int    `json:"code,omitempty"`
}

// Response struct used to encode response from command
type Response struct {
	Return interface{} `json:"return,omitempty"`
	Error  *Error      `json:"error,omitempty"`
	ID     string      `json:"id,omitempty"`
}

func Serve(ch Channel) error {
	var err error
	var n int
	var req qga.Request
	var res *qga.Response
	buffer := make([]byte, qga.MaxMessageLength)

	//for {
	n, err = ch.Read(buffer)
	if err != nil && err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}
	if n == 0 {
		return nil
	} else if n > 0 {
		//	break
	}
	//}

	if err = json.Unmarshal(buffer[:n], &req); err != nil {
		res = qga.ErrMessageFormat
	} else {
		res = qga.CmdRun(&req)
	}
	buffer, err = json.Marshal(res)
	buffer = append(buffer, byte('\n'))
	_, err = ch.Write(buffer)
	return err
}
