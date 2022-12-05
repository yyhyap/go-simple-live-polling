package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vote struct {
	ID          primitive.ObjectID `bson:"_id"`
	Vote_id     string             `json:"vote_id" bson:"vote_id"`
	Voter_ic_no string             `json:"voter_ic_no" bson:"voter_ic_no"`
	Voter_name  string             `json:"voter_name" bson:"voter_name"`
	Party       string             `json:"party" bson:"party"`
	Created_at  time.Time          `json:"created_at" bson:"created_at"`
}
