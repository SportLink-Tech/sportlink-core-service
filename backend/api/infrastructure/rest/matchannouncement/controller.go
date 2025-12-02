package matchannouncement

import (
	"sportlink/api/application"
	"sportlink/api/domain/matchannouncement"
	"sportlink/api/infrastructure/rest/matchannouncement/parser"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller interface {
	CreateMatchAnnouncement(c *gin.Context)
	FindMatchAnnouncements(c *gin.Context)
}

type DefaultController struct {
	createMatchAnnouncementUC application.UseCase[matchannouncement.Entity, matchannouncement.Entity]
	findMatchAnnouncementsUC  application.UseCase[matchannouncement.DomainQuery, []matchannouncement.Entity]
	validator                 *validator.Validate
	queryParser               parser.QueryParser
}

func NewController(
	createMatchAnnouncementUC application.UseCase[matchannouncement.Entity, matchannouncement.Entity],
	findMatchAnnouncementsUC application.UseCase[matchannouncement.DomainQuery, []matchannouncement.Entity],
	validator *validator.Validate,
) Controller {
	return NewControllerWithParser(createMatchAnnouncementUC, findMatchAnnouncementsUC, validator, nil)
}

// NewControllerWithParser creates a controller with an optional query parser.
// If queryParser is nil, a new DefaultQueryParser will be created.
// This allows for dependency injection in tests.
func NewControllerWithParser(
	createMatchAnnouncementUC application.UseCase[matchannouncement.Entity, matchannouncement.Entity],
	findMatchAnnouncementsUC application.UseCase[matchannouncement.DomainQuery, []matchannouncement.Entity],
	validator *validator.Validate,
	queryParser parser.QueryParser,
) Controller {
	if queryParser == nil {
		queryParser = parser.NewQueryParser()
	}
	return &DefaultController{
		createMatchAnnouncementUC: createMatchAnnouncementUC,
		findMatchAnnouncementsUC:  findMatchAnnouncementsUC,
		validator:                 validator,
		queryParser:               queryParser,
	}
}
