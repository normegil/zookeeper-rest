package mongo

// Session is a MongoSession, with an optional default database name. Methods map to a subset of mgo.v2 methods.
type Session interface {
	Copy() Session
	DB(name string) Database
	DatabaseNames() ([]string, error)
	DefaultDB() Database
	Close()
}

// Database is a mongo database. Methods map to a subset of mgo.v2 methods.
type Database interface {
	C(name string) Collection
	DropDatabase() error
}

// Collection is a mongo collection. Methods map to a subset of mgo.v2 methods.
type Collection interface {
	Find(query interface{}) Query
	Insert(docs ...interface{}) error
}

// Query is a mongo query. Methods map to a subset of mgo.v2 methods.
type Query interface {
	One(result interface{}) error
}
