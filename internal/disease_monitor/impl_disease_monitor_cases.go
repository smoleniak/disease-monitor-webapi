package disease_monitor

import (
	"net/http"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implDiseaseMonitorCasesAPI struct {
}

func NewDiseaseMonitorCasesApi() DiseaseMonitorCasesAPI {
	return &implDiseaseMonitorCasesAPI{}
}

func (o implDiseaseMonitorCasesAPI) CreateDiseaseCaseListEntry(c *gin.Context) {
	updateRegionFunc(c, func(c *gin.Context, region *Region) (*Region, interface{}, int) {
		var entry DiseaseCaseEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.Disease.Code == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Disease code is required",
			}, http.StatusBadRequest
		}

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(region.DiseaseCases, func(existing DiseaseCaseEntry) bool {
			return entry.Id == existing.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		region.DiseaseCases = append(region.DiseaseCases, entry)
		entryIndx := slices.IndexFunc(region.DiseaseCases, func(existing DiseaseCaseEntry) bool {
			return entry.Id == existing.Id
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return region, region.DiseaseCases[entryIndx], http.StatusOK
	})
}

func (o implDiseaseMonitorCasesAPI) DeleteDiseaseCaseEntry(c *gin.Context) {
	updateRegionFunc(c, func(c *gin.Context, region *Region) (*Region, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(region.DiseaseCases, func(waiting DiseaseCaseEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		region.DiseaseCases = append(region.DiseaseCases[:entryIndx], region.DiseaseCases[entryIndx+1:]...)
		// region.reconcileDiseaseCases()
		return region, nil, http.StatusNoContent
	})
}

func (o implDiseaseMonitorCasesAPI) GetDiseaseCaseEntries(c *gin.Context) {
	updateRegionFunc(c, func(c *gin.Context, region *Region) (*Region, interface{}, int) {
		result := region.DiseaseCases
		if result == nil {
			result = []DiseaseCaseEntry{}
		}
		// return nil region - no need to update it in db
		return nil, result, http.StatusOK
	})
}

func (o implDiseaseMonitorCasesAPI) GetDiseaseCaseEntry(c *gin.Context) {
	updateRegionFunc(c, func(c *gin.Context, region *Region) (*Region, interface{}, int) {
		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(region.DiseaseCases, func(waiting DiseaseCaseEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		// return nil region - no need to update it in db
		return nil, region.DiseaseCases[entryIndx], http.StatusOK
	})
}

func (o implDiseaseMonitorCasesAPI) UpdateDiseaseCaseEntry(c *gin.Context) {
	updateRegionFunc(c, func(c *gin.Context, region *Region) (*Region, interface{}, int) {
		var entry DiseaseCaseEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		entryId := c.Param("entryId")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(region.DiseaseCases, func(waiting DiseaseCaseEntry) bool {
			return entryId == waiting.Id
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		if entry.Id != "" {
			region.DiseaseCases[entryIndx].Id = entry.Id
		}

		if entry.Disease.Code != "" {
			region.DiseaseCases[entryIndx].Disease.Code = entry.Disease.Code
			region.DiseaseCases[entryIndx].Disease.Value = entry.Disease.Value
		}

		if entry.Patient.Name != "" {
			region.DiseaseCases[entryIndx].Patient.Name = entry.Patient.Name
		}

		// disease cannot start in the future
		if entry.DiseaseStart.Before(time.Time{}) {
			region.DiseaseCases[entryIndx].DiseaseStart = entry.DiseaseStart
		}

		if entry.DiseaseEnd.Before(time.Time{}) {
			region.DiseaseCases[entryIndx].DiseaseEnd = entry.DiseaseEnd
		}

		if entry.Latitude > 0 {
			region.DiseaseCases[entryIndx].Latitude = entry.Latitude
		}

		if entry.Longtitude > 0 {
			region.DiseaseCases[entryIndx].Longtitude = entry.Longtitude
		}

		// region.reconcileDiseaseCases()
		return region, region.DiseaseCases[entryIndx], http.StatusOK
	})
}
