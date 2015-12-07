/*
Package guest_fscheck - run check to test filesystem

Example:
        { "execute": "guest-fscheck", "arguments": {
            "path": string // optional, path to store tmp files
          }
        }
*/
package guest_fscheck

import "github.com/vtolstov/cloudagent/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fscheck",
		Func:    fnGuestFscheck,
		Enabled: true,
	})
}

func fnGuestFscheck(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}
	return res
}
