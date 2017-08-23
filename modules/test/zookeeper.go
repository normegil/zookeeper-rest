package test

import (
	"net"
	"testing"
)

func NewZookeeper(t testing.TB) (net.TCPAddr, func()) {
	mainBinding := PortBinding{"tcp", 2181, "[52180;52187]"}
	info, closeFn := NewDocker(t, DockerOptions{
		Name:  "Zookeeper",
		Image: "zookeeper:latest",
		Ports: []PortBinding{
			mainBinding,
			{"tcp", 2888, "[52188;52195]"},
			{"tcp", 3888, "[53888;53895]"},
		},
	})
	return net.TCPAddr{
		IP:   info.Address,
		Port: info.Ports[mainBinding],
	}, closeFn
}
