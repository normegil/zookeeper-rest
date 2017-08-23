package dao

import "github.com/normegil/zookeeper-rest/modules/model"

type UserDAO interface {
	Load(username string) (*model.UserImpl, error)
}
