package mgo

import (
	"context"
	"log"
	"time"

	"github.com/ggrrrr/bui_lib/config"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_URI      = "mongo.uri"
	MONGO_DB       = "mongo.db"
	MONGO_USERNAME = "mongo.username"
	MONGO_PASSWORD = "mongo.password"
)

var (
	Client   *mongo.Client
	Database *mongo.Database

	envParamsDefaults = []config.ParamValue{
		{
			Name:     MONGO_URI,
			Info:     "Mongo db uri mongodb://10.1.1.171:27017/",
			DefValue: "mongodb://10.1.1.171:27017/",
		},
		{
			Name:     MONGO_USERNAME,
			Info:     "mongo user",
			DefValue: "root",
		},
		{
			Name:     MONGO_PASSWORD,
			Info:     "pass",
			DefValue: "example",
		},
		{
			Name:     MONGO_DB,
			Info:     "database",
			DefValue: "test",
		},
	}
)

func init() {
	config.Configure(envParamsDefaults)

}

func Configure(rootCtx context.Context) error {
	user := viper.GetString(MONGO_USERNAME)
	pass := viper.GetString(MONGO_PASSWORD)
	uri := viper.GetString(MONGO_URI)
	db := viper.GetString(MONGO_DB)
	log.Printf("mongo: uri:%v, user: %v", uri, user)

	ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: user,
		Password: pass}
	var err error
	asd := options.Client()
	asd.ApplyURI(uri)
	asd.SetAuth(credential)
	Client, err := mongo.NewClient(asd)
	if err != nil {
		log.Printf("err: %v", err)
		return err
	}
	err = Client.Connect(ctx)
	if err != nil {
		return nil
	}
	d, err := Client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil
	}
	log.Println(d)

	Database = Client.Database(db)
	log.Printf("connected: db: %v", Database.Name())

	return nil
}
