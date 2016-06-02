/*
Package guest_set_user_password - sync host<->guest communication

Example:
        { "execute": "guest-set-user-password", "arguments": {
            "username": string // required, username to change password
            "password": string // required, base64 encoded password
            "crypted": bool // optional, specify that password already encrypted
          }
        }
*/
package guest_set_user_password

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/vtolstov/cloudagent/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-set-user-password",
		Func:      fnGuestSetUserPassword,
		Enabled:   true,
		Arguments: true,
	})
}

func setPwdChpasswd(user string, passwd []byte, crypted bool) error {
	args := []string{}

	if crypted {
		args = append(args, "-e")
	}
	cmd := exec.Command("chpasswd", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	arg := fmt.Sprintf("%s:%s", user, passwd)
	_, err = stdin.Write([]byte(arg))
	if err != nil {
		return err
	}
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func setPwdPasswd(user string, passwd []byte, crypted bool) error {
	args := []string{}

	args = append(args, "--stdin", user)
	cmd := exec.Command("passwd", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	arg := fmt.Sprintf("%s", passwd)
	_, err = stdin.Write([]byte(arg))
	if err != nil {
		return err
	}
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func fnGuestSetUserPassword(req *qga.Request) *qga.Response {
	res := &qga.Response{ID: req.ID}

	reqData := struct {
		User    string `json:"username"`
		Passwd  string `json:"password"`
		Crypted bool   `json:"crypted"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	passwd, err := base64.StdEncoding.DecodeString(reqData.Passwd)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	if reqData.Crypted {
		err = setPwdChpasswd(reqData.User, passwd, reqData.Crypted)
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		}
	} else {
		if err = setPwdChpasswd(reqData.User, passwd, reqData.Crypted); err != nil {
			err = setPwdPasswd(reqData.User, passwd, reqData.Crypted)
			if err != nil {
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			}
		}
	}
	return res
}
