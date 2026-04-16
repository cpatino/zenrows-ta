package main

import (
	"zenrows-ta/repository"
	"zenrows-ta/route"
	"zenrows-ta/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	db := repository.InitConnection()

	r := gin.Default()
	route.SetupTemplateRoute(r, service.NewTemplateService(db))

	authGroup := r.Group("/")
	authGroup.Use(service.NewUserService(db).Authenticate)
	route.SetupDeviceProfileRoutes(authGroup, service.NewDeviceProfileService(db))

	logrus.Info("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("Server failed to start")
	}
}
