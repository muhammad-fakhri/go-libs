package httpmiddleware

const (
	ExcludeLog = true
	IncludeLog = false
)

const (
	headerNameUserID    = "x-user-id"
	headerNameRequestID = "x-request-id"
	headerNameTenant    = "x-tenant"
	headerNameCountry   = "x-country"

	ContextUserIdKey  = "user_id"
	ContextEventIDKey = "event_id"

	EventPrefix  = "events"
	URLSeparator = "/"
)

const (
	wipedMessage = "-"
)
