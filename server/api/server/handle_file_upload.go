package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/utils"
)

func (s *Server) HandleFileUpload(w http.ResponseWriter, r *http.Request) error {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	payload := models.FileUploadPayload{}
	err = json.Unmarshal(reqBody, &payload)
	if err != nil {
		return utils.SendErrorResponse(w, "invalid payload", "invalid_payload", http.StatusBadRequest)
	}

	fileUploadOUtput, err := s.Cloudinary.UploadFile(payload)
	if err != nil {
		return utils.SendErrorResponse(w, "error uploading file", "file_upload_error", http.StatusInternalServerError)
	}

	data, _ := json.Marshal(&fileUploadOUtput)
	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}
