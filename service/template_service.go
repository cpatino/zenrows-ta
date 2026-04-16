package service

import (
	"net/http"

	"zenrows-ta/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateService struct {
	repository *repository.TemplateRepository
}

func NewTemplateService(db *mongo.Database) *TemplateService {
	repo := repository.NewTemplateRepository(db)
	return &TemplateService{repository: repo}
}

func (service *TemplateService) FindAllTemplates(ginCtx *gin.Context) {
	templates, err := service.repository.FindAllTemplates(ginCtx.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("Failed to find templates")
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load templates"})
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{"templates": templates})
}
