// +build !linux

/*
Package guest_fstrim - run fstrim on path

Example:
        { "execute": "guest-fstrim", "arguments": {
            "minimum": int // optional, minimum trimmed range
          }
        }
*/
package guest_fstrim

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/vtolstov/cloudagent/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fstrim",
		Func:    fnGuestFstrim,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestFstrim(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}

	reqData := struct {
		Minimum int `json:"minimum,omitempty"`
	}{}

	type resPath struct {
		Path    string `json:"path"`
		Trimmed *int   `json:"trimmed,omitempty"`
		Minimum *int   `json:"minimum,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	resData := struct {
		Paths []*resPath `json:"paths"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	fslist, err := qga.ListMountedFileSystems()
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	for _, fs := range fslist {
		switch fs.Type {
		case "ufs", "ffs":
			cmd : = exec.Command("fsck_"+fs.Type, "-B", "-E", fs.Path)
                        cmd.Env = append(cmd.Env, os.Environ()...)
                        err = cmd.Run()
		default:
			err = fmt.Errorf("unsupported")
		}
		rpath := &resPath{Path: fs.Path}
		if err != nil {
			rpath.Error = err.Error()
		}
		resData.Paths = append(resData.Paths, rpath)
	}

	res.Return = resData
	return res
}
