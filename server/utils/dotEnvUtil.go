package utils

import (
	"go-simple-live-polling/logger"
	"os"
	"regexp"
	"sync"

	"github.com/joho/godotenv"
)

// set to "" when Docker build
const projectDirName = "server"

// const projectDirName = ""

var (
	once       sync.Once
	dotEnvUtil *dotEnvUtilStruct
)

type IDotEnvUtil interface {
	GetEnvVariable(key string) string
}

type dotEnvUtilStruct struct{}

func GetDotEnvUtil() *dotEnvUtilStruct {
	if dotEnvUtil == nil {
		once.Do(func() {
			dotEnvUtil = &dotEnvUtilStruct{}
		})
	}
	return dotEnvUtil
}

func (d *dotEnvUtilStruct) GetEnvVariable(key string) string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	// logger.Logger.Info("Current work directory :" + currentWorkDirectory)
	// logger.Logger.Info("Path: " + string(rootPath) + `/.env`)

	if err != nil {
		logger.Logger.Fatal("Error loading .env file in dotEnvHelper.go " + err.Error())
	}

	return os.Getenv(key)
}
