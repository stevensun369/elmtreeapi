package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"

	// env
	"backend-go/env"
)

// main client
var Client *mongo.Client

// collections
var AverageMarks *mongo.Collection
var Grades *mongo.Collection
var Marks *mongo.Collection
var Parents *mongo.Collection
var Periods *mongo.Collection
var Schools *mongo.Collection
var Students *mongo.Collection
var Subjects *mongo.Collection
var Teachers *mongo.Collection
var TermMarks *mongo.Collection
var Truancies *mongo.Collection

// sort types
var DateSort interface{} = bson.D{
  {Key: "dateMonth", Value: 1}, 
  {Key: "dateDay", Value: 1},
}
var TermSort interface{} = bson.D{}
var EmptySort interface{} = bson.D{
  {Key: "term", Value: 1},
}
var PeriodSort interface{} = bson.D{
  {Key: "day", Value: 1}, 
  {Key: "interval", Value: 1},
}
var GradeSort interface{} = bson.D{
  {Key: "grade.gradeNumber", Value: 1}, 
  {Key: "grade.gradeLetter", Value: 1},
}

// initializing database
func InitDatabase() {
  var err error

  // client
  Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(env.MongoURI))

  if err != nil {
    log.Fatal(err)
  }

  // setting up collections
  AverageMarks = getCollection("averagemarks")
  Grades = getCollection("grades")
  Marks = getCollection("marks")
  Parents = getCollection("parents")
  Periods = getCollection("periods")
  Schools = getCollection("schools")
  Students = getCollection("students")
  Subjects = getCollection("subjects")
  Teachers = getCollection("teachers")
  TermMarks = getCollection("termmarks")
  Truancies = getCollection("truancies")
}

func GetCollection(collectionName string) (*mongo.Collection, error) {
  collection := Client.Database("elmtree").Collection(collectionName)

  return collection, nil
}

func getCollection(collectionName string) (*mongo.Collection) {
  return Client.Database("elmtree").Collection(collectionName)
}
