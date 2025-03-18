package main

import serialport "go.bug.st/serial"

func defaultPort() string {
	ports, err := serialport.GetPortsList()
	if err == nil && len(ports) > 0 {
		return ports[0]
	}

	return ""
}
