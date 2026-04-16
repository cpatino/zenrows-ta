package repository

import (
	"context"

	"zenrows-ta/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const templateCollectionName = "templates"

type TemplateRepository struct {
	collection *mongo.Collection
}

type TemplateDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeviceType  string             `bson:"deviceType" json:"deviceType"`
	WindowSize  WindowSizeDocument `bson:"windowSize" json:"windowSize"`
	UserAgent   string             `bson:"userAgent" json:"userAgent"`
	CountryCode string             `bson:"countryCode" json:"countryCode"`
}

func NewTemplateRepository(db *mongo.Database) *TemplateRepository {
	return &TemplateRepository{collection: db.Collection(templateCollectionName)}
}

func (repository *TemplateRepository) FindAllTemplates(ctx context.Context) ([]model.DeviceProfile, error) {
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templateDocuments []TemplateDocument
	if err := cursor.All(ctx, &templateDocuments); err != nil {
		return nil, err
	}

	templates := make([]model.DeviceProfile, len(templateDocuments))
	for i, templateDocument := range templateDocuments {
		templates[i] = model.DeviceProfile{
			ID:         templateDocument.ID.Hex(),
			DeviceType: templateDocument.DeviceType,
			WindowSize: model.WindowSize{
				Width:  templateDocument.WindowSize.Width,
				Height: templateDocument.WindowSize.Height,
			},
			UserAgent:   templateDocument.UserAgent,
			CountryCode: templateDocument.CountryCode,
		}
	}

	return templates, nil
}
