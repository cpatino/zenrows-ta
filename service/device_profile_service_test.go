package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"zenrows-ta/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func newGinTestContextWithBody(method, url, body string) (*gin.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func TestDeviceProfileService_FindDeviceProfiles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.deviceProfiles", mtest.FirstBatch,
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

		service := NewDeviceProfileService(mt.DB)
		req := httptest.NewRequest(http.MethodGet, "/deviceProfiles", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		service.FindDeviceProfiles(c, userID)
		if w.Code != http.StatusOK {
			mt.Fatalf("expected 200 status, got %d", w.Code)
		}

		var profiles []model.DeviceProfile
		if err := json.Unmarshal(w.Body.Bytes(), &profiles); err != nil {
			mt.Fatalf("failed to unmarshal response: %v", err)
		}
		if len(profiles) != 1 {
			mt.Fatalf("expected 1 profile, got %d", len(profiles))
		}
	})
}

func TestDeviceProfileService_FindDeviceProfile_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("not found", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.deviceProfiles", mtest.FirstBatch),
		)

		service := NewDeviceProfileService(mt.DB)
		c, w := newGinTestContextWithBody(http.MethodGet, "/deviceProfiles/000000000000000000000000", "{}")
		c.Params = gin.Params{{Key: "id", Value: primitive.NewObjectID().Hex()}}

		service.FindDeviceProfile(c, primitive.NewObjectID())
		if w.Code != http.StatusNotFound {
			mt.Fatalf("expected 404 status, got %d", w.Code)
		}
	})
}

func TestDeviceProfileService_SaveDeviceProfile_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("bad request", func(mt *mtest.T) {
		service := NewDeviceProfileService(mt.DB)
		c, w := newGinTestContextWithBody(http.MethodPost, "/deviceProfiles", "")

		service.SaveDeviceProfile(c, primitive.NewObjectID())
		if w.Code != http.StatusBadRequest {
			mt.Fatalf("expected 400 status, got %d", w.Code)
		}
	})
}

func TestDeviceProfileService_SaveDeviceProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		profileID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(bson.E{Key: "insertedId", Value: profileID}),
		)

		service := NewDeviceProfileService(mt.DB)
		payload := `{"deviceType":"desktop","windowSize":{"width":1920,"height":1080},"userAgent":"Mozilla/5.0","countryCode":"US"}`
		c, w := newGinTestContextWithBody(http.MethodPost, "/deviceProfiles", payload)

		service.SaveDeviceProfile(c, userID)
		if w.Code != http.StatusCreated {
			mt.Fatalf("expected 201 status, got %d", w.Code)
		}

		var result map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			mt.Fatalf("failed to unmarshal response: %v", err)
		}
		if result["id"] == "" {
			mt.Fatal("expected created id in response")
		}
	})
}

func TestDeviceProfileService_UpdateDeviceProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		profileID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "db.deviceProfiles", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: profileID},
				{Key: "userId", Value: userID},
				{Key: "deviceType", Value: "desktop"},
				{Key: "windowSize", Value: bson.D{
					{Key: "width", Value: 1920},
					{Key: "height", Value: 1080},
				}},
				{Key: "userAgent", Value: "Mozilla/5.0"},
				{Key: "countryCode", Value: "US"},
			}),
			mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: 1}, bson.E{Key: "n", Value: 1}, bson.E{Key: "nModified", Value: 1}),
		)

		service := NewDeviceProfileService(mt.DB)
		payload := `{"deviceType":"mobile","windowSize":{"width":428,"height":926},"userAgent":"Mozilla/5.0 Mobile","countryCode":"ES"}`
		c, w := newGinTestContextWithBody(http.MethodPut, "/deviceProfiles/"+profileID.Hex(), payload)
		c.Params = gin.Params{{Key: "id", Value: profileID.Hex()}}

		service.UpdateDeviceProfile(c, userID)
		if w.Code != http.StatusNoContent {
			mt.Fatalf("expected 204 status, got %d", w.Code)
		}
	})
}

func TestDeviceProfileService_DeleteDeviceProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		profileID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "db.deviceProfiles", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: profileID},
				{Key: "userId", Value: userID},
				{Key: "deviceType", Value: "desktop"},
				{Key: "windowSize", Value: bson.D{
					{Key: "width", Value: 1920},
					{Key: "height", Value: 1080},
				}},
				{Key: "userAgent", Value: "Mozilla/5.0"},
				{Key: "countryCode", Value: "US"},
			}),
			mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: 1}, bson.E{Key: "n", Value: 1}),
		)

		service := NewDeviceProfileService(mt.DB)
		c, w := newGinTestContextWithBody(http.MethodDelete, "/deviceProfiles/"+profileID.Hex(), "")
		c.Params = gin.Params{{Key: "id", Value: profileID.Hex()}}

		service.DeleteDeviceProfile(c, userID)
		if w.Code != http.StatusNoContent {
			mt.Fatalf("expected 204 status, got %d", w.Code)
		}
	})
}
