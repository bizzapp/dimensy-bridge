package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/model/seeder"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
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
	RetryProcess(token, externalID string, dto dto.PsreDocumentRetryProcess) ([]byte, int, error)
	VerifyDocument(file io.Reader, filename string, queryParams map[string]string) ([]byte, int, error)
}

type clientDocumentService struct {
	db                        *gorm.DB
	clientPsreSvc             service.ClientPsreService
	clientDocumentRepo        repository.ClientDocumentRepository
	clientDocumentProcessRepo repository.ClientDocumentProcessRepository
	// Add necessary dependencies here
}

func NewClientDocumentService(db *gorm.DB, clientPsreSvc service.ClientPsreService, clientDocumentRepo repository.ClientDocumentRepository, clientDocumentProcessRepo repository.ClientDocumentProcessRepository) ClientDocumentService {
	return &clientDocumentService{
		db:                        db,
		clientPsreSvc:             clientPsreSvc,
		clientDocumentRepo:        clientDocumentRepo,
		clientDocumentProcessRepo: clientDocumentProcessRepo,
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
			status = 400
			return fmt.Errorf("failed to calculate file size: %w", err)
		}

		quantity := 1
		_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
			MasterProductID: seeder.ID_PRODUCT_DOCUMENT_UPLOAD_COUNT_LIMIT,
			ClientID:        client.ID,
			Quantity:        int64(quantity),
		})
		if err != nil {
			message := fmt.Sprintf("Failed to use quota: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("failed to use quota: %w", err)
		}

		_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
			MasterProductID: seeder.ID_PRODUCT_DOCUMENT_STORAGE_KB,
			ClientID:        client.ID,
			Quantity:        int64(fileSize),
		})
		if err != nil {
			message := fmt.Sprintf("Failed to use quota: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("failed to use quota: %w", err)
		}

		quotaLimit, err := utils.NewQuotaUtils().QuotaLimit(tx, client.ID, seeder.ID_PRODUCT_UPLOAD_SINGLE_DOCUMENT)
		if err != nil {
			status = 400
			message := fmt.Sprintf("Failed to get quota limit: %v", err)
			respBody = utils.ResponseError(message, 400)
			return fmt.Errorf("failed to get quota limit: %w", err)
		}

		// Check max single upload limit
		if quotaLimit.MaxSingleUpload != nil {
			maxUploadMB := *quotaLimit.MaxSingleUpload      // MB dari database
			maxUploadKB := utils.MBToKB(int64(maxUploadMB)) // Convert MB ke KB
			fileSizeMB := utils.KBToMB(int64(fileSize))     // Convert file size ke MB
			if int64(fileSize) > maxUploadKB {
				message := fmt.Sprintf("Max single upload: %d MB (%d KB), Current file: %d KB (%d MB)",
					maxUploadMB, maxUploadKB, fileSize, fileSizeMB)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("file size exceeds the maximum single upload limit")
			}
		} else {
			message := "No max single upload limit set - allowing upload"
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("no max single upload limit set")
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
			UserID:      req.UserID,
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
	path := fmt.Sprintf("/document/preview/%s?option=LATEST", documentID)
	data, status, err := utils.PsreRequest("GET", path, nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre get document detail failed: %s", string(data))
	}
	return data, status, nil
}
func (s *clientDocumentService) UploadBulk(token, externalID string, req dto.PsreDocumentBulkFileRequest) ([]byte, int, error) {
	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		totalParticipants := len(req.Documents)
		typeDocument := "BULK"

		// Simpan dokumen yang berhasil dibuat
		var createdDocs []model.ClientDocument
		callbackURL := os.Getenv("APP_URL_WEBHOOK_NOTIFICATION")

		quotaLimit, err := utils.NewQuotaUtils().QuotaLimit(tx, client.ID, seeder.ID_PRODUCT_UPLOAD_BULK_DOCUMENT)
		if err != nil {
			status = 400
			message := fmt.Sprintf("Failed to get quota limit: %v", err)
			respBody = utils.ResponseError(message, 400)
			return fmt.Errorf("failed to get quota limit: %w", err)
		}

		if quotaLimit.MaxBulkUploadCount == nil || *quotaLimit.MaxBulkUploadCount <= 0 {
			message := "No bulk upload quota available"
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("no bulk upload quota available")
		}
		if int64(len(req.Documents)) > int64(*quotaLimit.MaxBulkUploadCount) {
			message := fmt.Sprintf("Exceeds max bulk upload count: %d, Current count: %d", *quotaLimit.MaxBulkUploadCount, len(req.Documents))
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("exceeds max bulk upload count: %d", *quotaLimit.MaxBulkUploadCount)
		}

		totalFileSizeKB := int64(0)
		for _, document := range req.Documents {
			// Hitung ukuran file dari base64
			fileSize, err := utils.CalculateBase64FileSize(document.Document)
			if err != nil {
				return fmt.Errorf("failed to calculate file size: %w", err)
			}

			quantity := 1
			_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
				MasterProductID: seeder.ID_PRODUCT_DOCUMENT_UPLOAD_COUNT_LIMIT,
				ClientID:        client.ID,
				Quantity:        int64(quantity),
			})
			if err != nil {
				message := fmt.Sprintf("Failed to use quota: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to use quota: %w", err)
			}

			_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
				MasterProductID: seeder.ID_PRODUCT_DOCUMENT_STORAGE_KB,
				ClientID:        client.ID,
				Quantity:        int64(fileSize),
			})
			if err != nil {
				message := fmt.Sprintf("Failed to use quota: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to use quota: %w", err)
			}

			if quotaLimit.MaxBulkUploadLimitPcs != nil {
				maxUploadMB := *quotaLimit.MaxBulkUploadLimitPcs // MB dari database
				maxUploadKB := utils.MBToKB(int64(maxUploadMB))  // Convert MB ke KB
				fileSizeMB := utils.KBToMB(int64(fileSize))      // Convert file size ke MB
				if int64(fileSize) > maxUploadKB {
					message := fmt.Sprintf("Max bulk upload per file: %d MB (%d KB), Current file: %d KB (%d MB)",
						maxUploadMB, maxUploadKB, fileSize, fileSizeMB)
					respBody = utils.ResponseError(message, 400)
					status = 400
					return fmt.Errorf("file size exceeds the maximum bulk upload per file limit")
				}
			} else {
				message := "No max bulk upload per file limit set - allowing upload"
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("no max bulk upload per file limit set")
			}
			totalFileSizeKB += int64(fileSize)
		}
		if quotaLimit.MaxBulkUploadLimitAll != nil {
			maxUploadMB := *quotaLimit.MaxBulkUploadLimitAll        // MB dari database
			maxUploadKB := utils.MBToKB(int64(maxUploadMB))         // Convert MB ke KB
			totalFileSizeMB := utils.KBToMB(int64(totalFileSizeKB)) // Convert file size ke MB
			if totalFileSizeKB > maxUploadKB {
				message := fmt.Sprintf("Max bulk upload total size: %d MB (%d KB), Current total size: %d KB (%d MB)",
					maxUploadMB, maxUploadKB, totalFileSizeKB, totalFileSizeMB)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("total file size exceeds the maximum bulk upload total size limit")
			}
		} else {
			message := "No max bulk upload total size limit set - allowing upload"
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("no max bulk upload total size limit set")
		}

		for _, document := range req.Documents {
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
				ClientCallbackURL: &req.CallbackURL,
			}

			if err := tx.Create(&clientDocument).Error; err != nil {
				return fmt.Errorf("failed create client document: %w", err)
			}

			createdDocs = append(createdDocs, clientDocument)
		}

		req.CallbackURL = callbackURL
		req.UserID = req.UserID
		// Call PSRE API
		data, st, err := utils.PsreRequest("POST", "/document/upload-bulk", req, token, nil)
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

func (s *clientDocumentService) RequestSign(token, externalID string, req dto.PsreDocumentSignRequest) ([]byte, int, error) {
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

		clientDocument, err := s.clientDocumentRepo.FindByExternalID(req.DocumentOrGroupID)
		if err != nil {
			status = 400
			message := fmt.Sprintf("Failed to find client document by external id: %v", err)
			respBody = utils.ResponseError(message, status)
			return fmt.Errorf("failed to find client document by external id: %w", err)
		}

		var userID *uuid.UUID
		if req.UserID != nil && *req.UserID != "" {
			if parsed, err := uuid.Parse(*req.UserID); err == nil {
				userID = &parsed
			}
		}
		var companyID *uuid.UUID
		if req.CompanyID != nil && *req.CompanyID != "" {
			if parsed, err := uuid.Parse(*req.CompanyID); err == nil {
				companyID = &parsed
			}
		}

		clientDocuentProcess, err := s.clientDocumentProcessRepo.FindByExternalIDAndExternalUserID(req.DocumentOrGroupID, userID)
		if err == nil && clientDocuentProcess != nil {
			status = 400
			message := "This user already requested to this document. create new document or using other user"
			respBody = utils.ResponseError(message, status)
			return fmt.Errorf("client document process already exists")
		}

		quantity := 1
		_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
			MasterProductID: seeder.ID_PRODUCT_OTP,
			ClientID:        client.ID,
			Quantity:        int64(quantity),
		})
		if err != nil {
			message := fmt.Sprintf("Failed to use quota: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("failed to use quota: %w", err)
		}

		clientDocumentProcess := &model.ClientDocumentProcess{
			ClientID:          client.ID,
			ClientDocumentID:  clientDocument.ID,
			ExternalID:        req.DocumentOrGroupID,
			ExternalGroupID:   &req.DocumentOrGroupID,
			ExternalUserID:    userID,
			ExternalCompanyID: companyID,
			Status:            model.ClientDocumentProcessStatusWaiting,
			ExpireTime:        &expireTime,
			Type:              model.TypeSignMeterai,
		}
		if err := tx.Create(&clientDocumentProcess).Error; err != nil {
			status = 400
			message := fmt.Sprintf("Failed to create client document process: %v", err)
			respBody = utils.ResponseError(message, status)
			return fmt.Errorf("failed create client document process: %w", err)
		}

		for _, position := range req.Positions {

			fileSize, err := utils.CalculateBase64FileSize(position.Image)
			if err != nil {
				status = 400
				message := fmt.Sprintf("Failed to calculate file size: %v", err)
				respBody = utils.ResponseError(message, status)
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
				status = 400
				message := fmt.Sprintf("Failed to create client document process detail: %v", err)
				respBody = utils.ResponseError(message, status)
				return fmt.Errorf("failed create client document process detail: %w", err)
			}
		}
		data, st, err := utils.PsreRequest("POST", "/document/request-sign", req, token, nil)
		respBody, status = data, st
		if err != nil {
			status = 400
			message := fmt.Sprintf("Failed to call psre api: %v", err)
			respBody = utils.ResponseError(message, status)
			return fmt.Errorf("failed call psre api: %w", err)
		}
		if status >= 400 {
			message := fmt.Sprintf("Psre api error: %s", string(data))
			respBody = utils.ResponseError(message, status)
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

func (s *clientDocumentService) ProcessSign(token, externalID string, req dto.PsreDocumentProcessSignRequest) ([]byte, int, error) {
	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Get client document process details to calculate quota usage
		var clientDocumentProcessDetails []model.ClientDocumentProcessDetail
		if err := tx.Joins("JOIN client_document_processes cdp ON cdp.id = client_document_process_details.client_document_process_id").
			Joins("JOIN client_documents cd ON cd.id = cdp.client_document_id").
			Where("cd.external_id = ?", req.DocumentOrGroupID).
			Or("cd.group_external_id = ?", req.DocumentOrGroupID).
			Find(&clientDocumentProcessDetails).Error; err != nil {
			return fmt.Errorf("failed to get client document process details: %w", err)
		}

		// Use quota based on client_document_process_details
		for _, detail := range clientDocumentProcessDetails {
			masterProductID := seeder.ID_PRODUCT_SIGN
			if detail.Type == model.StampTypeMeterai {
				masterProductID = seeder.ID_PRODUCT_METERAI
			}
			quantity := 1
			_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
				MasterProductID: masterProductID,
				ClientID:        client.ID,
				Quantity:        int64(quantity),
			})
			if err != nil {
				message := fmt.Sprintf("Failed to use quota: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to use quota: %w", err)
			}
		}

		data, st, err := utils.PsreRequest("POST", "/document/proccess-sign", req, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		if err := tx.Model(&model.ClientDocumentProcess{}).
			Where("external_id = ?", req.DocumentOrGroupID).
			Update("is_process", true).Error; err != nil {

			return fmt.Errorf("failed update is_process: %w", err)
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

func (s *clientDocumentService) RequestStamp(token, externalID string, req dto.PsreDocumentStampRequest) ([]byte, int, error) {
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

		clientDocument, err := s.clientDocumentRepo.FindByExternalID(req.DocumentOrGroupID)
		if err != nil {
			status = 400
			message := fmt.Sprintf("Failed to find client document by external id: %v", err)
			respBody = utils.ResponseError(message, status)
			return fmt.Errorf("failed to find client document by external id: %w", err)
		}

		clientDocuentProcess, err := s.clientDocumentProcessRepo.FindByExternalIDAndExternalUserID(req.DocumentOrGroupID, &req.UserID)
		if err == nil && clientDocuentProcess != nil {
			status = 400
			message := "This user already requested to this document. create new document or using other user"
			respBody = utils.ResponseError(message, status)
			return fmt.Errorf("client document process already exists")
		}

		quantity := 1
		_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
			MasterProductID: seeder.ID_PRODUCT_OTP,
			ClientID:        client.ID,
			Quantity:        int64(quantity),
		})
		if err != nil {
			message := fmt.Sprintf("Failed to use quota: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("failed to use quota: %w", err)
		}

		clientDocumentProcess := &model.ClientDocumentProcess{
			ClientID:          client.ID,
			ClientDocumentID:  clientDocument.ID,
			ExternalID:        req.DocumentOrGroupID,
			ExternalUserID:    &req.UserID,
			ExternalCompanyID: &req.CompanyID,
			Status:            model.ClientDocumentProcessStatusWaiting,
			ExpireTime:        &expireTime,
			Type:              model.TypeStamp,
		}
		if err := tx.Create(&clientDocumentProcess).Error; err != nil {
			return fmt.Errorf("failed create client document process: %w", err)
		}

		for _, position := range req.Positions {

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

		data, st, err := utils.PsreRequest("POST", "/document/request-stamp", req, token, nil)
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

func (s *clientDocumentService) ProcessStamp(token, externalID string, req dto.PsreDocumentProcessStampRequest) ([]byte, int, error) {

	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Get client document process details to calculate quota usage
		var clientDocumentProcessDetails []model.ClientDocumentProcessDetail
		if err := tx.Joins("JOIN client_document_processes cdp ON cdp.id = client_document_process_details.client_document_process_id").
			Joins("JOIN client_documents cd ON cd.id = cdp.client_document_id").
			Where("cd.external_id = ?", req.DocumentOrGroupID).
			Or("cd.group_external_id = ?", req.DocumentOrGroupID).
			Find(&clientDocumentProcessDetails).Error; err != nil {
			return fmt.Errorf("failed to get client document process details: %w", err)
		}

		// Use quota based on client_document_process_details for stamp
		for range clientDocumentProcessDetails {
			masterProductID := seeder.ID_PRODUCT_STAMP
			quantity := 1
			_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
				MasterProductID: masterProductID,
				ClientID:        client.ID,
				Quantity:        int64(quantity),
			})
			if err != nil {
				message := fmt.Sprintf("Failed to use quota: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to use quota: %w", err)
			}
		}

		data, st, err := utils.PsreRequest("POST", "/document/proccess-stamp", req, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		if err := tx.Model(&model.ClientDocumentProcess{}).
			Where("external_id = ?", req.DocumentOrGroupID).
			Update("is_process", true).Error; err != nil {

			return fmt.Errorf("failed update is_process: %w", err)
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

func (s *clientDocumentService) RequestOtpSign(token, externalID string, req dto.PsreDocumentOtpSignRequest) ([]byte, int, error) {

	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}
	var (
		respBody []byte
		status   int
	)
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Convert string IDs to UUID pointers

		var userID *uuid.UUID
		if req.UserID != nil {
			if parsed, err := uuid.Parse(*req.UserID); err == nil {
				userID = &parsed
			}
		}
		var companyID *uuid.UUID
		if req.CompanyID != nil {
			if parsed, err := uuid.Parse(*req.CompanyID); err == nil {
				companyID = &parsed
			}
		}

		clientDocumentResendOtp := &model.ClientDocumentResendOtp{
			ClientID:          client.ID,
			ExternalID:        req.DocumentOrGroupID,
			ExternalUserID:    userID,
			ExternalCompanyID: companyID,
			// Type:              dto.DocumentType,
		}
		if err := tx.Create(&clientDocumentResendOtp).Error; err != nil {
			return fmt.Errorf("failed create client document process: %w", err)
		}

		data, st, err := utils.PsreRequest("POST", "/document/request-otp-sign", req, token, nil)
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

func (s *clientDocumentService) RetryProcess(token, externalID string, req dto.PsreDocumentRetryProcess) ([]byte, int, error) {

	client, err := s.clientPsreSvc.GetByExternalID(externalID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed get client psre: %w", err)
	}

	var (
		respBody []byte
		status   int
	)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Get client document process details to calculate quota usage
		var clientDocumentProcessDetails []model.ClientDocumentProcessDetail
		if err := tx.Joins("JOIN client_document_processes cdp ON cdp.id = client_document_process_details.client_document_process_id").
			Joins("JOIN client_documents cd ON cd.id = cdp.client_document_id").
			Where("cd.external_id = ?", req.DocumentOrGroupID).
			Or("cd.group_external_id = ?", req.DocumentOrGroupID).
			Find(&clientDocumentProcessDetails).Error; err != nil {
			return fmt.Errorf("failed to get client document process details: %w", err)
		}

		// Use quota based on client_document_process_details for stamp
		for _, detail := range clientDocumentProcessDetails {
			masterProductID := seeder.ID_PRODUCT_STAMP
			if detail.Type == model.StampTypeMeterai {
				masterProductID = seeder.ID_PRODUCT_METERAI
			}
			if detail.Type == model.StampTypeSign {
				masterProductID = seeder.ID_PRODUCT_SIGN
			}
			quantity := 1
			_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
				MasterProductID: masterProductID,
				ClientID:        client.ID,
				Quantity:        int64(quantity),
			})
			if err != nil {
				message := fmt.Sprintf("Failed to use quota: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to use quota: %w", err)
			}
		}

		data, st, err := utils.PsreRequest("POST", "/document/retry-process", req, token, nil)
		respBody, status = data, st
		if err != nil {
			return fmt.Errorf("failed call psre api: %w", err)
		}

		if status >= 400 {
			return fmt.Errorf("psre api error: %s", string(data))
		}

		if err := tx.Model(&model.ClientDocumentProcess{}).
			Where("external_id = ?", req.DocumentOrGroupID).
			Update("is_process", true).Error; err != nil {

			return fmt.Errorf("failed update is_process: %w", err)
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

func (s *clientDocumentService) VerifyDocument(file io.Reader, filename string, queryParams map[string]string) ([]byte, int, error) {
	data, status, err := utils.PsreMultipartRequest("POST", "/document/verify", file, filename, "", queryParams)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	return data, status, nil
}
