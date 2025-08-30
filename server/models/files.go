package models

type FileUploadPayload struct {
	Files []struct {
		Base64 string `json:"base64String"`
		Type   string `json:"type"`
		Folder string `json:"folder"`
	} `json:"files"`
}

type FileOutputs struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
}
