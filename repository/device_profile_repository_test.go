package repository

import (
	"context"
	"testing"

	"zenrows-ta/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestDeviceProfileRepository_InsertAndFindDeviceProfile_Success(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		profileID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(bson.E{Key: "insertedId", Value: profileID}),
		)

		repo := NewDeviceProfileRepository(mt.DB)
		profile := model.DeviceProfile{
			DeviceType:    "desktop",
			WindowSize:    model.WindowSize{Width: 1920, Height: 1080},
			UserAgent:     "Mozilla/5.0",
			CountryCode:   "US",
			CustomHeaders: map[string]string{"x-test": "value"},
			UserID:        userID.Hex(),
		}
		insertedID, err := repo.InsertDeviceProfile(context.Background(), profile)
		if err != nil {
			mt.Fatalf("InsertDeviceProfile returned error: %v", err)
		}
		if insertedID == nil || insertedID.IsZero() {
			mt.Fatal("expected non-zero inserted ID")
		}
	})
}

func TestDeviceProfileRepository_FindDeviceProfiles_Success(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "db.deviceProfiles", mtest.FirstBatch,
				bson.D{
					{Key: "userId", Value: userID},
					{Key: "deviceType", Value: "desktop"},
					{Key: "windowSize", Value: bson.D{
						{Key: "width", Value: 1280},
						{Key: "height", Value: 720},
					}},
					{Key: "userAgent", Value: "Mozilla/5.0"},
					{Key: "countryCode", Value: "US"},
				},
			),
		)

		repo := NewDeviceProfileRepository(mt.DB)
		results, err := repo.FindDeviceProfiles(context.Background(), userID)
		if err != nil {
			mt.Fatalf("FindDeviceProfiles returned error: %v", err)
		}
		if len(results) != 1 {
			mt.Fatalf("expected 1 device profile, got %d", len(results))
		}
	})
}

func TestDeviceProfileRepository_UpdateDeviceProfile_Success(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		profileID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(bson.E{Key: "insertedId", Value: profileID}),
			mtest.CreateSuccessResponse(bson.E{Key: "modifiedCount", Value: 1}),
		)

		repo := NewDeviceProfileRepository(mt.DB)
		profile := model.DeviceProfile{
			DeviceType:    "desktop",
			WindowSize:    model.WindowSize{Width: 1024, Height: 768},
			UserAgent:     "Mozilla/5.0",
			CountryCode:   "US",
			CustomHeaders: map[string]string{"x-test": "value"},
			UserID:        userID.Hex(),
		}
		insertedID, err := repo.InsertDeviceProfile(context.Background(), profile)
		if err != nil {
			mt.Fatalf("InsertDeviceProfile returned error: %v", err)
		}

		updatedProfile := &model.DeviceProfile{
			ID:            insertedID.Hex(),
			DeviceType:    "mobile",
			WindowSize:    model.WindowSize{Width: 428, Height: 926},
			UserAgent:     "Mozilla/5.0 Mobile",
			CountryCode:   "CA",
			CustomHeaders: map[string]string{"x-updated": "true"},
			UserID:        userID.Hex(),
		}

		if err := repo.UpdateDeviceProfile(context.Background(), updatedProfile); err != nil {
			mt.Fatalf("UpdateDeviceProfile returned error: %v", err)
		}
	})
}

func TestDeviceProfileRepository_UpdateDeviceProfile_InvalidID(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("invalid id", func(mt *mtest.T) {
		repo := NewDeviceProfileRepository(mt.DB)
		err := repo.UpdateDeviceProfile(context.Background(), &model.DeviceProfile{ID: "invalid-id"})
		if err == nil {
			mt.Fatal("expected error for invalid object ID")
		}
	})
}

func TestDeviceProfileRepository_DeleteProfile_Success(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		profileID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(bson.E{Key: "insertedId", Value: profileID}),
			mtest.CreateSuccessResponse(bson.E{Key: "deletedCount", Value: 1}),
		)

		repo := NewDeviceProfileRepository(mt.DB)
		profile := model.DeviceProfile{
			DeviceType:    "desktop",
			WindowSize:    model.WindowSize{Width: 1280, Height: 720},
			UserAgent:     "Mozilla/5.0",
			CountryCode:   "US",
			CustomHeaders: map[string]string{"x-test": "delete"},
			UserID:        userID.Hex(),
		}
		insertedID, err := repo.InsertDeviceProfile(context.Background(), profile)
		if err != nil {
			mt.Fatalf("InsertDeviceProfile returned error: %v", err)
		}

		if err := repo.DeleteProfile(context.Background(), *insertedID); err != nil {
			mt.Fatalf("DeleteProfile returned error: %v", err)
		}
	})
}

func TestDeviceProfileRepository_InsertDeviceProfile_InvalidUserID(t *testing.T) {
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("invalid user id", func(mt *mtest.T) {
		repo := NewDeviceProfileRepository(mt.DB)
		_, err := repo.InsertDeviceProfile(context.Background(), model.DeviceProfile{UserID: "not-an-object-id"})
		if err == nil {
			mt.Fatal("expected error when saving profile with invalid user ID")
		}
	})
}
