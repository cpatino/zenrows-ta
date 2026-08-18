package repository_test

import (
	"context"
	"testing"
	"zenrows-ta/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestTemplateRepository_FindAllTemplates_Success(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(2, "db.templates", mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: "deviceType", Value: "desktop"},
					{Key: "windowSize", Value: bson.D{
						{Key: "width", Value: 1920},
						{Key: "height", Value: 1080},
					}},
					{Key: "userAgent", Value: "Mozilla/5.0"},
					{Key: "countryCode", Value: "US"},
				},
			),
			mtest.CreateCursorResponse(0, "db.templates", mtest.NextBatch,
				bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: "deviceType", Value: "mobile"},
					{Key: "windowSize", Value: bson.D{
						{Key: "width", Value: 375},
						{Key: "height", Value: 812},
					}},
					{Key: "userAgent", Value: "Mozilla/5.0 Mobile"},
					{Key: "countryCode", Value: "GB"},
				},
			),
		)

		repo := repository.NewTemplateRepository(mt.DB)
		results, err := repo.FindAllTemplates(context.Background())
		if err != nil {
			mt.Fatalf("FindAllTemplates returned error: %v", err)
		}
		if len(results) != 2 {
			mt.Fatalf("expected 2 templates, got %d", len(results))
		}
	})
}

func TestTemplateRepository_FindAllTemplates_Error(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(
			bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "connection error"}},
		)

		repo := repository.NewTemplateRepository(mt.DB)
		_, err := repo.FindAllTemplates(context.Background())
		if err == nil {
			mt.Fatal("expected error in response")
		}
	})
}
