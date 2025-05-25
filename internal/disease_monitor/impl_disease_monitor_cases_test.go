package disease_monitor

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smoleniak/disease-monitor-webapi/internal/db_service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DiseaseMonitorCasesSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Region]
}

func TestDiseaseMonitorCasesSuite(t *testing.T) {
	suite.Run(t, new(DiseaseMonitorCasesSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *DiseaseMonitorCasesSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Region]{}

	// Compile time Assert that the mock is of type db_service.DbService[Region]
	var _ db_service.DbService[Region] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Region{
				Id: "test-region",
				DiseaseCases: []DiseaseCaseEntry{
					{
						Id: "test-entry",
						Disease: Disease{
							Code:  "covid",
							Value: "SARS-CoV-19",
						},
						Patient: Patient{
							Id:   "a123",
							Name: "Joe Test",
						},
						DiseaseStart: time.Now(),
						Latitude:     48.45,
						Longtitude:   17.32,
					},
				},
			},
			nil,
		)
}

func (suite *DiseaseMonitorCasesSuite) Test_UpdateCaseList_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
        "id": "test-entry",
        "patientId": "test-patient",
        "estimatedDurationMinutes": 42
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "regionId", Value: "test-region"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("POST", "/region/test-region/waitinglist/test-entry", strings.NewReader(json))

	sut := implDiseaseMonitorCasesAPI{}

	// ACT
	sut.UpdateDiseaseCaseEntry(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-region", mock.Anything)
}
