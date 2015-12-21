package main

import (
	"fmt"

	"github.com/vtolstov/cloudagent/channel"
)

func slave() error {
	var ch channel.Channel
	var err error

	switch options.Method {
	case "virtio-serial":
		if ch, err = channel.NewVirtioChannel(options.Path); err != nil {
			return err
		}
		err = ch.Open()
		/*
			case "isa-serial":
				if ch, err = qga.NewIsaChannel(options.Path); err != nil {
					return err
				}
				err = ch.Open()
		*/
	default:
		return fmt.Errorf("unsupported method %s", options.Method)
	}
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Poll()
}
