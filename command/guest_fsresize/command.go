/*
Package guest_fsresize - run resize file system

Example:
        { "execute": "guest-fsresize", "arguments": {
            "path": string // optional, mounted filesystem path
          }
        }
*/
package guest_fsresize

import (
	"encoding/json"

	"github.com/vtolstov/cloudagent/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fsresize",
		Func:    fnGuestFsresize,
		Enabled: true,
	})
}

func fnGuestFsresize(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}

	reqData := struct {
		Path string `json:"path,omitempty"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	if reqData.Path == "" {
		reqData.Path = "/"
	}

	if err = resizefs(reqData.Path); err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
	}
	return res
}
