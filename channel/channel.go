package channel

import (
	"encoding/json"
	"io"

	"github.com/vtolstov/cloudagent/qga"
)

/*
var (
	// ErrMessageFormat message
	ErrMessageFormat = &Response{Error: &Error{Code: -1, Desc: "Invalid Message Format"}}
)
*/

// Channel interface provide communication channel with cloudagent
type Channel interface {
	Open() error
	Close() error
	Reset() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Poll() error
}

// Serve run i/o on channel
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
