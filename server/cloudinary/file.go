package cloudinary

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
)

func (cloudinaryClient *CloudinaryClient) UploadFile(fileUploadInput models.FileUploadPayload) ([]models.FileOutputs, error) {

	var fileUploadOutputs []models.FileOutputs

	for _, file := range fileUploadInput.Files {
		fileString := strings.Split(file.Base64, "base64,")
		decodedData, err := base64.StdEncoding.DecodeString(fileString[1])
		if err != nil {
			return []models.FileOutputs{}, err
		}

		constructedName := generateUUID()

		reader := bytes.NewReader(decodedData)

		uploadResult, err := cloudinaryClient.Upload.Upload(context.Background(), reader, uploader.UploadParams{PublicID: constructedName})
		if err != nil {
			return nil, err
		}
		fileUploadOutputs = append(fileUploadOutputs, models.FileOutputs{
			SecureURL: uploadResult.SecureURL,
			URL:       uploadResult.URL,
		})
	}
	return fileUploadOutputs, nil
}

func generateUUID() string {
	uuid := uuid.New().String()
	return uuid
}
