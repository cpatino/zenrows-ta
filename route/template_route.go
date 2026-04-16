package route

import (
	"zenrows-ta/service"

	"github.com/gin-gonic/gin"
)

func SetupTemplateRoute(r *gin.Engine, templateService *service.TemplateService) {
	r.GET("/template", templateService.FindAllTemplates)
	r.GET("/templates", templateService.FindAllTemplates)
}
