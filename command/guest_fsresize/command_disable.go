// +build openbsd windows

/*
Package guest_fsresize - run resize file system

Example:
        { "execute": "guest-fsresize", "arguments": {
            "path": string // optional, mounted filesystem path
          }
        }
*/
package guest_fsresize

import "github.com/vtolstov/cloudagent/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fsresize",
		Func:    fnGuestFsresize,
		Enabled: false,
	})
}

func fnGuestFsresize(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}
	return res
}
