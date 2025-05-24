package disease_monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implDiseaseTypesAPI struct {
}

func NewDiseaseTypesApi() DiseaseTypesAPI {
	return &implDiseaseTypesAPI{}
}

func (o implDiseaseTypesAPI) GetDiseases(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
