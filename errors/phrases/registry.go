package phrases

import "sync"

// Registry holds error phrases for different services
type Registry struct {
	mu      sync.RWMutex
	phrases map[Language]map[MessagePhrase]string
}

var globalRegistry *Registry
var once sync.Once

// GetRegistry returns the global error phrase registry
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			phrases: make(map[Language]map[MessagePhrase]string),
		}
		// Initialize with default messages
		globalRegistry.initDefaults()
	})
	return globalRegistry
}

// initDefaults initializes the registry with default framework messages
func (r *Registry) initDefaults() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.phrases[Fa] = map[MessagePhrase]string{
		DefaultValidationID:       "درخواست نامعتبر است",
		DefaultNotFoundID:         "منبع یافت نشد",
		DefaultConflictID:          "درخواست با وضعیت فعلی در تضاد است",
		DefaultUnauthorizedID:     "احراز هویت نشده است",
		DefaultForbiddenID:         "دسترسی مجاز نیست",
		DefaultTimeoutID:          "درخواست منقضی شده است",
		DefaultRateLimitID:        "تعداد درخواست بیش از حد مجاز است",
		DefaultTooLargeID:         "حجم درخواست بیش از حد مجاز است",
		DefaultInternalID:         "خطای داخلی سرور",
		DefaultMethodNotAllowedID:  "متد HTTP مجاز نیست",
		OperationCanNot:           "عملیات موفق آمیز نبود. لطفا دوباره تلاش بفرمایید",
		FailedParseJson:           "خطا در تجزیه JSON: %s",
		FailedParseQuery:          "خطا در تجزیه Query: %s",
		FailedParseForm:           "خطا در تجزیه Form: %s",
	}

	r.phrases[En] = map[MessagePhrase]string{
		DefaultValidationID:       "Invalid request",
		DefaultNotFoundID:         "Resource not found",
		DefaultConflictID:          "Request conflicts with current state",
		DefaultUnauthorizedID:     "Unauthorized",
		DefaultForbiddenID:         "Forbidden",
		DefaultTimeoutID:          "Request timeout",
		DefaultRateLimitID:        "Too many requests",
		DefaultTooLargeID:         "Request entity too large",
		DefaultInternalID:         "Internal server error",
		DefaultMethodNotAllowedID:  "Method not allowed",
		OperationCanNot:           "Operation was not successful. Please try again",
		FailedParseJson:           "Failed to parse json: %s",
		FailedParseQuery:          "Failed to parse query: %s",
		FailedParseForm:           "Failed to parse form: %s",
	}
}

// Register registers error phrases for a service
// phrases should be a map of language to map of MessagePhrase to message string
func (r *Registry) Register(phrases map[Language]map[MessagePhrase]string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for lang, langPhrases := range phrases {
		if r.phrases[lang] == nil {
			r.phrases[lang] = make(map[MessagePhrase]string)
		}
		for phrase, message := range langPhrases {
			r.phrases[lang][phrase] = message
		}
	}
}

// GetMessage retrieves a message for a given phrase and language
func (r *Registry) GetMessage(phrase MessagePhrase, lan Language) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if lan == "" {
		lan = Fa
	}

	// First, try to find message in requested language
	if langPhrases, ok := r.phrases[lan]; ok {
		if msg, found := langPhrases[phrase]; found {
			return msg
		}
	}

	// If not found in requested language, search in all other languages
	for lang, langPhrases := range r.phrases {
		if lang != lan {
			if msg, found := langPhrases[phrase]; found {
				return msg
			}
		}
	}

	// If not found in any language, return default message
	return "Unknown error"
}

