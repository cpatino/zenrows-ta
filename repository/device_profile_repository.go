package repository

import (
	"context"
	"time"

	"zenrows-ta/model"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const deviceProfileCollectionName = "deviceProfiles"
const userIdAttribute = "userId"
const idAttributeName = "_id"

type DeviceProfileRepository struct {
	collection *mongo.Collection
}

type DeviceProfileDocument struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	DeviceType    string             `bson:"deviceType" json:"deviceType"`
	WindowSize    WindowSizeDocument `bson:"windowSize" json:"windowSize"`
	UserAgent     string             `bson:"userAgent" json:"userAgent"`
	CountryCode   string             `bson:"countryCode" json:"countryCode"`
	CustomHeaders map[string]string  `bson:"customHeaders,omitempty" json:"customHeaders,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	ModifiedAt    *time.Time         `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
}

type WindowSizeDocument struct {
	Width  int `bson:"width" json:"width"`
	Height int `bson:"height" json:"height"`
}

func NewDeviceProfileRepository(db *mongo.Database) *DeviceProfileRepository {
	return &DeviceProfileRepository{collection: db.Collection(deviceProfileCollectionName)}
}

func (repository *DeviceProfileRepository) buildDeviceProfile(deviceProfileDocument *DeviceProfileDocument) model.DeviceProfile {
	return model.DeviceProfile{
		ID:         deviceProfileDocument.ID.Hex(),
		DeviceType: deviceProfileDocument.DeviceType,
		WindowSize: model.WindowSize{
			Width:  deviceProfileDocument.WindowSize.Width,
			Height: deviceProfileDocument.WindowSize.Height,
		},
		UserAgent:     deviceProfileDocument.UserAgent,
		CountryCode:   deviceProfileDocument.CountryCode,
		CustomHeaders: deviceProfileDocument.CustomHeaders,
		UserID:        deviceProfileDocument.UserID.Hex(),
	}
}

func (repository *DeviceProfileRepository) FindDeviceProfiles(ctx context.Context, userID primitive.ObjectID) ([]model.DeviceProfile, error) {
	cursor, err := repository.collection.Find(ctx, bson.M{userIdAttribute: userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deviceProfileDocuments []DeviceProfileDocument
	if err := cursor.All(ctx, &deviceProfileDocuments); err != nil {
		return nil, err
	}

	deviceProfiles := make([]model.DeviceProfile, len(deviceProfileDocuments))
	for i, deviceProfileDocument := range deviceProfileDocuments {
		deviceProfiles[i] = repository.buildDeviceProfile(&deviceProfileDocument)
	}

	return deviceProfiles, nil
}

func (repository *DeviceProfileRepository) FindDeviceProfile(ctx context.Context, id primitive.ObjectID) (model.DeviceProfile, error) {
	var deviceProfileDocument DeviceProfileDocument
	err := repository.collection.FindOne(ctx, bson.M{idAttributeName: id}).Decode(&deviceProfileDocument)
	if err != nil {
		logrus.WithError(err).Error("Failed to find device profile by ID")
		return model.DeviceProfile{}, err
	}
	return repository.buildDeviceProfile(&deviceProfileDocument), nil
}

func (repository *DeviceProfileRepository) InsertDeviceProfile(ctx context.Context, deviceProfile model.DeviceProfile) (*primitive.ObjectID, error) {
	userID, err := primitive.ObjectIDFromHex(deviceProfile.UserID)
	if err != nil {
		return nil, err
	}

	profile := &DeviceProfileDocument{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		DeviceType: deviceProfile.DeviceType,
		WindowSize: WindowSizeDocument{
			Width:  deviceProfile.WindowSize.Width,
			Height: deviceProfile.WindowSize.Height,
		},
		UserAgent:     deviceProfile.UserAgent,
		CountryCode:   deviceProfile.CountryCode,
		CustomHeaders: deviceProfile.CustomHeaders,
		CreatedAt:     time.Now(),
	}

	result, err := repository.collection.InsertOne(ctx, profile)
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	return &insertedID, nil
}

func (repository *DeviceProfileRepository) UpdateDeviceProfile(ctx context.Context, deviceProfile *model.DeviceProfile) error {
	now := time.Now()

	id, err := primitive.ObjectIDFromHex(deviceProfile.ID)
	if err != nil {
		return err
	}

	_, err = repository.collection.UpdateByID(
		ctx,
		id,
		bson.M{"$set": bson.M{
			"deviceType": deviceProfile.DeviceType,
			"windowSize": WindowSizeDocument{
				Width:  deviceProfile.WindowSize.Width,
				Height: deviceProfile.WindowSize.Height,
			},
			"userAgent":     deviceProfile.UserAgent,
			"countryCode":   deviceProfile.CountryCode,
			"customHeaders": deviceProfile.CustomHeaders,
			"modifiedAt":    &now,
		}},
	)
	return err
}

// DeleteProfile removes a device profile by ID
func (repository *DeviceProfileRepository) DeleteProfile(ctx context.Context, id primitive.ObjectID) error {
	_, err := repository.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
