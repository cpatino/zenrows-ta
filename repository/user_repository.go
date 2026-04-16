package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollectionName = "users"

type UserDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{collection: db.Collection(userCollectionName)}
}

func (repository *UserRepository) FindUser(ctx context.Context, username, password string) (*UserDocument, error) {
	var user UserDocument
	err := repository.collection.FindOne(ctx, bson.M{
		"username": username,
		"password": password,
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
