package models

import (
	"fmt"
	"winddies/manage-api/global"

	"gopkg.in/mgo.v2"
)

var Session *mgo.Session

func MongoInit() {
	fmt.Println(global.Conf.Mongo.Addr)
	var err error
	Session, err = mgo.Dial("mongodb://" + global.Conf.Mongo.Addr)
	if err != nil {
		panic(err)
	}
	Session.SetMode(mgo.Monotonic, true)
}

type mongo struct{}

var mgoQuery *mongo

//
func (_ *mongo) Query(collectionName string, operation func(*mgo.Collection)) {
	sessionCp := Session.Copy()
	defer sessionCp.Close()

	database := sessionCp.DB(global.Conf.Mongo.DB)
	c := database.C(collectionName)
	operation(c)

}
