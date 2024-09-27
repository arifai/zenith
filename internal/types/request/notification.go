package request

import "github.com/google/uuid"

type (
	// NotificationMarkAsReadRequest FIXME: Handle if the UUID has empty string.
	NotificationMarkAsReadRequest struct {
		ID uuid.UUID `json:"id" validate:"required,uuid"`
	}
)
