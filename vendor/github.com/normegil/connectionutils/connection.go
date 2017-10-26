// Library of utilities used to manage connections
package connectionutils

import (
	"net"

	"github.com/normegil/interval"
)

// TCPPortAvalaible check if the given port is available for listening on current machine.
func TCPPortAvalaible(addr *net.TCPAddr) bool {
	conn, err := net.ListenTCP("tcp", addr)
	if nil != err {
		return false
	}
	defer conn.Close()
	return true
}

// SelectPort will select an available port in the given interval
func SelectPort(addr net.IP, possibilities interval.IntervalInteger) net.TCPAddr {
	return SelectPortExcluding(addr, possibilities, make([]int, 0))
}

// SelectPortExcluding will select an available port in the given interval, exluding the specified ports
func SelectPortExcluding(addr net.IP, possibilities interval.IntervalInteger, excluding []int) net.TCPAddr {
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
