package route

import (
	"zenrows-ta/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const DeviceProfilesBasePath = "/deviceProfiles"

func SetupDeviceProfileRoutes(authRoutes gin.IRoutes, deviceProfileService *service.DeviceProfileService) {
	authRoutes.GET(DeviceProfilesBasePath, func(ginCtx *gin.Context) {
		deviceProfileService.FindDeviceProfiles(ginCtx, getUserIDFromContext(ginCtx))
	})

	authRoutes.GET(DeviceProfilesBasePath+"/:id", func(ginCtx *gin.Context) {
		deviceProfileService.FindDeviceProfile(ginCtx, getUserIDFromContext(ginCtx))
	})

	authRoutes.POST(DeviceProfilesBasePath, func(ginCtx *gin.Context) {
		deviceProfileService.SaveDeviceProfile(ginCtx, getUserIDFromContext(ginCtx))
	})

	authRoutes.PUT(DeviceProfilesBasePath+"/:id", func(ginCtx *gin.Context) {
		deviceProfileService.UpdateDeviceProfile(ginCtx, getUserIDFromContext(ginCtx))
	})

	authRoutes.DELETE(DeviceProfilesBasePath+"/:id", func(ginCtx *gin.Context) {
		deviceProfileService.DeleteDeviceProfile(ginCtx, getUserIDFromContext(ginCtx))
	})
}

func getUserIDFromContext(ginCtx *gin.Context) primitive.ObjectID {
	userIDStr := ginCtx.GetString("userID")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		logrus.WithError(err).WithField("userID", userIDStr).Error("Invalid user ID in context")
		return primitive.NilObjectID
	}
	return userID
}
