package mongo

import (
	"strconv"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

func NewMongo(url string, port int, database, user, pass string) (*Mongo, error) {
	mongo := &Mongo{
		url: url,
		port: port,
		database: database,
		user: user,
		pass: pass,
	}

	s, err := mgo.Dial(mongo.Address())
	if nil != err {
		return nil, errors.Wrap(err, "Cannot dial MongoDB at address " + mongo.Address())
	}
	mongo.session = s
	return mongo, nil
}

type Mongo struct {
	session *mgo.Session
	url string
	port int
	database string
	user string
	pass string
}

func (m *Mongo) Session() *mgo.Session {
	return m.session.Copy()
}

func (m Mongo) URL() string {
	return m.url
}

func (m Mongo) Port() int {
	return m.port
}

func (m Mongo) Address() string {
 return m.URL() + ":" + strconv.Itoa(m.Port())
}

func (m Mongo) Database() string {
	return m.database
}

func (m Mongo) User() string {
	return m.user
}

func (m Mongo) Password() string {
	return m.pass
}

func (m *Mongo) Close() {
	m.session.Close()
}
