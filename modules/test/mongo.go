package test

import (
	"net"
	"testing"

	"github.com/normegil/docker"
)

const mongoInternalPort = 27017
const MONGO_PORTS string = "[50017;50037]"

func NewMongo(t testing.TB) (MongoInfo, func()) {
	mainPortBinding := docker.PortBinding{"tcp", mongoInternalPort, MONGO_PORTS}
	info, close := NewDocker(t, docker.Options{
		Name:  "MongoDB",
		Image: "mongo:latest",
		Ports: []docker.PortBinding{mainPortBinding},
	})
	return MongoInfo{
		Address: info.Address,
		Port:    info.Ports[mainPortBinding],
	}, close
}

type MongoInfo struct {
	Address net.IP
	Port    int
}
