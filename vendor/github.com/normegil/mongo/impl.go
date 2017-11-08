package mongo

import (
	mgo "gopkg.in/mgo.v2"
)

// NewSession will create a new Session using an mgo.Session and a default Database
func NewSession(session *mgo.Session, databaseName string) Session {
	return &sessionImpl{
		Session:      session,
		DatabaseName: databaseName,
	}
}

type sessionImpl struct {
	Session      *mgo.Session
	DatabaseName string
}

func (s *sessionImpl) Copy() Session {
	return &sessionImpl{
		Session:      s.Session.Copy(),
		DatabaseName: s.DatabaseName,
	}
}

func (s *sessionImpl) DB(name string) Database {
	return &databaseImpl{
		Database: s.Session.DB(name),
	}
}

func (s *sessionImpl) DatabaseNames() ([]string, error) {
	return s.Session.DatabaseNames()
}

func (s *sessionImpl) DefaultDB() Database {
	return s.DB(s.DatabaseName)
}

func (s *sessionImpl) Close() {
	s.Session.Close()
}

type databaseImpl struct {
	*mgo.Database
}

func (d *databaseImpl) C(name string) Collection {
	return &collectionImpl{
		Collection: d.Database.C(name),
	}
}

func (d *databaseImpl) DropDatabase() error {
	return d.Database.DropDatabase()
}

type collectionImpl struct {
	*mgo.Collection
}

func (c *collectionImpl) Find(query interface{}) Query {
	return queryImpl{
		Query: c.Collection.Find(query),
	}
}

func (c *collectionImpl) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs...)
}

type queryImpl struct {
	*mgo.Query
}
