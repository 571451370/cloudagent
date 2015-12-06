package qga

import "fmt"

// Command struct contains supported commands
type Command struct {
	Enabled   bool                     `json:"enabled"`          // flag to enable command
	Name      string                   `json:"name"`             // command name
	Func      func(*Request) *Response `json:"-"`                // command execution function
	Returns   bool                     `json:"success-response"` // flag for command returned value on success
	Arguments bool                     `json:"-"`                // flag for comand that it needs arguments
}

var (
	commands = make(map[string]*Command)
)

// RegisterCommand registers command to process inside worker
func RegisterCommand(cmd *Command) error {
	if _, ok := commands[cmd.Name]; ok {
		return fmt.Errorf("command %s already registered", cmd.Name)
	}
	commands[cmd.Name] = cmd
	return nil
}

// ListCommands returns commands
func ListCommands() []*Command {
	var commands []*Command
	for _, cmd := range commands {
		commands = append(commands, cmd)
	}
	return commands
}

// CmdRun executes command
func CmdRun(req *Request) *Response {
	if req == nil || req.Execute == "" {
		return &Response{Error: &Error{Class: "CommandNotFound", Desc: fmt.Sprintf("invalid command")}}
	}

	if cmd, ok := commands[req.Execute]; ok && cmd.Func != nil {
		if cmd.Arguments && req.RawArgs == nil {
			return &Response{Error: &Error{Class: "CommandNotFound", Desc: fmt.Sprintf("invalid request for %s", req.Execute)}}
		}
		res := cmd.Func(req)
		if cmd.Returns || res.Error != nil {
			return res
		}
		return &Response{Return: struct{}{}}

	}

	return &Response{Error: &Error{Class: "CommandNotFound", Desc: fmt.Sprintf("command %s not found", req.Execute)}}
}
