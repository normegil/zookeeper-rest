package mongo

type Session interface {
	Copy() Session
	DB(name string) Database
	DatabaseNames() ([]string, error)
	DefaultDB() Database
	Close()
}

type Database interface {
	C(name string) Collection
	DropDatabase() error
}

type Collection interface {
	Find(query interface{}) Query
	Insert(docs ...interface{}) error
}

type Query interface {
	One(result interface{}) error
}
