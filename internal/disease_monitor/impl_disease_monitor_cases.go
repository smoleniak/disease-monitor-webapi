package disease_monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implDiseaseMonitorCasesAPI struct {
}

func NewDiseaseMonitorCasesApi() DiseaseMonitorCasesAPI {
	return &implDiseaseMonitorCasesAPI{}
}

func (o implDiseaseMonitorCasesAPI) CreateDiseaseCaseListEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDiseaseMonitorCasesAPI) DeleteDiseaseCaseEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDiseaseMonitorCasesAPI) GetDiseaseCaseEntries(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDiseaseMonitorCasesAPI) GetDiseaseCaseEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDiseaseMonitorCasesAPI) UpdateDiseaseCaseEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
