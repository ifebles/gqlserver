package controller

import (
	"fmt"
	"gqlserver/database"
	"gqlserver/graph/model"
	"gqlserver/util"
	"log"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ControllerBook interface {
	Save(book *model.Book)
	FindAll() []*model.Book
}

func Save(input *model.NewBook) (*model.Book, any) {
	document := &model.Book{
		ID:    fmt.Sprintf("T%d", rand.Int()),
		Title: input.Title,
		Author: &model.User{
			ID:   input.UserID,
			Name: input.Name,
		},
	}

	client, ctx, cancel, err := database.Connect(util.Config.Server.URI)

	if err != nil {
		panic(err)
	}

	defer database.Close(client, ctx, cancel)
	cursor, err := database.SaveOne(
		client,
		ctx,
		util.Config.Server.Database.Name,
		util.Config.Server.Database.Collection,
		document,
	)

	if err != nil {
		panic(err)
	}

	return document, cursor.InsertedID
}

func FindAll() []*model.Book {
	client, ctx, cancel, err := database.Connect(util.Config.Server.URI)

	if err != nil {
		panic(err)
	}

	defer database.Close(client, ctx, cancel)
	cursor, err := database.Query(
		client,
		ctx,
		util.Config.Server.Database.Name,
		util.Config.Server.Database.Collection,
		primitive.D{},
	)

	if err != nil {
		panic(err)
	}

	///

	results := []*model.Book{}

	for cursor.Next(ctx) {
		var v *model.Book
		if err := cursor.Decode(&v); err != nil {
			log.Fatal(err)
		}

		results = append(results, v)
	}

	return results
}
