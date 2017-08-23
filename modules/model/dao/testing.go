package dao

import (
	"testing"

	"github.com/normegil/zookeeper-rest/modules/model"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_Load(t *testing.T, dao UserDAO) {
	t.Run("Not Found", func(t *testing.T) {
		username := uuid.NewV4().String()
		user, err := dao.Load(username)
		if nil != err {
			t.Fatal(errors.Wrap(err, "Error while loading inexisting data (NotFoud should've return 'nil')"))
		}
		if nil != user {
			t.Fatal("User was found but shouldn't {UserID:" + username + "}")
		}
	})
}

type Test_UserDAO struct {
	Users []model.User
}

func (d *Test_UserDAO) Load(username string) (*model.UserImpl, error) {
	for _, user := range d.Users {
		if user.Name() == username {
			return &model.UserImpl{
				Username: user.Name(),
				Pass:     user.Password(),
			}, nil
		}
	}
	return nil, nil
}
