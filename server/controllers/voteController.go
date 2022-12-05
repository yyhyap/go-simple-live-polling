package controllers

import (
	"context"
	"fmt"
	"go-simple-live-polling/database"
	"go-simple-live-polling/enums"
	"go-simple-live-polling/logger"
	"go-simple-live-polling/models"
	"go-simple-live-polling/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	voteController *voteControllerStruct

	mongoDbClient     *mongo.Client               = database.MongoDbClient
	mongoDbConnection database.IMongoDbConnection = database.GetMongoDbConnection()
	voteCollection    *mongo.Collection           = mongoDbConnection.OpenCollection(mongoDbClient, "vote")

	dataValidationUtil utils.IDataValidationUtil = utils.GetDataValidationUtil()

	validate = validator.New()

	once sync.Once

	upgrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			// logger.Logger.Info(fmt.Sprintf("%s %s%s %v\n", r.Method, r.Host, r.URL, r.Proto))
			return r.Method == http.MethodGet
		},
	}

	clients = make(map[*websocket.Conn]struct{})

	mux = &sync.RWMutex{}
)

type IVoteController interface {
	HandleVoteWebSocket(c *gin.Context)
	CreateNewVote(c *gin.Context)
}

type voteControllerStruct struct{}

type NewVoteRequest struct {
	Voter_ic_no string `validate:"required"`
	Voter_name  string `validate:"required"`
	Party       string `validate:"required"`
}

func GetVoteController() *voteControllerStruct {
	if voteController == nil {
		once.Do(func() {
			voteController = &voteControllerStruct{}
		})
	}
	return voteController
}

func (v *voteControllerStruct) HandleVoteWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		logger.Logger.Error(fmt.Errorf("error on websocket connection: %v", err).Error())
		conn.Close()
		return
	}

	mux.Lock()
	clients[conn] = struct{}{}
	mux.Unlock()

	// get all votes and send to client
	groupStageByParty := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$party"},
			{Key: "totalCount", Value: bson.M{
				"$count": bson.M{},
			}},
		}},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "totalCount", Value: -1},
		}},
	}

	groupStage2 := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "data", Value: bson.M{
				"$push": "$$ROOT",
			}},
		}},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "data", Value: 1},
		}},
	}

	aggregateResult, err := voteCollection.Aggregate(context.Background(), mongo.Pipeline{groupStageByParty, sortStage, groupStage2, projectStage})

	if err != nil {
		logger.Logger.Error(err.Error())
		conn.Close()
		return
	}

	var data []bson.M

	err = aggregateResult.All(context.Background(), &data)

	if err != nil {
		logger.Logger.Error(err.Error())
		conn.Close()
		return
	}

	if len(data) > 0 {
		if len(data[0]) > 0 {
			dataBytes, err := json.Marshal(data[0])

			if err != nil {
				logger.Logger.Error(err.Error())
				conn.Close()
				return
			}

			conn.WriteMessage(websocket.TextMessage, dataBytes)
		}
	}

	logger.Logger.Info(fmt.Sprintf("Added a new connection. Total connection: %v", len(clients)))

	v.reader(conn)
}

func (v *voteControllerStruct) reader(conn *websocket.Conn) {
	for {
		// _, message, err := conn.ReadMessage()
		_, _, err := conn.ReadMessage()
		if err != nil {
			logger.Logger.Warn(err.Error())
			mux.Lock()
			delete(clients, conn)
			mux.Unlock()
			logger.Logger.Info(fmt.Sprintf("A client has been disconnected. Total connection: %v", len(clients)))
			return
		}
	}
}

func (v *voteControllerStruct) broadcast(msg []byte) {
	// use Lock instead of RLock, to prevent next broadcast request to interfere current map read operation (writing message to client)
	mux.Lock()
	wg := new(sync.WaitGroup)
	for client := range clients {
		wg.Add(1)
		go func(conn *websocket.Conn, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Logger.Error(err.Error())
			}
		}(client, wg)
	}
	wg.Wait()
	mux.Unlock()
}

func (v *voteControllerStruct) CreateNewVote(c *gin.Context) {
	var newVoteReq NewVoteRequest

	err := c.BindJSON(&newVoteReq)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = validate.Struct(newVoteReq)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match := dataValidationUtil.IsIcNoValid(newVoteReq.Voter_ic_no)

	if !match {
		logger.Logger.Warn("invalid ic no")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ic no"})
		return
	}

	partySelected, exist := enums.ParseString(newVoteReq.Party)

	if !exist {
		logger.Logger.Warn("invalid party")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// check ic no
	icNoCount, err := voteCollection.CountDocuments(ctx, bson.M{"voter_ic_no": newVoteReq.Voter_ic_no})

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// already voted
	if icNoCount > 0 {
		logger.Logger.Warn("voter with ic no already voted")
		c.JSON(http.StatusBadRequest, gin.H{"error": "voter with ic no already voted"})
		return
	}

	var vote models.Vote

	vote.ID = primitive.NewObjectID()
	vote.Vote_id = vote.ID.Hex()
	vote.Voter_ic_no = newVoteReq.Voter_ic_no
	vote.Voter_name = newVoteReq.Voter_name
	vote.Party = partySelected.String()

	// get local time
	created_at, err := time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	vote.Created_at = created_at

	insertResult, err := voteCollection.InsertOne(ctx, vote)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groupStageByParty := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$party"},
			{Key: "totalCount", Value: bson.M{
				"$count": bson.M{},
			}},
		}},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "totalCount", Value: -1},
		}},
	}

	groupStage2 := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "data", Value: bson.M{
				"$push": "$$ROOT",
			}},
		}},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "data", Value: 1},
		}},
	}

	aggregateResult, err := voteCollection.Aggregate(ctx, mongo.Pipeline{groupStageByParty, sortStage, groupStage2, projectStage})

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      fmt.Errorf("vote already submitted, another error: %v", err).Error(),
			"InsertedID": insertResult.InsertedID,
		})
		return
	}

	var data []bson.M

	err = aggregateResult.All(ctx, &data)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      fmt.Errorf("vote already submitted, another error: %v", err).Error(),
			"InsertedID": insertResult.InsertedID,
		})
		return
	}

	msg, err := json.Marshal(data[0])

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      fmt.Errorf("vote already submitted, another error: %v", err).Error(),
			"InsertedID": insertResult.InsertedID,
		})
		return
	}

	go v.broadcast(msg)

	c.JSON(http.StatusOK, insertResult)
}
