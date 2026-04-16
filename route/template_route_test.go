package route

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"zenrows-ta/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSetupTemplateRoute_GetTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("get template", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "db.templates", mtest.FirstBatch,
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
		)

		templateService := service.NewTemplateService(mt.DB)
		r := gin.New()
		SetupTemplateRoute(r, templateService)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/template", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			mt.Fatalf("expected 200 status, got %d", w.Code)
		}

		var payload map[string][]map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
			mt.Fatalf("failed to unmarshal response: %v", err)
		}
		if _, ok := payload["templates"]; !ok {
			mt.Fatal("expected templates key in response")
		}
	})
}

func TestSetupTemplateRoute_GetTemplatesAlias(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("templates alias", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "db.templates", mtest.FirstBatch,
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

		templateService := service.NewTemplateService(mt.DB)
		r := gin.New()
		SetupTemplateRoute(r, templateService)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/templates", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			mt.Fatalf("expected 200 status, got %d", w.Code)
		}
	})
}

func TestSetupTemplateRoute_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("internal error", func(mt *mtest.T) {
		mt.AddMockResponses(
			bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "connection error"}},
		)

		templateService := service.NewTemplateService(mt.DB)
		r := gin.New()
		SetupTemplateRoute(r, templateService)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/template", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			mt.Fatalf("expected 500 status, got %d", w.Code)
		}
	})
}
