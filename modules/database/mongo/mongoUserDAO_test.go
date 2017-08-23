package mongo_test

import (
	"testing"

	"github.com/normegil/zookeeper-rest/modules/database/mongo"
	"github.com/normegil/zookeeper-rest/modules/model"
	daoPkg "github.com/normegil/zookeeper-rest/modules/model/dao"
	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func TestMongoUserDAO(t *testing.T) {
	if testing.Short() {
		t.Skip("Test require docker (Which make it slow). Skipping.")
	}

	user := model.UserImpl{"test", "pass"}

	dbName := "zookeeper-rest"
	dbNameWithoutCollection := "zookeeper-rest-no-collection"
	session, closeFn := mongo.Test_NewMongoSession(t)
	defer closeFn()
	mongo.Test_Insert(t, session, []mongo.Test_MongoData{
		{
			Database:   dbName,
			Collection: "users",
			Datas:      []interface{}{user},
		},
		{
			Database: dbNameWithoutCollection,
		},
	})

	t.Run("Load", func(t *testing.T) {
		dao := mongo.MongoUserDAO{session, dbName}
		result, err := dao.Load(user.Name())
		if nil != err {
			t.Fatal(errors.Wrap(err, "Error while loading user '"+user.Name()+"'"))
		}
		if nil == result {
			t.Error("Loaded user is nil")
		} else if *result != user {
			t.Error(test.Format("Loaded user is not the expected one.", user.String(), result.String()))
		}
	})

	t.Run("Load inexisting database", func(t *testing.T) {
		dao := mongo.MongoUserDAO{session, uuid.NewV4().String()}
		result, err := dao.Load(uuid.NewV4().String())
		if nil != err {
			t.Fatal(errors.Wrap(err, "Error while loading user"))
		}

		if nil != result {
			t.Errorf("Loaded user should be nil. Found %+v", result)
		}
	})

	t.Run("Inexisting collection", func(t *testing.T) {
		dao := mongo.MongoUserDAO{session, dbNameWithoutCollection}
		result, err := dao.Load(uuid.NewV4().String())
		if nil != err {
			t.Fatal(errors.Wrap(err, "Error while loading user"))
		}
		if nil != result {
			t.Errorf("Loaded user should be nil. Found %+v", result)
		}
	})

	t.Run("Inexisting user", func(t *testing.T) {
		dao := mongo.MongoUserDAO{session, dbName}
		result, err := dao.Load(uuid.NewV4().String())
		if nil != err {
			t.Fatal(errors.Wrap(err, "Error while loading user"))
		}

		if nil != result {
			t.Errorf("Loaded user should be nil. Found %+v", result)
		}
	})

	t.Run("Conformity to interface", func(t *testing.T) {
		dao := &mongo.MongoUserDAO{session, dbName}
		daoPkg.Test_Load(t, dao)
	})

}
