package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type KycService interface {
	IdValidationFull(token, externalID, id string, req *dto.ClientUserKYCRequest) ([]byte, int, error)
	LivenessFull(token, externalID, id string, req *dto.ClientUserKYCRequest) ([]byte, int, error)
}

type kycService struct {
	clientKYCHistoryRepo repository.ClientKYCHistoryRepository
}

func NewKycService(clientKYCHistoryRepo repository.ClientKYCHistoryRepository) KycService {
	return &kycService{
		clientKYCHistoryRepo: clientKYCHistoryRepo,
	}
}

func (s *kycService) IdValidationFull(token, externalID, id string, req *dto.ClientUserKYCRequest) ([]byte, int, error) {
	path := fmt.Sprintf("/kyc/id-validation/full/%s", id)
	data, status, err := utils.PsreRequest("POST", path, req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre id-validation failed: %s", string(data))
	}
	return data, status, nil
}

func (s *kycService) LivenessFull(token, externalID, id string, req *dto.ClientUserKYCRequest) ([]byte, int, error) {
	path := fmt.Sprintf("/kyc/liveness/full/%s", id)
	data, status, err := utils.PsreRequest("POST", path, req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre liveness failed: %s", string(data))
	}
	return data, status, nil
}
