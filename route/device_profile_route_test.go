package route

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"zenrows-ta/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetUserIDFromContext_ValidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", primitive.NewObjectID().Hex())

	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		t.Fatal("expected valid user ID from context")
	}
}

func TestGetUserIDFromContext_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "invalid-id")

	userID := getUserIDFromContext(c)
	if !userID.IsZero() {
		t.Fatal("expected nil object ID for invalid userID string")
	}
}

func TestSetupDeviceProfileRoutes_PostAndList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("post and list", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		profileID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(bson.E{Key: "insertedId", Value: profileID}),
			mtest.CreateCursorResponse(0, "db.deviceProfiles", mtest.FirstBatch,
				bson.D{
					{Key: "userId", Value: userID},
					{Key: "deviceType", Value: "desktop"},
					{Key: "windowSize", Value: bson.D{
						{Key: "width", Value: 1920},
						{Key: "height", Value: 1080},
					}},
					{Key: "userAgent", Value: "Mozilla/5.0"},
					{Key: "countryCode", Value: "US"},
				},
			),
			mtest.CreateCursorResponse(0, "db.deviceProfiles", mtest.FirstBatch,
				bson.D{
					{Key: "userId", Value: userID},
					{Key: "deviceType", Value: "desktop"},
					{Key: "windowSize", Value: bson.D{
						{Key: "width", Value: 1920},
						{Key: "height", Value: 1080},
					}},
					{Key: "userAgent", Value: "Mozilla/5.0"},
					{Key: "countryCode", Value: "US"},
				},
			),
		)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("userID", userID.Hex())
			c.Next()
		})

		deviceProfileService := service.NewDeviceProfileService(mt.DB)
		SetupDeviceProfileRoutes(r, deviceProfileService)

		payload := `{"deviceType":"desktop","windowSize":{"width":1920,"height":1080},"userAgent":"Mozilla/5.0","countryCode":"US"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/deviceProfiles", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			mt.Fatalf("expected 201 status, got %d", w.Code)
		}

		var created map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			mt.Fatalf("failed to parse created response: %v", err)
		}
		if created["id"] == "" {
			mt.Fatal("expected created profile id")
		}

		w = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/deviceProfiles", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			mt.Fatalf("expected 200 status for list, got %d", w.Code)
		}

		var profiles []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &profiles); err != nil {
			mt.Fatalf("failed to parse list response: %v", err)
		}
		if len(profiles) != 1 {
			mt.Fatalf("expected 1 profile in list, got %d", len(profiles))
		}
	})
}
