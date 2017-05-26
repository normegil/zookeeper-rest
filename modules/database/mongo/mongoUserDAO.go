package mongo

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"github.com/normegil/zookeeper-rest/modules/security"
)

type MongoUserDAO struct {
	Connection *Mongo
}

func (m *MongoUserDAO) Load(username string) (*security.User, error) {
	s := m.Connection.Session()
	defer s.Close()

	usrCollection := s.DB(m.Connection.Database()).C("users")
	result := security.User{}
	err := usrCollection.Find(bson.M{"name": username}).One(&result)
	if nil != err {
		return nil, errors.Wrap(err, "Loading user ("+username+") from MongoDB")
	}
	return &result, nil
}
