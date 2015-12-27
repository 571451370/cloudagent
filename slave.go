package main

import (
	"fmt"

	"github.com/vtolstov/cloudagent/channel"
)

func slave() error {
	var ch channel.Channel
	var err error

	switch options.Method {
	/*
		case "fd-serial":
			if ch, err = channel.NewFdChannel(options.Path); err != nil {
				return err
			}
	*/
	case "virtio-serial":
		if ch, err = channel.NewVirtioChannel(options.Path); err != nil {
			return err
		}
		/*
			case "isa-serial":
				if ch, err = channel.NewIsaChannel(options.Path); err != nil {
					return err
				}
		*/
	default:
		return fmt.Errorf("unsupported method %s", options.Method)
	}
	if err = ch.Open(); err != nil {
		return err
	}
	defer ch.Close()

	return ch.Poll()
}
