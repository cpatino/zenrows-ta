package service

import (
	"net/http"

	"zenrows-ta/model"
	"zenrows-ta/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeviceProfileService struct {
	repository *repository.DeviceProfileRepository
}

func NewDeviceProfileService(db *mongo.Database) *DeviceProfileService {
	repo := repository.NewDeviceProfileRepository(db)
	return &DeviceProfileService{repository: repo}
}

func getIdFromGinContext(ginCtx *gin.Context) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(ginCtx.Param("id"))
	if err != nil {
		logrus.WithError(err).Warn("Invalid ID format in request")
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return primitive.NilObjectID
	}
	return id
}

func isForbiddenDeviceProfile(ginCtx *gin.Context, deviceProfileDocument model.DeviceProfile, userID primitive.ObjectID) bool {
	if deviceProfileDocument.UserID != userID.Hex() {
		logrus.Warn("Access denied to device profile")
		ginCtx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return true
	}
	return false
}

func mapInputDeviceProfile(ginCtx *gin.Context) (*model.DeviceProfile, bool) {
	if ginCtx.Request.ContentLength == 0 {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Request body is required"})
		return nil, false
	}

	var deviceProfile model.DeviceProfile
	if err := ginCtx.ShouldBindJSON(&deviceProfile); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, false
	}

	if deviceProfile.WindowSize.Width <= 0 || deviceProfile.WindowSize.Height <= 0 {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Window size is required and dimensions must be positive integers"})
		return nil, false
	}

	return &deviceProfile, true
}

func (service *DeviceProfileService) FindDeviceProfiles(ginCtx *gin.Context, userID primitive.ObjectID) {
	deviceProfiles, err := service.repository.FindDeviceProfiles(ginCtx.Request.Context(), userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to list device profiles for user")
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list device profiles"})
		return
	}

	ginCtx.JSON(http.StatusOK, deviceProfiles)
}

func (service *DeviceProfileService) FindDeviceProfile(ginCtx *gin.Context, userID primitive.ObjectID) {
	existingDeviceProfile, ok := service.findOwnedDeviceProfile(ginCtx, getIdFromGinContext(ginCtx), userID)
	if !ok {
		return
	}

	ginCtx.JSON(http.StatusOK, existingDeviceProfile)
}

func (service *DeviceProfileService) SaveDeviceProfile(ginCtx *gin.Context, userID primitive.ObjectID) {
	inputDeviceProfile, ok := mapInputDeviceProfile(ginCtx)
	if !ok {
		return
	}

	inputDeviceProfile.UserID = userID.Hex()

	id, err := service.repository.InsertDeviceProfile(ginCtx.Request.Context(), *inputDeviceProfile)
	if err != nil {
		logrus.WithError(err).Error("Failed to create device profile")
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save device profile, contact support"})
		return
	}

	ginCtx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (service *DeviceProfileService) UpdateDeviceProfile(ginCtx *gin.Context, userID primitive.ObjectID) {
	inputDeviceProfile, ok := mapInputDeviceProfile(ginCtx)
	if !ok {
		return
	}

	existingDeviceProfile, ok := service.findOwnedDeviceProfile(ginCtx, getIdFromGinContext(ginCtx), userID)
	if !ok {
		return
	}

	existingDeviceProfile.DeviceType = inputDeviceProfile.DeviceType
	existingDeviceProfile.WindowSize = inputDeviceProfile.WindowSize
	existingDeviceProfile.UserAgent = inputDeviceProfile.UserAgent
	existingDeviceProfile.CountryCode = inputDeviceProfile.CountryCode
	existingDeviceProfile.CustomHeaders = inputDeviceProfile.CustomHeaders

	err := service.repository.UpdateDeviceProfile(ginCtx.Request.Context(), existingDeviceProfile)
	if err != nil {
		logrus.WithError(err).Error("Failed to update device profile")
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device profile, contact support"})
		return
	}

	ginCtx.JSON(http.StatusNoContent, nil)
}

func (service *DeviceProfileService) findOwnedDeviceProfile(ginCtx *gin.Context, id primitive.ObjectID, userID primitive.ObjectID) (*model.DeviceProfile, bool) {
	existingDeviceProfile, err := service.repository.FindDeviceProfile(ginCtx.Request.Context(), id)
	if err != nil {
		ginCtx.JSON(http.StatusNotFound, gin.H{"error": "Device Profile not found"})
		return nil, false
	}

	if isForbiddenDeviceProfile(ginCtx, existingDeviceProfile, userID) {
		return nil, false
	}
	return &existingDeviceProfile, true
}

func (service *DeviceProfileService) DeleteDeviceProfile(ginCtx *gin.Context, userID primitive.ObjectID) {
	id := getIdFromGinContext(ginCtx)

	_, ok := service.findOwnedDeviceProfile(ginCtx, id, userID)
	if !ok {
		return
	}

	err := service.repository.DeleteProfile(ginCtx.Request.Context(), id)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete device profile")
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device profile, contact support"})
		return
	}

	ginCtx.JSON(http.StatusNoContent, nil)
}
