package node

import (
	"github.com/normegil/resterrors"
	"github.com/normegil/zookeeper-rest/modules/middleware"
	"github.com/normegil/zookeeper-rest/modules/model"
	"github.com/normegil/zookeeper-rest/router"
	"github.com/samuel/go-zookeeper/zk"
)

const ADDRESS string = "127.0.0.1"
const BASE_PATH = "/rest/node"

type Controller struct {
	model.Env
}

func (c Controller) Routes() []router.Route {
	errCodeAssigner := middleware.ErrorHandlerFactory{
		ErrorCodeAssignerFunc: c.assignCode,
		Logger:                c.Log(),
	}
	return []router.Route{
		router.NewRoute(GET_METHOD, GET_PATH, errCodeAssigner.New(c.load).Handle),
		router.NewRoute(GET_ALL_METHOD, GET_ALL_PATH, errCodeAssigner.New(c.load).Handle),
		router.NewRoute(CREATE_METHOD, CREATE_PATH, errCodeAssigner.New(c.create).Handle),
		router.NewRoute("POST", CREATE_PATH, errCodeAssigner.New(c.create).Handle),
		router.NewRoute(UPDATE_METHOD, UPDATE_PATH, errCodeAssigner.New(c.update).Handle),
		router.NewRoute("PUT", UPDATE_PATH, errCodeAssigner.New(c.update).Handle),
		router.NewRoute(DELETE_METHOD, DELETE_PATH, errCodeAssigner.New(c.remove).Handle),
	}
}

func (c Controller) assignCode(err error) error {
	switch err {
	default:
		return err
	case zk.ErrNoServer:
		return resterrors.NewErrWithCode(50301, err)
	}
}
