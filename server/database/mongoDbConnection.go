package database

import (
	"context"
	"go-simple-live-polling/logger"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once sync.Once

	MongoDbClient     *mongo.Client = DBClient()
	mongoDbConnection *mongoDbConnectionStruct
)

// set to "" when Docker build
const projectDirName = "server"

// const projectDirName = ""

type IMongoDbConnection interface {
	OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection
}

type mongoDbConnectionStruct struct{}

func GetMongoDbConnection() *mongoDbConnectionStruct {
	if mongoDbConnection == nil {
		once.Do(func() {
			mongoDbConnection = &mongoDbConnectionStruct{}
		})
	}
	return mongoDbConnection
}

func DBClient() *mongo.Client {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	// logger.Logger.Info("Current work directory :" + currentWorkDirectory)
	// logger.Logger.Info("Path: " + string(rootPath) + `/.env`)

	if err != nil {
		logger.Logger.Fatal("Error loading .env file in databaseConnection.go " + err.Error())
	}

	MongoDb := os.Getenv("MONGODB_URL")

	// MongoDB connection
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		logger.Logger.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		logger.Logger.Fatal(err.Error())
	}

	return client
}

func (m *mongoDbConnectionStruct) OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("go_live_polling_db").Collection(collectionName)
	return collection
}
