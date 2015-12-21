// +build !linux

/*
Package guest_fsfreeze_freeze - run fsfreeze on all mounted file systems

Example:
        { "execute": "guest-fsfreeze-freeze", "arguments": {} }
*/
package guest_fsfreeze_freeze

import "github.com/vtolstov/cloudagent/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fsfreeze-freeze",
		Func:    fnGuestFsfreezeFreeze,
		Enabled: false,
		Returns: true,
	})
}

func fnGuestFsfreezeFreeze(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}
	return res
}
