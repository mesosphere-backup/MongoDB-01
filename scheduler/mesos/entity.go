package mesos

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
)

type Summary struct {
	Cpu        float64
	Mem        float64
	PortRanges []*mesos.Value_Range
}

type Used struct {
	Cpu   float64
	Mem   float64
	Ports []uint64
}

func (u *Used) addPort(p uint64) {
	if u.Ports == nil {
		u.Ports = make([]uint64, 0)
	}
	u.Ports = append(u.Ports, p)
}

func (u *Used) isPortUsed(p uint64) bool {
	if u.Ports == nil {
		u.Ports = make([]uint64, 0)
	}
	for _, port := range u.Ports {
		if port == p {
			return true
		}
	}

	return false
}
