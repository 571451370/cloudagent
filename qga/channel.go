package qga

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Channel interface provide communication channel with cloudagent
type Channel interface {
	Open() error
	Close() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Reset() error
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

// WorkerIO read/write bytes from channel
func WorkerIO(ch Channel) error {
	var wg sync.WaitGroup

	//      wg.Add(2)

	// channel read
	go func() {
		defer wg.Done()
		var err error
		buffer := make([]byte, QGA_MAX_MESSAGE_LEN)
		var n int

		for {
			n, err = ch.Read(buffer)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			fmt.Printf("%s\n", buffer[:n])
		}
	}()

	// channel write
	go func() {
		defer wg.Done()
		var err error
		buffer := make([]byte, QGA_MAX_MESSAGE_LEN)
		//var n int

		if _, err = ch.Write(buffer); err == nil {
			//			w <- buffer[:n]
		} else {
			fmt.Printf(err.Error())
		}
	}()

	wg.Wait()

	return nil
}

// WorkerMsg read/write messages
func WorkerMsg() error {
	return nil
}
