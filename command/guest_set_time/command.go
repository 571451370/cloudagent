// +build linux freebsd openbsd netbsd

/*
Package guest_set_time - set guest time

Example:
        { "execute": "guest-set-time", "arguments": {
            "time": int // optional, time to set
          }
        }
*/
package guest_set_time

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/vtolstov/cloudagent/qga"

	"golang.org/x/sys/unix"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-set-time",
		Func:    fnGuestSetTime,
		Enabled: true,
	})
}

func fnGuestSetTime(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}

	reqData := reqDataSetTime{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	args := []string{}
	if reqData.Time != 0 {
		tv := newTimeval(reqData.Time)

		tv.Sec = tv.Sec / 1000000000
		tv.Usec = tv.Usec % 1000000000 / 1000

		if err = unix.Settimeofday(tv); err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		args = append(args, "-w")
	} else {
		args = append(args, "-s")
	}

	cmd := exec.Command("hwclock", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	err = cmd.Run()
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	return res
}
