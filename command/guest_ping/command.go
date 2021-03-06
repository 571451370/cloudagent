/*
Package guest_ping - ping guest

Example:
        { "execute": "guest-ping", "arguments": {}}
*/
package guest_ping

import "github.com/vtolstov/cloudagent/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-ping",
		Func:    fnGuestPing,
		Enabled: true,
	})
}

func fnGuestPing(req *qga.Request) *qga.Response {
	return &qga.Response{}
}
