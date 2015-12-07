/*
Package guest_fsfreeze_status - get status of file systems

Example:
        { "execute": "guest-fsfreeze-status", "arguments": {} }
*/
package guest_fsfreeze_status

import "github.com/vtolstov/cloudagent/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fsfreeze-status",
		Func:    fnGuestFsfreezeStatus,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestFsfreezeStatus(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}

	resData := "thawed"

	if _, ok := qga.StoreGet("guest-fsfreeze", "paths"); ok {
		resData = "frozen"
	}

	res.Return = resData
	return res
}
