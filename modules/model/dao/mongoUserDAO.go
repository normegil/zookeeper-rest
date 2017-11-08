package dao

import (
	"github.com/normegil/mongo"
	"github.com/normegil/zookeeper-rest/modules/model"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoUserDAO struct {
	Connection mongo.Session
	Database   string
}

func (m *MongoUserDAO) Load(username string) (*model.UserImpl, error) {
	s := m.Connection.Copy()
	defer s.Close()

	db := s.DefaultDB()
	if "" != m.Database {
		db = s.DB(m.Database)
	}
	usrCollection := db.C("users")
	result := model.UserImpl{}
	err := usrCollection.Find(bson.M{"name": username}).One(&result)
	if nil != err {
		if mgo.ErrNotFound == err {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Loading user ("+username+") from MongoDB")
	}
	return &result, nil
}
