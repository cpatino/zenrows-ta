package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"zenrows-ta/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestTemplateService_FindAllTemplates_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
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

		service := NewTemplateService(mt.DB)
		req := httptest.NewRequest(http.MethodGet, "/template", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		service.FindAllTemplates(c)
		if w.Code != http.StatusOK {
			mt.Fatalf("expected 200 status, got %d", w.Code)
		}

		var payload map[string][]model.DeviceProfile
		if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
			mt.Fatalf("failed to unmarshal response: %v", err)
		}

		templates, ok := payload["templates"]
		if !ok {
			mt.Fatal("expected templates key in response body")
		}
		if len(templates) != 2 {
			mt.Fatalf("expected 2 templates, got %d", len(templates))
		}
	})
}

func TestTemplateService_FindAllTemplates_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(
			bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "connection error"}},
		)

		service := NewTemplateService(mt.DB)
		req := httptest.NewRequest(http.MethodGet, "/templates", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		service.FindAllTemplates(c)
		if w.Code != http.StatusInternalServerError {
			mt.Fatalf("expected 500 status, got %d", w.Code)
		}
	})
}
