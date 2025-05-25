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
	updateRegionFunc(c, func(
		c *gin.Context,
		region *Region,
	) (updatedRegion *Region, responseContent interface{}, status int) {
		result := region.PredefinedDiseases
		if result == nil {
			result = []Disease{}
		}
		return nil, result, http.StatusOK
	})
}
