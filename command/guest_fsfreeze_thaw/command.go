// +build !linux

/*
Package guest_fsfreeze_thaw - run unfreeze on all mounted file systems

Example:
        { "execute": "guest-fsfreeze-thaw", "arguments": {} }
*/
package guest_fsfreeze_thaw

import "github.com/vtolstov/cloudagent/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fsfreeze-thaw",
		Func:    fnGuestFsfreezeThaw,
		Enabled: false,
		Returns: true,
	})
}

func fnGuestFsfreezeThaw(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}
	return res
}
