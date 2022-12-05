package main

import (
	"go-simple-live-polling/controllers"
	"go-simple-live-polling/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	dotEnvUtil utils.IDotEnvUtil = utils.GetDotEnvUtil()

	voteController controllers.IVoteController = controllers.GetVoteController()
)

func main() {
	port := dotEnvUtil.GetEnvVariable("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	// default cors for localhost
	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ws/live-vote", voteController.HandleVoteWebSocket)
	router.POST("/api/create-vote", voteController.CreateNewVote)

	router.Run(":" + port)
}
