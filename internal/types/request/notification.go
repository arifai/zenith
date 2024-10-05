package request

type (
	// NotificationMarkAsReadRequest represents a request to mark a notification as read.
	NotificationMarkAsReadRequest struct {
		ID string `json:"id" validate:"required,uuid" reason:"required:ID is required;uuid:ID must be a valid UUID"`
	}
)
