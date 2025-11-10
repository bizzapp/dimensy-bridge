package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type ClientDocumentService interface {
	UploadSingle(token, externalID string, dto dto.PsreDocumentSingleFileRequest) ([]byte, int, error)
	UploadBulk(token, externalID string, dto dto.PsreDocumentBulkFileRequest) ([]byte, int, error)
	Preview(token, externalID, documentID string) ([]byte, int, error)
	RequestSign(token, externalID string, dto dto.PsreDocumentSignRequest) ([]byte, int, error)
	ProcessSign(token, externalID string, dto dto.PsreDocumentProcessSignRequest) ([]byte, int, error)
	RequestStamp(token, externalID string, dto dto.PsreDocumentStampRequest) ([]byte, int, error)
	ProcessStamp(token, externalID string, dto dto.PsreDocumentProcessStampRequest) ([]byte, int, error)
	RequestOtpSign(token, externalID string, dto dto.PsreDocumentOtpSignRequest) ([]byte, int, error)
}

type clientDocumentService struct {
	db                        *gorm.DB
	clientPsreSvc             service.ClientPsreService
	clientDocumentRepo        repository.ClientDocumentRepository
	clientDocumentProcessRepo repository.ClientDocumentProcessRepository
	// Add necessary dependencies here
}

func NewClientDocumentService(db *gorm.DB, clientPsreSvc service.ClientPsreService, clientDocumentRepo repository.ClientDocumentRepository) ClientDocumentService {
	return &clientDocumentService{
		db:                 db,
		clientPsreSvc:      clientPsreSvc,
		clientDocumentRepo: clientDocumentRepo,
	}
}

func (s *clientDocumentService) UploadSingle(token, externalID string, req dto.PsreDocumentSingleFileRequest) ([]byte, int, error) {
	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Calculate file size from base64 document
		fileSize, err := utils.CalculateBase64FileSize(req.Document)
		if err != nil {
			return fmt.Errorf("failed to calculate file size: %w", err)
		}

		totalParticipants := 1
		typeDocument := "SINGLE"
		callbackURL := os.Getenv("APP_URL_WEBHOOK_NOTIFICATION")
		// fmt.Println(dto.CallbackURL)
		clientDocument := &model.ClientDocument{
			ClientID:          client.ID,
			FileName:          req.FileName,
			Type:              typeDocument,
			Status:            model.DOCUMENT_STATUS_PENDING,
			TotalParticipants: &totalParticipants,
			FileSizeKB:        &fileSize,
			CallbackURL:       &callbackURL,
			ClientCallbackURL: &req.CallbackURL,
		}
		if err := tx.Create(&clientDocument).Error; err != nil {
			return fmt.Errorf("failed create client document: %w", err)
		}

		// dto.CallbackURL = callbackURL
		dtoSingleFile := dto.PsreDocumentSingleFileRequest{
			FileName:    req.FileName,
			Document:    req.Document,
			CallbackURL: callbackURL,
		}
		data, st, err := utils.PsreRequest("POST", "/document/upload", dtoSingleFile, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		var psreResp struct {
			Code       int    `json:"code"`
			Message    string `json:"message"`
			DocumentID string `json:"documentId"`
		}
		if err := json.Unmarshal(data, &psreResp); err != nil {
			return errors.New("invalid psre response format")
		}

		if psreResp.Code != 0 {
			return fmt.Errorf("psre register failed: %s", psreResp.Message)
		}

		clientDocument.ExternalID = &psreResp.DocumentID
		if err := tx.Save(&clientDocument).Error; err != nil {
			return fmt.Errorf("failed update external_id: %w", err)
		}
		return nil
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *clientDocumentService) Preview(token, externalID, documentID string) ([]byte, int, error) {
	path := fmt.Sprintf("/document/preview/%s", documentID)
	data, status, err := utils.PsreRequest("GET", path, nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre get document detail failed: %s", string(data))
	}
	return data, status, nil
}
func (s *clientDocumentService) UploadBulk(token, externalID string, dto dto.PsreDocumentBulkFileRequest) ([]byte, int, error) {
	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		totalParticipants := len(dto.Documents)
		typeDocument := "BULK"

		// Simpan dokumen yang berhasil dibuat
		var createdDocs []model.ClientDocument
		callbackURL := os.Getenv("APP_URL_WEBHOOK_NOTIFICATION")

		for _, document := range dto.Documents {
			// Hitung ukuran file dari base64
			fileSize, err := utils.CalculateBase64FileSize(document.Document)
			if err != nil {
				return fmt.Errorf("failed to calculate file size: %w", err)
			}

			clientDocument := model.ClientDocument{
				ClientID:          client.ID,
				FileName:          document.FileName,
				Type:              typeDocument,
				Status:            model.DOCUMENT_STATUS_PENDING,
				TotalParticipants: &totalParticipants,
				FileSizeKB:        &fileSize,
				CallbackURL:       &callbackURL,
				ClientCallbackURL: &dto.CallbackURL,
			}

			if err := tx.Create(&clientDocument).Error; err != nil {
				return fmt.Errorf("failed create client document: %w", err)
			}

			createdDocs = append(createdDocs, clientDocument)
		}

		dto.CallbackURL = callbackURL
		// Call PSRE API
		data, st, err := utils.PsreRequest("POST", "/document/upload-bulk", dto, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		// Parse response PSRE
		var psreResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			GroupID string `json:"groupId"`
		}
		if err := json.Unmarshal(data, &psreResp); err != nil {
			return errors.New("invalid psre response format")
		}

		if psreResp.Code != 0 {
			return fmt.Errorf("psre register failed: %s", psreResp.Message)
		}

		// Update semua dokumen yang berhasil dibuat dengan GroupExternalID
		for _, doc := range createdDocs {
			if err := tx.Model(&model.ClientDocument{}).
				Where("id = ?", doc.ID).
				Update("group_external_id", psreResp.GroupID).Error; err != nil {
				return fmt.Errorf("failed update group_external_id: %w", err)
			}
		}

		return nil
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *clientDocumentService) RequestSign(token, externalID string, dto dto.PsreDocumentSignRequest) ([]byte, int, error) {
	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}
	var (
		respBody []byte
		status   int
	)
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Calculate expire time based on current time + DocumentProcessExpiredMinutes
		expireTime := time.Now().Add(time.Duration(model.DocumentProcessExpiredHour) * time.Hour)

		clientDocumentProcess := &model.ClientDocumentProcess{
			ClientID:          client.ID,
			ExternalID:        dto.DocumentOrGroupID,
			ExternalUserID:    dto.UserID,
			ExternalCompanyID: dto.CompanyID,
			Status:            model.ClientDocumentProcessStatusWaiting,
			ExpireTime:        &expireTime,
			Type:              model.TypeSignMeterai,
		}
		if err := tx.Create(&clientDocumentProcess).Error; err != nil {
			return fmt.Errorf("failed create client document process: %w", err)
		}
		for _, position := range dto.Positions {

			fileSize, err := utils.CalculateBase64FileSize(position.Image)
			if err != nil {
				return fmt.Errorf("failed to calculate file size: %w", err)
			}

			clientDocumenProcessDetail := &model.ClientDocumentProcessDetail{
				ClientID:                client.ID,
				ClientDocumentProcessID: clientDocumentProcess.ID,
				Type:                    position.StampType,
				Reason:                  position.Reason,
				Location:                position.Location,
				X:                       *position.X,
				Y:                       *position.Y,
				W:                       *position.W,
				H:                       *position.H,
				Page:                    position.Page,
				ImageFileSizeKB:         &fileSize,
			}
			if err := tx.Create(&clientDocumenProcessDetail).Error; err != nil {
				return fmt.Errorf("failed create client document process detail: %w", err)
			}
		}
		data, st, err := utils.PsreRequest("POST", "/document/request-sign", dto, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}
		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}
		return nil
	})

	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *clientDocumentService) ProcessSign(token, externalID string, dto dto.PsreDocumentProcessSignRequest) ([]byte, int, error) {
	_, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		data, st, err := utils.PsreRequest("POST", "/document/proccess-sign", dto, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		if err := tx.Model(&model.ClientDocumentProcess{}).
			Where("external_id = ?", dto.DocumentOrGroupID).
			Update("is_process", true).Error; err != nil {

			return fmt.Errorf("failed update group_external_id: %w", err)
		}
		return nil
	})
	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *clientDocumentService) RequestStamp(token, externalID string, dto dto.PsreDocumentStampRequest) ([]byte, int, error) {
	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}
	var (
		respBody []byte
		status   int
	)
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		expireTime := time.Now().Add(time.Duration(model.DocumentProcessExpiredHour) * time.Hour)

		clientDocumentProcess := &model.ClientDocumentProcess{
			ClientID:          client.ID,
			ExternalID:        dto.DocumentOrGroupID,
			ExternalUserID:    &dto.UserID,
			ExternalCompanyID: &dto.CompanyID,
			Status:            model.ClientDocumentProcessStatusWaiting,
			ExpireTime:        &expireTime,
			Type:              model.TypeStamp,
		}
		if err := tx.Create(&clientDocumentProcess).Error; err != nil {
			return fmt.Errorf("failed create client document process: %w", err)
		}
		for _, position := range dto.Positions {

			fileSize, err := utils.CalculateBase64FileSize(position.Image)
			if err != nil {
				return fmt.Errorf("failed to calculate file size: %w", err)
			}
			clientDocumenProcessDetail := &model.ClientDocumentProcessDetail{
				ClientID:                client.ID,
				ClientDocumentProcessID: clientDocumentProcess.ID,
				Reason:                  position.Reason,
				Location:                position.Location,
				X:                       *position.X,
				Y:                       *position.Y,
				W:                       *position.W,
				H:                       *position.H,
				Page:                    position.Page,
				ImageFileSizeKB:         &fileSize,
			}
			if err := tx.Create(&clientDocumenProcessDetail).Error; err != nil {
				return fmt.Errorf("failed create client document process detail: %w", err)
			}
		}

		data, st, err := utils.PsreRequest("POST", "/document/request-stamp", dto, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}
		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}
		return nil
	})
	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *clientDocumentService) ProcessStamp(token, externalID string, dto dto.PsreDocumentProcessStampRequest) ([]byte, int, error) {

	_, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		data, st, err := utils.PsreRequest("POST", "/document/proccess-stamp", dto, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		if err := tx.Model(&model.ClientDocumentProcess{}).
			Where("external_id = ?", dto.DocumentOrGroupID).
			Update("is_process", true).Error; err != nil {

			return fmt.Errorf("failed update external_id: %w", err)
		}
		return nil
	})
	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *clientDocumentService) RequestOtpSign(token, externalID string, dto dto.PsreDocumentOtpSignRequest) ([]byte, int, error) {

	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}
	var (
		respBody []byte
		status   int
	)
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		clientDocumentResendOtp := &model.ClientDocumentResendOtp{
			ClientID:          client.ID,
			ExternalID:        dto.DocumentOrGroupID,
			ExternalUserID:    dto.UserID,
			ExternalCompanyID: dto.CompanyID,
			Type:              dto.DocumentType,
		}
		if err := tx.Create(&clientDocumentResendOtp).Error; err != nil {
			return fmt.Errorf("failed create client document process: %w", err)
		}

		data, st, err := utils.PsreRequest("POST", "/document/request-otp-sign", dto, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}
		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}
		return nil
	})
	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		errJSON, _ := json.Marshal(map[string]any{"code": 400, "message": txErr.Error()})
		return errJSON, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}
