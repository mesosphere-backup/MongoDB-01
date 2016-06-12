package mesos

import (
	log "github.com/Sirupsen/logrus"
	"net"
)

func getDockerbrIP() string {
	return getInterfaceIP("dockerbr")
}

func getInterfaceIP(interfaceName string) string {
	iff, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Errorf("getInterfaceIP(%v) fail,%v", interfaceName, err)
		return ""
	}

	as, _ := iff.Addrs()
	for _, a := range as {
		ipnet, _ := a.(*net.IPNet)
		return ipnet.IP.String()
	}
	return ""
}
