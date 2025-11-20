package phrases

type Language string

const (
	Fa Language = "fa"
	En Language = "en"
)

// GetMessage retrieves a message from the global registry
func GetMessage(phrase MessagePhrase, lan Language) string {
	return GetRegistry().GetMessage(phrase, lan)
}
