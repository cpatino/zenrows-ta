package service

import (
	"net/http"
	"zenrows-ta/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(db *mongo.Database) *UserService {
	repository := repository.NewUserRepository(db)
	return &UserService{repository: repository}
}

func (service *UserService) Authenticate(ginCtx *gin.Context) {
	username, password, ok := ginCtx.Request.BasicAuth()
	if !ok {
		logrus.Warn("No basic auth credentials provided")
		ginCtx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		ginCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
		return
	}

	userDocument, err := service.repository.FindUser(ginCtx.Request.Context(), username, password)
	if err != nil {
		logrus.WithField("username", username).Warn("Failed to authenticate user")
		ginCtx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ginCtx.Set("userID", userDocument.ID.Hex())
	ginCtx.Next()
}
