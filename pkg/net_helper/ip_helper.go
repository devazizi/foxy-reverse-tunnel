package net_helper

import "net"

func GetIp(remoteAddr net.Addr) string {
	ip, _, _ := net.SplitHostPort(remoteAddr.String())

	return ip
}
