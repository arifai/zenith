package response

import "github.com/google/uuid"

type (
	// PushNotificationResponse represents a request to send a push notification.
	PushNotificationResponse struct {
		ID       uuid.UUID                     `json:"id" binding:"required"`
		Tokens   []string                      `json:"tokens" binding:"required"`
		Platform string                        `json:"platform" binding:"required"`
		Message  string                        `json:"message" binding:"required"`
		Title    string                        `json:"title" binding:"required"`
		Retry    int                           `json:"retry" binding:"required"`
		Data     *PushNotificationDataResponse `json:"data"`
	}

	// PushNotificationDataResponse represents the data required for sending a push notification.
	PushNotificationDataResponse struct {
		AccountID string   `json:"account_id" binding:"required"`
		IsRead    bool     `json:"is_read" binding:"required"`
		Tokens    []string `json:"tokens" binding:"required"`
		To        string   `json:"to" binding:"required"`
		Title     string   `json:"title" binding:"required"`
		Message   string   `json:"message" binding:"required"`
		Link      string   `json:"link"`
	}

	// FCMNotificationResponse represents a payload structure for sending notifications via FCM (Firebase Cloud Messaging).
	FCMNotificationResponse struct {
		Title       string `json:"title" binding:"required"`
		Body        string `json:"body" binding:"required"`
		ChannelID   string `json:"channel_id" binding:"required"`
		ClickAction string `json:"click_action" binding:"required"`
	}

	// PushNotificationsResponse represents a request to send multiple push notifications.
	PushNotificationsResponse struct {
		Notifications []PushNotificationResponse `json:"notifications" binding:"required"`
	}

	// AndroidNotificationResponse represents the response structure for sending Android notifications via a specific API.
	AndroidNotificationResponse struct {
		APIKey                string                   `json:"api_key" binding:"required"`
		To                    string                   `json:"to" binding:"required"`
		CollapseKey           string                   `json:"collapse_key" binding:"required"`
		DelayWhileIdle        bool                     `json:"delay_while_idle" binding:"required"`
		TimeToLive            int                      `json:"time_to_live" binding:"required"`
		RestrictedPackageName string                   `json:"restricted_package_name" binding:"required"`
		DryRun                bool                     `json:"dry_run" binding:"required"`
		Condition             string                   `json:"condition" binding:"required"`
		Notification          *FCMNotificationResponse `json:"notification" binding:"required"`
	}

	// IOSNotificationResponse represents the structure of a response for an iOS notification.
	IOSNotificationResponse struct {
		ApnsID     string                        `json:"apns_id" binding:"required"`
		CollapseID string                        `json:"collapse_id" binding:"required"`
		ThreadID   int                           `json:"thread_id" binding:"required"`
		URLArgs    string                        `json:"url_args" binding:"required"`
		Alert      *AlerterResponse              `json:"alert"`
		APNS       *PushNotificationDataResponse `json:"apns"`
	}

	// AlerterResponse represents the structure of an alert content response.
	AlerterResponse struct {
		Action          string `json:"action" binding:"required"`
		Body            string `json:"body" binding:"required"`
		Title           string `json:"title" binding:"required"`
		Subtitle        string `json:"subtitle" binding:"required"`
		SummeryArg      string `json:"summery_arg" binding:"required"`
		SummeryArgCount int    `json:"summery_arg_count" binding:"required"`
	}
)
