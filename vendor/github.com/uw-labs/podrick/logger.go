package podrick

// Logger must be implemented to log events. See
// https://logur.dev/logur for some adapters
// for popular logging libraries.
type Logger interface {
	Trace(msg string, fields ...map[string]interface{})
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
}
