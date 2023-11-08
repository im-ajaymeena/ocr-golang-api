package common

type ImageRequest struct {
	SingleImageData string   `json:"image_data,omitempty"`      // For single image
	ImageDataList   []string `json:"image_data_list,omitempty"` // For list of images
}

type SyncTextResponse struct {
	Text string `json:"text"`
}

type AsyncTextResponse struct {
	Text string `json:"task_id"`
}

type TaskIDMessage struct {
	TaskID string `json:"task_id"`
}

type TaskStatus string

// Define constants for valid status of tasks values
const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusFailed    TaskStatus = "FAILED"
)

type OCRTaskInput struct {
	TaskID    string
	ImageData string
}

type OCRTaskOutput struct {
	Status         TaskStatus
	RecognizedText string
}

type Base64Response struct {
	Base64 string `json:"base64"`
}
