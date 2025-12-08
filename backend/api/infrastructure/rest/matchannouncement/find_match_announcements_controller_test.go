package matchannouncement_test

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

	"sportlink/api/application/matchannouncement/usecases"
	"sportlink/api/domain/common"
	domain "sportlink/api/domain/matchannouncement"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/matchannouncement"
	amocks "sportlink/mocks/api/application"
	pmocks "sportlink/mocks/api/infrastructure/rest/matchannouncement/parser"
)

// FindUseCaseMock is a type alias for the find match announcements use case mock
type FindUseCaseMock = amocks.UseCase[domain.DomainQuery, usecases.FindMatchAnnouncementResult]

func TestFindMatchAnnouncements(t *testing.T) {
	validator := validator.New()

	testCases := []struct {
		name        string
		queryParams map[string]string
		on          func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser)
		then        func(t *testing.T, responseCode int, response interface{})
	}{
		{
			name: "given valid query parameters when finding announcements then returns list of announcements",
			queryParams: map[string]string{
				"sports":     "Paddle",
				"categories": "4,5",
				"statuses":   "PENDING",
				"fromDate":   "2025-12-01",
			},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "Paddle").Return([]common.Sport{common.Paddle}, nil)
				parserMock.On("Categories", "4,5").Return([]common.Category{common.L4, common.L5}, nil)
				parserMock.On("Statuses", "PENDING").Return([]domain.Status{domain.StatusPending}, nil)
				parserMock.On("Date", "2025-12-01").Return(time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC), nil)
				parserMock.On("Date", "").Return(time.Time{}, nil)
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				expectedAnnouncement := createTestAnnouncement("Boca", common.Paddle, domain.StatusPending)
				expectedResult := &usecases.FindMatchAnnouncementResult{
					Entities: []domain.Entity{expectedAnnouncement},
					Page: usecases.PageInfo{
						Number: 1,
						OutOf:  1,
						Total:  1,
					},
				}

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(query domain.DomainQuery) bool {
					return len(query.Sports) == 1 &&
						query.Sports[0] == common.Paddle &&
						len(query.Categories) == 2 &&
						query.Categories[0] == common.L4 &&
						query.Categories[1] == common.L5 &&
						len(query.Statuses) == 1 &&
						query.Statuses[0] == domain.StatusPending &&
						!query.FromDate.IsZero()
				})).Return(expectedResult, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				responseMap := response.(map[string]interface{})
				data := responseMap["data"].([]interface{})
				assert.Greater(t, len(data), 0)
				firstAnnouncement := data[0].(map[string]interface{})
				assert.Equal(t, "Boca", firstAnnouncement["team_name"])
				assert.Equal(t, "Paddle", firstAnnouncement["sport"])
				pagination := responseMap["pagination"].(map[string]interface{})
				assert.Equal(t, float64(1), pagination["number"])
				assert.Equal(t, float64(1), pagination["out_of"])
			},
		},
		{
			name:        "given no query parameters when finding announcements then returns all announcements",
			queryParams: map[string]string{},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "").Return(nil, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "").Return(time.Time{}, nil).Twice()
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				expectedResult := &usecases.FindMatchAnnouncementResult{
					Entities: []domain.Entity{
						createTestAnnouncement("Boca", common.Paddle, domain.StatusPending),
						createTestAnnouncement("River", common.Football, domain.StatusConfirmed),
					},
					Page: usecases.PageInfo{
						Number: 1,
						OutOf:  1,
						Total:  2,
					},
				}

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(query domain.DomainQuery) bool {
					return len(query.Sports) == 0 &&
						len(query.Categories) == 0 &&
						len(query.Statuses) == 0 &&
						query.Location == nil
				})).Return(expectedResult, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				responseMap := response.(map[string]interface{})
				data := responseMap["data"].([]interface{})
				assert.Equal(t, 2, len(data))
			},
		},
		{
			name: "given invalid category format when finding announcements then returns validation error",
			queryParams: map[string]string{
				"categories": "invalid",
			},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "").Return(nil, nil)
				parserMock.On("Categories", "invalid").Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", responseMap["code"])
			},
		},
		{
			name: "given invalid date format when finding announcements then returns validation error",
			queryParams: map[string]string{
				"fromDate": "invalid-date",
			},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "").Return(nil, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "invalid-date").Return(time.Time{}, assert.AnError)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "request_validation_failed", responseMap["code"])
			},
		},
		{
			name: "given use case returns error when finding announcements then returns error",
			queryParams: map[string]string{
				"sports": "Paddle",
			},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "Paddle").Return([]common.Sport{common.Paddle}, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "").Return(time.Time{}, nil).Twice()
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(query domain.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Paddle
				})).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusConflict, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "use_case_execution_error", responseMap["code"])
			},
		},
		{
			name: "given no announcements found when finding announcements then returns not found",
			queryParams: map[string]string{
				"sports": "Paddle",
			},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "Paddle").Return([]common.Sport{common.Paddle}, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "").Return(time.Time{}, nil).Twice()
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				emptyResult := &usecases.FindMatchAnnouncementResult{
					Entities: []domain.Entity{},
					Page: usecases.PageInfo{
						Number: 1,
						OutOf:  0,
						Total:  0,
					},
				}
				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(query domain.DomainQuery) bool {
					return len(query.Sports) == 1 && query.Sports[0] == common.Paddle
				})).Return(emptyResult, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusNotFound, responseCode)
				responseMap := response.(map[string]interface{})
				assert.Equal(t, "not_found", responseMap["code"])
			},
		},
		{
			name: "given location parameters when finding announcements then returns filtered announcements",
			queryParams: map[string]string{
				"country":  "Argentina",
				"province": "Buenos Aires",
				"locality": "Palermo",
			},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "").Return(nil, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				parserMock.On("Date", "").Return(time.Time{}, nil).Twice()
				location := domain.NewLocation("Argentina", "Buenos Aires", "Palermo")
				parserMock.On("Location", "Argentina", "Buenos Aires", "Palermo").Return(&location)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				expectedResult := &usecases.FindMatchAnnouncementResult{
					Entities: []domain.Entity{
						createTestAnnouncement("Boca", common.Paddle, domain.StatusPending),
					},
					Page: usecases.PageInfo{
						Number: 1,
						OutOf:  1,
						Total:  1,
					},
				}

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(query domain.DomainQuery) bool {
					return query.Location != nil &&
						query.Location.Country == "Argentina" &&
						query.Location.Province == "Buenos Aires" &&
						query.Location.Locality == "Palermo"
				})).Return(expectedResult, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				responseMap := response.(map[string]interface{})
				data := responseMap["data"].([]interface{})
				assert.Equal(t, 1, len(data))
				firstAnnouncement := data[0].(map[string]interface{})
				location := firstAnnouncement["location"].(map[string]interface{})
				assert.Equal(t, "Argentina", location["country"])
				assert.Equal(t, "Buenos Aires", location["province"])
				assert.Equal(t, "Palermo", location["locality"])
			},
		},
		{
			name: "given date range when finding announcements then returns announcements in range",
			queryParams: map[string]string{
				"fromDate": "2025-12-01",
				"toDate":   "2025-12-31",
			},
			on: func(t *testing.T, useCaseMock *FindUseCaseMock, parserMock *pmocks.QueryParser) {
				parserMock.On("Sports", "").Return(nil, nil)
				parserMock.On("Categories", "").Return(nil, nil)
				parserMock.On("Statuses", "").Return(nil, nil)
				fromDate := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
				toDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
				parserMock.On("Date", "2025-12-01").Return(fromDate, nil)
				parserMock.On("Date", "2025-12-31").Return(toDate, nil)
				parserMock.On("Location", "", "", "").Return(nil)
				parserMock.On("Limit", "").Return(0, nil)
				parserMock.On("Offset", "").Return(0, nil)

				expectedResult := &usecases.FindMatchAnnouncementResult{
					Entities: []domain.Entity{
						createTestAnnouncement("Boca", common.Paddle, domain.StatusPending),
					},
					Page: usecases.PageInfo{
						Number: 1,
						OutOf:  1,
						Total:  1,
					},
				}

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(query domain.DomainQuery) bool {
					return !query.FromDate.IsZero() &&
						!query.ToDate.IsZero() &&
						query.FromDate.Equal(fromDate) &&
						query.ToDate.Equal(toDate)
				})).Return(expectedResult, nil)
			},
			then: func(t *testing.T, responseCode int, response interface{}) {
				assert.Equal(t, http.StatusOK, responseCode)
				responseMap := response.(map[string]interface{})
				data := responseMap["data"].([]interface{})
				assert.Equal(t, 1, len(data))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			useCaseMock := amocks.NewUseCase[domain.DomainQuery, usecases.FindMatchAnnouncementResult](t)
			parserMock := pmocks.NewQueryParser(t)
			controller := matchannouncement.NewControllerWithParser(nil, useCaseMock, validator, parserMock)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.GET("/match-announcement", controller.FindMatchAnnouncements)

			// Given
			tc.on(t, useCaseMock, parserMock)
			req := httptest.NewRequest("GET", "/match-announcement", nil)
			q := req.URL.Query()
			for key, value := range tc.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			resp := httptest.NewRecorder()

			// When
			router.ServeHTTP(resp, req)

			// Then
			response := createResponse(resp)
			tc.then(t, resp.Code, response)
		})
	}
}

// Helper functions

func createTestAnnouncement(teamName string, sport common.Sport, status domain.Status) domain.Entity {
	timeSlot, _ := domain.NewTimeSlot(
		time.Date(2025, 12, 10, 18, 0, 0, 0, time.UTC),
		time.Date(2025, 12, 10, 20, 0, 0, 0, time.UTC),
	)
	location := domain.NewLocation("Argentina", "Buenos Aires", "Palermo")
	categoryRange := domain.NewGreaterThanCategory(common.L4)

	return domain.NewMatchAnnouncement(
		teamName,
		sport,
		time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		timeSlot,
		location,
		categoryRange,
		status,
		time.Now(),
	)
}

func createResponse(resp *httptest.ResponseRecorder) interface{} {
	var response interface{}
	if resp.Body.Len() > 0 {
		json.Unmarshal(resp.Body.Bytes(), &response)
	}
	return response
}
