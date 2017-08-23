package tools

import "net"

func TCPPortAvalaible(addr *net.TCPAddr) bool {
	conn, err := net.ListenTCP("tcp", addr)
	if nil != err {
		return false
	}
	defer conn.Close()
	return true
}

func SelectPort(addr net.IP, possibilities IntervalInteger) net.TCPAddr {
	return SelectPortExcluding(addr, possibilities, make([]int, 0))
}

func SelectPortExcluding(addr net.IP, possibilities IntervalInteger, excluding []int) net.TCPAddr {
	for i := possibilities.LowestNumberIncluded(); i <= possibilities.HighestNumberIncluded(); i++ {
		if !contains(excluding, i) {
			toTest := net.TCPAddr{addr, i, ""}
			if TCPPortAvalaible(&toTest) {
				return toTest
			}
		}
	}
	return net.TCPAddr{IP: addr}
}

func contains(slice []int, value int) bool {
	for _, sliceValue := range slice {
		if value == sliceValue {
			return true
		}
	}
	return false
}
