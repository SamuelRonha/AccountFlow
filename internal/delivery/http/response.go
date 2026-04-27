package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"AccountFlow/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// validationTag maps validator tag names to human-readable messages.
var validationTag = map[string]string{
	"required": "is required",
	"min":      "must be greater than or equal to the minimum",
	"ne":       "must not be zero",
	"gt":       "must be greater than 0",
	"gte":      "must be greater than or equal to 0",
	"email":    "must be a valid e-mail",
	"uuid":     "must be a valid UUID",
}

// bindError translates Gin/validator binding errors into clean field-level messages.
func bindError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := make([]string, 0, len(ve))
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			hint, ok := validationTag[fe.Tag()]
			if !ok {
				hint = fmt.Sprintf("failed validation '%s'", fe.Tag())
			}
			msgs = append(msgs, fmt.Sprintf("'%s' %s", field, hint))
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": msgs})
		return
	}
	// JSON syntax / type mismatch
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
}

func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		return http.StatusNotFound, "account not found"
	case errors.Is(err, domain.ErrDocumentAlreadyUsed):
		return http.StatusConflict, "document already registered"
	case errors.Is(err, domain.ErrOperationTypeNotFound):
		return http.StatusUnprocessableEntity, "operation type not found"
	case errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusBadRequest, "amount must not be zero"
	case errors.Is(err, domain.ErrInvalidField):
		return http.StatusBadRequest, "one or more fields are invalid"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func writeError(c *gin.Context, err error) {
	code, msg := mapDomainError(err)
	c.JSON(code, gin.H{"error": msg})
}
