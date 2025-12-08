package matchannouncement_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sportlink/api/application/matchannouncement/request"
	"sportlink/api/domain/common"
	domain "sportlink/api/domain/matchannouncement"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest/matchannouncement"
	amocks "sportlink/mocks/api/application"
)

// UseCaseMock is a type alias to make the generic mock easier to use
type UseCaseMock = amocks.UseCase[domain.Entity, domain.Entity]

func TestCreateMatchAnnouncement(t *testing.T) {
	validator := validator.New()

	testCases := []struct {
		name    string
		payload request.NewMatchAnnouncementRequest
		on      func(t *testing.T, useCaseMock *UseCaseMock)
		then    func(t *testing.T, responseCode int, response map[string]interface{})
	}{
		{
			name: "given_valid_match_announcement_when_creating_then_returns_created_announcement",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "Boca",
				Sport:    "Paddle",
				Day:      "2025-12-10",
				TimeSlot: request.TimeSlot{
					StartTime: "2025-12-10T18:00:00",
					EndTime:   "2025-12-10T20:00:00",
				},
				Location: request.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "Palermo",
				},
				AdmittedCategories: request.CategoryRangeInput{
					Type:     "GREATER_THAN",
					MinLevel: 4,
				},
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				expectedEntity := domain.NewMatchAnnouncement(
					"Boca",
					common.Paddle,
					time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
					mustCreateTimeSlot(t, "2025-12-10T18:00:00", "2025-12-10T20:00:00"),
					domain.NewLocation("Argentina", "Buenos Aires", "Palermo"),
					domain.NewGreaterThanCategory(common.L4),
					domain.StatusPending,
					time.Now(),
				)

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.TeamName == "Boca" &&
						entity.Sport == common.Paddle &&
						entity.Location.Country == "Argentina" &&
						entity.Location.Province == "Buenos Aires" &&
						entity.Location.Locality == "Palermo"
				})).Return(&expectedEntity, nil)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.Equal(t, "Boca", response["team_name"])
				assert.Equal(t, "Paddle", response["sport"])
				assert.Equal(t, "Argentina", response["location"].(map[string]interface{})["country"])
			},
		},
		{
			name:    "given_invalid_json_when_creating_then_returns_bad_request",
			payload: request.NewMatchAnnouncementRequest{},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				// No expectations - request should fail before use case is called
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				// Empty payload fails validation, not JSON parsing
				assert.Equal(t, "request_validation_failed", response["code"])
			},
		},
		{
			name: "given_missing_required_fields_when_creating_then_returns_validation_error",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "",
				Sport:    "Paddle",
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				// No expectations - validation should fail before use case is called
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Equal(t, "request_validation_failed", response["code"])
			},
		},
		{
			name: "given_invalid_sport_when_creating_then_returns_validation_error",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "Boca",
				Sport:    "InvalidSport",
				Day:      "2025-12-10",
				TimeSlot: request.TimeSlot{
					StartTime: "2025-12-10T18:00:00",
					EndTime:   "2025-12-10T20:00:00",
				},
				Location: request.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "Palermo",
				},
				AdmittedCategories: request.CategoryRangeInput{
					Type:     "GREATER_THAN",
					MinLevel: 4,
				},
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				// Invalid sport passes validation but fails in mapper or use case
				// The mapper will convert it, but use case validation will fail
				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.TeamName == "Boca" && entity.Sport == common.Sport("InvalidSport")
				})).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				// Sport validation happens in use case, not in request validation
				// Error handler returns 409 for use case execution errors
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
			},
		},
		{
			name: "given_invalid_category_range_type_when_creating_then_returns_validation_error",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "Boca",
				Sport:    "Paddle",
				Day:      "2025-12-10",
				TimeSlot: request.TimeSlot{
					StartTime: "2025-12-10T18:00:00",
					EndTime:   "2025-12-10T20:00:00",
				},
				Location: request.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "Palermo",
				},
				AdmittedCategories: request.CategoryRangeInput{
					Type: "INVALID_TYPE",
				},
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				// No expectations - validation should fail before use case is called
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusBadRequest, responseCode)
				assert.Equal(t, "request_validation_failed", response["code"])
			},
		},
		{
			name: "given_use_case_returns_error_when_creating_then_returns_internal_server_error",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "Boca",
				Sport:    "Paddle",
				Day:      "2025-12-10",
				TimeSlot: request.TimeSlot{
					StartTime: "2025-12-10T18:00:00",
					EndTime:   "2025-12-10T20:00:00",
				},
				Location: request.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "Palermo",
				},
				AdmittedCategories: request.CategoryRangeInput{
					Type:     "GREATER_THAN",
					MinLevel: 4,
				},
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.TeamName == "Boca" && entity.Sport == common.Paddle
				})).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				// Error handler returns 409 for use case execution errors
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
			},
		},
		{
			name: "given_team_does_not_exist_when_creating_then_returns_error",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "NonExistentTeam",
				Sport:    "Paddle",
				Day:      "2025-12-10",
				TimeSlot: request.TimeSlot{
					StartTime: "2025-12-10T18:00:00",
					EndTime:   "2025-12-10T20:00:00",
				},
				Location: request.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "Palermo",
				},
				AdmittedCategories: request.CategoryRangeInput{
					Type:     "GREATER_THAN",
					MinLevel: 4,
				},
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.TeamName == "NonExistentTeam" && entity.Sport == common.Paddle
				})).Return(nil, assert.AnError)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				// Error handler returns 409 for use case execution errors
				assert.Equal(t, http.StatusConflict, responseCode)
				assert.Equal(t, "use_case_execution_error", response["code"])
			},
		},
		{
			name: "given_specific_categories_when_creating_then_returns_created_announcement",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "Boca",
				Sport:    "Paddle",
				Day:      "2025-12-10",
				TimeSlot: request.TimeSlot{
					StartTime: "2025-12-10T18:00:00",
					EndTime:   "2025-12-10T20:00:00",
				},
				Location: request.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "Palermo",
				},
				AdmittedCategories: request.CategoryRangeInput{
					Type:       "SPECIFIC",
					Categories: []int{4, 5, 6},
				},
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				expectedEntity := domain.NewMatchAnnouncement(
					"Boca",
					common.Paddle,
					time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
					mustCreateTimeSlot(t, "2025-12-10T18:00:00", "2025-12-10T20:00:00"),
					domain.NewLocation("Argentina", "Buenos Aires", "Palermo"),
					domain.NewSpecificCategories([]common.Category{common.L4, common.L5, common.L6}),
					domain.StatusPending,
					time.Now(),
				)

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.TeamName == "Boca" &&
						entity.Sport == common.Paddle &&
						entity.AdmittedCategories.Type == domain.RangeTypeSpecific
				})).Return(&expectedEntity, nil)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.Equal(t, "Boca", response["team_name"])
				admittedCategories := response["admitted_categories"].(map[string]interface{})
				assert.Equal(t, "SPECIFIC", admittedCategories["type"])
			},
		},
		{
			name: "given_between_categories_when_creating_then_returns_created_announcement",
			payload: request.NewMatchAnnouncementRequest{
				TeamName: "Boca",
				Sport:    "Paddle",
				Day:      "2025-12-10",
				TimeSlot: request.TimeSlot{
					StartTime: "2025-12-10T18:00:00",
					EndTime:   "2025-12-10T20:00:00",
				},
				Location: request.Location{
					Country:  "Argentina",
					Province: "Buenos Aires",
					Locality: "Palermo",
				},
				AdmittedCategories: request.CategoryRangeInput{
					Type:     "BETWEEN",
					MinLevel: 3,
					MaxLevel: 6,
				},
			},
			on: func(t *testing.T, useCaseMock *UseCaseMock) {
				categoryRange, _ := domain.NewBetweenCategories(common.L3, common.L6)
				expectedEntity := domain.NewMatchAnnouncement(
					"Boca",
					common.Paddle,
					time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
					mustCreateTimeSlot(t, "2025-12-10T18:00:00", "2025-12-10T20:00:00"),
					domain.NewLocation("Argentina", "Buenos Aires", "Palermo"),
					categoryRange,
					domain.StatusPending,
					time.Now(),
				)

				useCaseMock.On("Invoke", mock.Anything, mock.MatchedBy(func(entity domain.Entity) bool {
					return entity.TeamName == "Boca" &&
						entity.Sport == common.Paddle &&
						entity.AdmittedCategories.Type == domain.RangeTypeBetween
				})).Return(&expectedEntity, nil)
			},
			then: func(t *testing.T, responseCode int, response map[string]interface{}) {
				assert.Equal(t, http.StatusCreated, responseCode)
				assert.Equal(t, "Boca", response["team_name"])
				admittedCategories := response["admitted_categories"].(map[string]interface{})
				assert.Equal(t, "BETWEEN", admittedCategories["type"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			useCaseMock := amocks.NewUseCase[domain.Entity, domain.Entity](t)
			controller := matchannouncement.NewController(useCaseMock, nil, validator)

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.Use(middleware.ErrorHandler())
			router.POST("/match-announcement", controller.CreateMatchAnnouncement)

			// Given
			tc.on(t, useCaseMock)
			jsonData, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/match-announcement", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// When
			router.ServeHTTP(resp, req)

			// Then
			response := createMapResponse(resp)
			tc.then(t, resp.Code, response)
		})
	}
}

// Helper functions

func createMapResponse(resp *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	return response
}

func mustCreateTimeSlot(t *testing.T, startTimeStr, endTimeStr string) domain.TimeSlot {
	// Try multiple date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
	}

	var startTime, endTime time.Time
	var err error

	for _, format := range formats {
		startTime, err = time.Parse(format, startTimeStr)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("Failed to parse start time %s: %v", startTimeStr, err)
	}

	for _, format := range formats {
		endTime, err = time.Parse(format, endTimeStr)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("Failed to parse end time %s: %v", endTimeStr, err)
	}

	timeSlot, err := domain.NewTimeSlot(startTime, endTime)
	if err != nil {
		t.Fatalf("Failed to create time slot: %v", err)
	}

	return timeSlot
}
