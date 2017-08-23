package mongo

import (
	"net"
	"testing"

	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

type Test_MongoData struct {
	Database   string
	Collection string
	Datas      []interface{}
}

func Test_NewMongoSession(t testing.TB) (Session, func()) {
	mongoInfo, closeFn := test.NewMongo(t)
	host := net.TCPAddr{mongoInfo.Address, mongoInfo.Port, ""}
	session, err := mgo.Dial(host.String())
	if nil != err {
		defer closeFn()
		t.Fatal(errors.Wrap(err, "Could not create a new mgo session {Host:"+host.String()+"}"))
	}
	return NewSession(session, ""), func() {
		session.Close()
		closeFn()
	}
}

func Test_Insert(t testing.TB, mongo Session, datas []Test_MongoData) {
	for _, data := range datas {
		t.Log("Inserting data: " + data.Database + " " + data.Collection)
		collection := mongo.DB(data.Database).C(data.Collection)
		for _, toInsert := range data.Datas {
			t.Logf("Inserting data - %+v", toInsert)
			err := collection.Insert(toInsert)
			if nil != err {
				t.Fatal(errors.Wrapf(err, "Could not insert %+v"))
			}
		}
	}
}

func Test_Clean(t testing.TB, mongo Session) {
	names, err := mongo.DatabaseNames()
	if err != nil {
		t.Fatal(errors.Wrapf(err, "Loading databases names"))
	}
	for _, name := range names {
		if err = mongo.DB(name).DropDatabase(); nil != err {
			t.Fatal(errors.Wrapf(err, "Droping DB %s", name))
		}
	}
}
