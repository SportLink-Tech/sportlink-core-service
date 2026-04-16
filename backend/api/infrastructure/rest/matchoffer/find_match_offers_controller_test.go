package matchoffer_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchoffer/usecases"
	"sportlink/api/domain/common"
	domain "sportlink/api/domain/matchoffer"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/matchoffer"
	amocks "sportlink/mocks/api/application"
	pmocks "sportlink/mocks/api/infrastructure/rest/matchoffer/parser"
)

type SearchUseCaseMock = amocks.UseCase[usecases.SearchMatchOffersInput, usecases.FindMatchOfferResult]

func TestSearchMatchOffers(t *testing.T) {
	validator := validator.New()

	testCases := []struct {
		name        string
		accountID   string
		queryParams map[string]string
		on          func(t *testing.T, ucMock *SearchUseCaseMock, parserMock *pmocks.QueryParser)
		then        func(t *testing.T, responseCode int, response interface{})
	}{
		{
			name:      "given valid query parameters when searching then returns available offers",
			accountID: "account-1",
			queryParams: map[string]string{
				"sports":    "Paddle",
				"from_date": "2025-12-01",
			},
			on: func(t *testing.T, ucMock *SearchUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "Paddle").Return([]common.Sport{common.Paddle}, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "2025-12-01").Return(time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC), nil)
				parserMock.On("Date", "").Return(time.Time{}, nil)
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("GeoFilter", "", "", "").Return(nil, nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				expectedResult := &usecases.FindMatchOfferResult{
					Entities: []domain.Entity{createTestOffer("Boca", common.Paddle, domain.StatusPending)},
					Page:     usecases.PageInfo{Number: 1, OutOf: 1, Total: 1},
				}

				ucMock.On("Invoke", mock.Anything, mock.MatchedBy(func(input usecases.SearchMatchOffersInput) bool {
					return input.ViewerAccountID == "account-1" &&
						len(input.Query.Sports) == 1 &&
						input.Query.Sports[0] == common.Paddle &&
						!input.Query.FromDate.IsZero()
				})).Return(expectedResult, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				responseMap := response.(map[string]interface{})
				data := responseMap["data"].([]interface{})
				assert.Len(t, data, 1)
				assert.Equal(t, "Boca", data[0].(map[string]interface{})["team_name"])
				pagination := responseMap["pagination"].(map[string]interface{})
				assert.Equal(t, float64(1), pagination["number"])
				assert.Equal(t, float64(1), pagination["out_of"])
			},
		},
		{
			name:        "given no query parameters when searching then returns all available offers",
			accountID:   "account-2",
			queryParams: map[string]string{},
			on: func(t *testing.T, ucMock *SearchUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "").Return(nil, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "").Return(time.Time{}, nil).Twice()
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("GeoFilter", "", "", "").Return(nil, nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				expectedResult := &usecases.FindMatchOfferResult{
					Entities: []domain.Entity{
						createTestOffer("Boca", common.Paddle, domain.StatusPending),
						createTestOffer("River", common.Football, domain.StatusPending),
					},
					Page: usecases.PageInfo{Number: 1, OutOf: 1, Total: 2},
				}

				ucMock.On("Invoke", mock.Anything, mock.MatchedBy(func(input usecases.SearchMatchOffersInput) bool {
					return input.ViewerAccountID == "account-2"
				})).Return(expectedResult, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				responseMap := response.(map[string]interface{})
				data := responseMap["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:      "given invalid category format when searching then returns validation error",
			accountID: "account-3",
			queryParams: map[string]string{
				"categories": "invalid",
			},
			on: func(t *testing.T, ucMock *SearchUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "").Return(nil, nil)
				parserMock.On("Categories", "invalid").Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Equal(t, "request_validation_failed", response.(map[string]interface{})["code"])
			},
		},
		{
			name:      "given use case returns error when searching then returns error",
			accountID: "account-4",
			queryParams: map[string]string{
				"sports": "Paddle",
			},
			on: func(t *testing.T, ucMock *SearchUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "Paddle").Return([]common.Sport{common.Paddle}, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "").Return(time.Time{}, nil).Twice()
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("GeoFilter", "", "", "").Return(nil, nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				ucMock.On("Invoke", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response.(map[string]interface{})["code"])
			},
		},
		{
			name:      "given no available offers when searching then returns not found",
			accountID: "account-5",
			queryParams: map[string]string{
				"sports": "Paddle",
			},
			on: func(t *testing.T, ucMock *SearchUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "Paddle").Return([]common.Sport{common.Paddle}, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "").Return(time.Time{}, nil).Twice()
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("GeoFilter", "", "", "").Return(nil, nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				ucMock.On("Invoke", mock.Anything, mock.Anything).Return(&usecases.FindMatchOfferResult{
					Entities: []domain.Entity{},
					Page:     usecases.PageInfo{Number: 1, OutOf: 0, Total: 0},
				}, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusNotFound, responseCode)
				assert.Equal(t, "not_found", response.(map[string]interface{})["code"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ucMock := amocks.NewUseCase[usecases.SearchMatchOffersInput, usecases.FindMatchOfferResult](t)
			parserMock := pmocks.NewQueryParser(t)
			controller := matchoffer.NewControllerWithParser(nil, nil, ucMock, nil, nil, nil, validator, parserMock)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.GET("/account/:account_id/match-offer/search", controller.SearchMatchOffers)

			tc.on(t, ucMock, parserMock)

			req := httptest.NewRequest("GET", "/account/"+tc.accountID+"/match-offer/search", nil)
			q := req.URL.Query()
			for key, value := range tc.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			tc.then(t, resp.Code, createResponse(resp))
		})
	}
}

// Helper functions shared across matchoffer controller tests.

func createTestOffer(teamName string, sport common.Sport, status domain.Status) domain.Entity {
	timeSlot, _ := domain.NewTimeSlot(
		time.Date(2025, 12, 10, 18, 0, 0, 0, time.UTC),
		time.Date(2025, 12, 10, 20, 0, 0, 0, time.UTC),
	)
	location := domain.NewLocation("Argentina", "Buenos Aires", "Palermo")
	categoryRange := domain.NewGreaterThanCategory(common.L4)

	return domain.NewMatchOffer(
		teamName,
		sport,
		time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		timeSlot,
		location,
		categoryRange,
		status,
		time.Now(),
		"", 0,
	)
}

func createResponse(resp *httptest.ResponseRecorder) interface{} {
	var response interface{}
	if resp.Body.Len() > 0 {
		json.Unmarshal(resp.Body.Bytes(), &response)
	}
	return response
}
