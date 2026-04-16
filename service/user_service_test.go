package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func newGinTestContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func TestUserService_Authenticate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userID := primitive.NewObjectID()
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "db.users", mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: userID},
					{Key: "name", Value: "Test User"},
					{Key: "username", Value: "alice"},
					{Key: "password", Value: "secret"},
					{Key: "createdAt", Value: primitive.NewDateTimeFromTime(time.Now())},
				},
			),
		)

		service := NewUserService(mt.DB)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth("alice", "secret")
		c, _ := newGinTestContext(req)

		service.Authenticate(c)
		if c.IsAborted() {
			mt.Fatal("expected authentication to succeed")
		}
	})
}

func TestUserService_Authenticate_NoCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("no credentials", func(mt *mtest.T) {
		service := NewUserService(mt.DB)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		c, w := newGinTestContext(req)

		service.Authenticate(c)
		if !c.IsAborted() {
			mt.Fatal("expected authentication to be aborted for missing credentials")
		}
		if w.Code != http.StatusForbidden {
			mt.Fatalf("expected 403 status, got %d", w.Code)
		}
	})
}

func TestUserService_Authenticate_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mt := newMongoTest(t)
	defer mt.Close()

	mt.Run("invalid credentials", func(mt *mtest.T) {
		// No mock response, so FindUser will fail

		service := NewUserService(mt.DB)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth("alice", "wrong")
		c, w := newGinTestContext(req)

		service.Authenticate(c)
		if !c.IsAborted() {
			mt.Fatal("expected authentication to be aborted for invalid credentials")
		}
		if w.Code != http.StatusUnauthorized {
			mt.Fatalf("expected 401 status, got %d", w.Code)
		}
	})
}
