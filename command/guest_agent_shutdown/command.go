/*
Package guest_agent_shutdown - shutdown cloudagent inside vm

Example:
        { "execute": "guest-agent-shutdown", "arguments": {
            "timeout": int // optional, wait time for shutdown
          }
        }
*/
package guest_agent_shutdown

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vtolstov/cloudagent/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-agent-shutdown",
		Func:      fnGuestAgentShutdown,
		Enabled:   true,
		Arguments: true,
	})
}

func fnGuestAgentShutdown(req *qga.Request) *qga.Response {
	res := &qga.Response{}
	var dt time.Duration

	reqData := struct {
		Timeout string `json:"timeout,omitempty"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	if reqData.Timeout == "" {
		dt = 1
	}

	dt, err = time.ParseDuration(fmt.Sprintf("%ds", reqData.Timeout))
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	time.Sleep(dt)

	return res
}
