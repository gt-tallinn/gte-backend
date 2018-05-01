package main

import (
	"net/http"
	"github.com/gt-tallinn/request-explorer/handlers/add"
	"github.com/sirupsen/logrus"
	"github.com/mongodb/mongo-go-driver/mongo"
	storage "github.com/gt-tallinn/request-explorer/storage/mongo"
	"os"
)

func main() {
	logger := logrus.WithField("app", "request-explorer")
	logrus.StandardLogger().SetLevel(5)

	mongoClient, err := mongo.NewClient(os.Getenv("EXPLORER_MONGODB"))
	if err != nil {
		logger.WithError(err).Fatal("Can't connect to mongodb")
	}
	mongoStorage := storage.New(mongoClient.Database("explorer").Collection("data"))
	add.New(logger, mongoStorage)
	http.Handle("/add", add.New(logger, mongoStorage))

	port := "9999"
	logger.Infof("start listen on :%s", port)
	err = http.ListenAndServe(":" + port, http.DefaultServeMux)
	if err != nil {
		logger.WithError(err).Fatal("Stop listen and serve")
	}
}
