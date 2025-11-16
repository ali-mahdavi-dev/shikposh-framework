package http

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strings"

	"github.com/ali-mahdavi-dev/framework/errors"
	"github.com/ali-mahdavi-dev/framework/errors/phrases"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cast"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// formatValidationError formats validator errors into a readable message
func formatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldError := range validationErrors {
			message := fmt.Sprintf("field '%s' failed validation: %s", fieldError.Field(), getValidationMessage(fieldError))
			messages = append(messages, message)
		}
		return strings.Join(messages, "; ")
	}
	return err.Error()
}

// getValidationMessage returns a user-friendly message for validation errors
func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fieldError.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fieldError.Param())
	case "email":
		return "must be a valid email address"
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fieldError.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fieldError.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fieldError.Param())
	default:
		return fmt.Sprintf("failed validation rule '%s'", fieldError.Tag())
	}
}

var statusMap = map[string]int{
	http.StatusText(http.StatusBadRequest):            http.StatusBadRequest,
	http.StatusText(http.StatusUnauthorized):          http.StatusUnauthorized,
	http.StatusText(http.StatusPaymentRequired):       http.StatusPaymentRequired,
	http.StatusText(http.StatusForbidden):             http.StatusForbidden,
	http.StatusText(http.StatusNotFound):              http.StatusNotFound,
	http.StatusText(http.StatusMethodNotAllowed):      http.StatusMethodNotAllowed,
	http.StatusText(http.StatusConflict):              http.StatusConflict,
	http.StatusText(http.StatusRequestEntityTooLarge): http.StatusRequestEntityTooLarge,
	http.StatusText(http.StatusRequestTimeout):        http.StatusRequestTimeout,
	http.StatusText(http.StatusTooManyRequests):       http.StatusTooManyRequests,
	http.StatusText(http.StatusInternalServerError):   http.StatusInternalServerError,
}

// statusTextToCode converts HTTP status text to status code
func statusTextToCode(statusText string) int {
	if code, ok := statusMap[statusText]; ok {
		return code
	}
	return http.StatusInternalServerError
}

// Token
func GetToken(c fiber.Ctx) string {
	auth := c.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = auth
	}

	if token == "" {
		token = c.Query("accessToken")
	}

	return token
}

// Parsing
func ParseJSON(c fiber.Ctx, obj interface{}) error {
	// Parse JSON body
	if err := c.Bind().Body(obj); err != nil {
		return errors.Validation(phrases.FailedParseJson, err.Error())
	}

	// Validate struct
	if err := validate.Struct(obj); err != nil {
		return errors.Validation(phrases.FailedParseJson, formatValidationError(err))
	}

	return nil
}

func ParseQuery(c fiber.Ctx, obj interface{}) error {
	// Parse query parameters
	if err := c.Bind().Query(obj); err != nil {
		return errors.Validation(phrases.FailedParseQuery, err.Error())
	}

	// Validate struct
	if err := validate.Struct(obj); err != nil {
		return errors.Validation(phrases.FailedParseQuery, formatValidationError(err))
	}

	return nil
}

func ParsePaginationQueryParam(c fiber.Ctx, obj *PaginationResult) error {
	if err := c.Bind().Query(obj); err != nil {
		return errors.Validation(phrases.FailedParseQuery, err.Error())
	}
	if obj.Limit < 1 {
		obj.Limit = 10
	}
	return nil
}

func ParseForm(c fiber.Ctx, obj interface{}) error {
	// Parse form data
	if err := c.Bind().Body(obj); err != nil {
		return errors.Validation(phrases.FailedParseForm, err.Error())
	}

	// Validate struct
	if err := validate.Struct(obj); err != nil {
		return errors.Validation(phrases.FailedParseForm, formatValidationError(err))
	}

	return nil
}

// Responses
func ResJSON(c fiber.Ctx, status int, v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Set(ResBodyKey, string(buf))
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Status(status)
	return c.Send(buf)
}

func ResSuccess(c fiber.Ctx, v interface{}) error {
	return ResJSON(c, fiber.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
	})
}

func ResOK(c fiber.Ctx) error {
	return ResJSON(c, fiber.StatusOK, ResponseResult{
		Success: true,
	})
}

func CalculatePagination(total, limit, skip int64) (int64, int64) {
	if limit <= 0 {
		limit = 1 // Prevent division by zero
	}
	pages := int64(math.Ceil(float64(total) / float64(limit)))
	page := (skip / limit) + 1
	return pages, page
}

func ResPage(c fiber.Ctx, v interface{}, pr *PaginationResult) error {
	var total, pages, page int64
	if pr != nil {
		total = pr.Total
		pages, page = CalculatePagination(total, pr.Limit, pr.Skip)
	}
	if page < 1 {
		page = 1
	}
	if pages < 1 {
		pages = 1
	}
	if page > pages {
		page = pages
	}

	reflectValue := reflect.Indirect(reflect.ValueOf(v))
	if reflectValue.IsNil() {
		v = make([]interface{}, 0)
	}

	return ResJSON(c, fiber.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
		Total:   total,
		Page:    page,
		Pages:   pages,
	})
}

func ResError(c fiber.Ctx, err error) error {
	var httpErr Error
	var statusCode int

	// Check if it's an app error (errors.Error)
	if appErr, ok := errors.As(err); ok {
		httpErr = ErrorToHTTP(appErr)
		statusCode = errorTypeToHTTPStatus(appErr.Type())
	} else if e, ok := err.(Error); ok {
		// Already an HTTP error - need to get status code from status text
		httpErr = e
		statusCode = statusTextToCode(httpErr.Status())
	} else {
		// Convert to internal error
		appErr := errors.Internal(cast.ToString(err))
		httpErr = ErrorToHTTP(appErr)
		statusCode = errorTypeToHTTPStatus(appErr.Type())
	}

	// Convert Error interface to HTTPError struct for JSON serialization
	httpErrorStruct := ToHTTPError(httpErr)

	// Add request_id from context if available
	// Import middleware package to get request_id
	requestID := getRequestIDFromContext(c)
	if requestID != "" {
		httpErrorStruct.RequestID = requestID
	}

	return ResJSON(c, statusCode, ResponseResult{Error: httpErrorStruct})
}

// getRequestIDFromContext extracts request ID from Fiber context
// This avoids circular dependency with middleware package
func getRequestIDFromContext(c fiber.Ctx) string {
	if requestID, ok := c.Locals("request_id").(string); ok {
		return requestID
	}
	return ""
}
