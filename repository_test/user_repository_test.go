package repository_test

import (
	"context"
	"testing"
	"time"
	"zenrows-ta/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUserRepository_FindUser_Success(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "db.users", mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: userID},
					{Key: "name", Value: "Test User"},
					{Key: "username", Value: "alice"},
					{Key: "password", Value: "secret"},
					{Key: "createdAt", Value: primitive.NewDateTimeFromTime(time.Now())},
				},
			),
		)

		repo := repository.NewUserRepository(mt.DB)
		got, err := repo.FindUser(context.Background(), "alice", "secret")
		if err != nil {
			mt.Fatalf("FindUser returned error: %v", err)
		}
		if got.Username != "alice" {
			mt.Fatalf("expected username alice, got %q", got.Username)
		}
	})
}

func TestUserRepository_FindUser_NotFound(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("not found", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "db.users", mtest.FirstBatch),
		)

		repo := repository.NewUserRepository(mt.DB)
		_, err := repo.FindUser(context.Background(), "alice", "wrong")
		if err == nil {
			mt.Fatal("expected not found error")
		}
		if err != mongo.ErrNoDocuments {
			mt.Fatalf("expected ErrNoDocuments, got %v", err)
		}
	})
}
