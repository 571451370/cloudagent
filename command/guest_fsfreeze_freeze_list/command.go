/*
guest-fsfreeze-freeze - run fsfreeze on all mounted file systems

Example:
        { "execute": "guest-fsfreeze-freeze", "arguments": {} }
*/
package guest_fsfreeze_freeze

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/vtolstov/cloudagent/qga"
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
	action := F_FREEZE_FS
	r := 0
	for _, fs := range fslist {
		f, err := os.Open(fs.Path)
		if err != nil {
			unfreeze()
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		_, _, err = syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(f.Fd()), uintptr(action), uintptr(unsafe.Pointer(&r)))
		if err != nil {
			unfreeze()
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		freezed = append(freezed, fs.Path)
		qga.StoreSet("guest-fsfreeze", "paths", freezed)
		resData += 1
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
				syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(f.Fd()), uintptr(F_THAW_FS), uintptr(unsafe.Pointer(&r)))
			}
		}
	}
}
