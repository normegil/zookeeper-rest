package test

import (
	"net"
	"testing"

	"github.com/normegil/docker"
)

func NewZookeeper(t testing.TB) (net.TCPAddr, func()) {
	mainBinding := docker.PortBinding{"tcp", 2181, "[52180;52187]"}
	info, closeFn := NewDocker(t, docker.Options{
		Name:  "Zookeeper",
		Image: "zookeeper:latest",
		Ports: []docker.PortBinding{
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
