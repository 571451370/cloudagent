/*
Package guest_fsfreeze_freeze - run fsfreeze on all mounted file systems

Example:
        { "execute": "guest-fsfreeze-freeze", "arguments": {} }
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
		Name:    "guest-fsfreeze-freeze",
		Func:    fnGuestFsfreezeFreeze,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestFsfreezeFreeze(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}

	resData := 0

	fslist, err := qga.ListMountedFileSystems()
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	freezed := []string{}
	r := 0

	for _, fs := range fslist {
		f, err := os.Open(fs.Path)
		if err != nil {
			unfreeze()
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		err = ioctl.Fifreeze(uintptr(f.Fd()), uintptr(unsafe.Pointer(&r)))
		f.Close()
		if err != nil {
			unfreeze()
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		freezed = append(freezed, fs.Path)
		qga.StoreSet("guest-fsfreeze", "paths", freezed)
		resData++
	}

	res.Return = resData
	return res
}

func unfreeze() {
	r := 0
	defer qga.StoreDel("guest-fsfreeze", "paths")
	if paths, ok := qga.StoreGet("guest-fsfreeze", "paths"); ok {
		for _, path := range paths.([]string) {
			if f, err := os.Open(path); err == nil {
				ioctl.Fithaw(uintptr(f.Fd()), uintptr(unsafe.Pointer(&r)))
				f.Close()
			}
		}
	}
}
