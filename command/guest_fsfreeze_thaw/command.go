/*
guest-fsfreeze-thaw - run unfreeze on all mounted file systems

Example:
        { "execute": "guest-fsfreeze-thaw", "arguments": {} }
*/
package guest_fsfreeze_freeze

import (
	"os"
	"unsafe"

	"github.com/vtolstov/cloudagent/qga"
	"github.com/vtolstov/go-ioctl"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-fsfreeze-thaw",
		Func:    fnGuestFsfreezeThaw,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestFsfreezeThaw(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}

	resData := 0

	fslist, err := qga.ListMountedFileSystems()
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	r := 0

	for _, fs := range fslist {
		f, err := os.Open(fs.Path)
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		err = ioctl.Fithaw(uintptr(f.Fd()), uintptr(unsafe.Pointer(&r)))
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		resData += 1
	}

	res.Return = resData
	return res
}
