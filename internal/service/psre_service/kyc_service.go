package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type KycService interface {
	IdValidationFull(token, externalID, id string) ([]byte, int, error)
	LivenessFull(token, externalID, id string) ([]byte, int, error)
	FailedKyc(token, externalID, id string) ([]byte, int, error)
	FailedKycByBody(token, externalID string, req *dto.UserCompanyPicDTO) ([]byte, int, error)
}

type kycService struct {
	clientKYCHistoryRepo repository.ClientKYCHistoryRepository
}

func NewKycService(clientKYCHistoryRepo repository.ClientKYCHistoryRepository) KycService {
	return &kycService{
		clientKYCHistoryRepo: clientKYCHistoryRepo,
	}
}

func (s *kycService) IdValidationFull(token, externalID, id string) ([]byte, int, error) {
	path := fmt.Sprintf("/kyc/id-validation/full/%s", id)
	data, status, err := utils.PsreRequest("GET", path, nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre id-validation failed: %s", string(data))
	}
	return data, status, nil
}

func (s *kycService) LivenessFull(token, externalID, id string) ([]byte, int, error) {
	path := fmt.Sprintf("/kyc/liveness/full/%s", id)
	data, status, err := utils.PsreRequest("GET", path, nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre liveness failed: %s", string(data))
	}
	return data, status, nil
}

func (s *kycService) FailedKyc(token, externalID, id string) ([]byte, int, error) {
	path := fmt.Sprintf("/kyc/failed-kyc/%s", id)
	data, status, err := utils.PsreRequest("POST", path, nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre failed-kyc failed: %s", string(data))
	}

	if s.clientKYCHistoryRepo != nil {
		_ = s.clientKYCHistoryRepo.UpdateIsRejectStatus(id, true)
	}

	return data, status, nil
}

func (s *kycService) FailedKycByBody(token, externalID string, req *dto.UserCompanyPicDTO) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/kyc/failed-kyc", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre failed-kyc failed: %s", string(data))
	}

	// KYC history disimpan dengan userId, atau userPicId untuk company (lihat RequestKYC)
	userID := req.UserID
	if userID == nil {
		userID = req.UserPicID
	}
	if s.clientKYCHistoryRepo != nil && userID != nil {
		_ = s.clientKYCHistoryRepo.UpdateIsRejectStatusByExternalUserID(*userID, true)
	}

	return data, status, nil
}
