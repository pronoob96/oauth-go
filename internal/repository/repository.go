package repository

import (
	"context"
	"log"
	m "oauth/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	AddUser(context.Context, *m.User) (*m.User, error)
	GetUserByUsername(context.Context, string) (*m.User, error)
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(collection *mongo.Collection) UserRepository {
	return &userRepo{
		collection: collection,
	}
}

func (r *userRepo) AddUser(ctx context.Context, user *m.User) (*m.User, error) {
	insertUser, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Inserted a User document: ", insertUser.InsertedID)
	user.ID = insertUser.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*m.User, error) {

	filter := bson.D{
		bson.E{Key: "username", Value: username},
	}

	var user m.User

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Fetched User: ", user)
	return &user, nil
}
