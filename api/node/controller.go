package node

import (
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/middleware"
	"github.com/normegil/zookeeper-rest/router"
	"github.com/samuel/go-zookeeper/zk"
)

const ADDRESS string = "127.0.0.1"

type Controller struct {
	environment.Env
}

func (c Controller) Routes() []router.Route {
	errCodeAssigner := middleware.ErrorHandlerFactory{
		ErrorCodeAssignerFunc: c.assignCode,
		Logger:                c.Log(),
	}
	return []router.Route{
		router.NewRoute("GET", "/rest/node", errCodeAssigner.New(c.load).Handle),
		router.NewRoute("GET", "/rest/node/:nodeID", errCodeAssigner.New(c.load).Handle),
		router.NewRoute("POST", "/rest/node", errCodeAssigner.New(c.create).Handle),
		router.NewRoute("PUT", "/rest/node", errCodeAssigner.New(c.create).Handle),
		router.NewRoute("POST", "/rest/node/:nodeID", errCodeAssigner.New(c.update).Handle),
		router.NewRoute("PUT", "/rest/node/:nodeID", errCodeAssigner.New(c.update).Handle),
		router.NewRoute("DELETE", "/rest/node/:nodeID", errCodeAssigner.New(c.remove).Handle),
	}
}

func (c Controller) assignCode(err error) error {
	switch err {
	default:
		return err
	case zk.ErrNoServer:
		return errors.NewErrWithCode(50301, err)
	}
}
