package main

import (
	"fmt"

	"github.com/vtolstov/cloudagent/qga"
)

func slave() error {
	var ch qga.Channel
	var err error

	switch options.Method {
	case "virtio-serial":
		if ch, err = qga.NewVirtioChannel(options.Path); err != nil {
			return err
		}
		err = ch.Open()
	case "isa-serial":
		if ch, err = qga.NewIsaChannel(options.Path); err != nil {
			return err
		}
		err = ch.Open()
	default:
		return fmt.Errorf("unsupported method %s", options.Method)
	}
	if err != nil {
		return err
	}
	defer ch.Close()

	return qga.WorkerIO(ch)
}
