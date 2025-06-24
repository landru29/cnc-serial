package serial

import (
	"strings"

	serialport "go.bug.st/serial"
)

type PortName string

func (p *PortName) Decode(value string) error {
	if strings.ToLower(value) == "default" {
		ports, err := serialport.GetPortsList()
		if err == nil && len(ports) > 0 {
			*p = PortName(ports[0])
		}

		return nil
	}

	*p = PortName(value)

	return nil
}

func (p PortName) String() string {
	return string(p)
}

func (p *PortName) Set(value string) error {
	return p.Decode(value)
}

func (p PortName) Type() string {
	return "portName"
}
