package psreservice

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type RequestChangeService interface {
	RequestChangePhone(token, externalID string, req *dto.RequestChangePhoneRequest) ([]byte, int, error)
	ResendPhone(token, externalID string, req *dto.ResendPhoneRequest) ([]byte, int, error)
	ValidatePhone(token, externalID string, req *dto.ValidatePhoneRequest) ([]byte, int, error)

	RequestChangeEmail(token, externalID string, req *dto.RequestChangeEmailRequest) ([]byte, int, error)
	ResendEmail(token, externalID string, req *dto.ResendEmailRequest) ([]byte, int, error)
	ValidateEmail(token, externalID string, req *dto.ValidateEmailRequest) ([]byte, int, error)
}

type requestChangeService struct{}

func NewRequestChangeService() RequestChangeService {
	return &requestChangeService{}
}

func (s *requestChangeService) RequestChangePhone(token, externalID string, req *dto.RequestChangePhoneRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/request-change/phone", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre request change phone failed: %s", string(data))
	}

	return data, status, nil
}

func (s *requestChangeService) ResendPhone(token, externalID string, req *dto.ResendPhoneRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/request-change/resend-phone", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre resend phone failed: %s", string(data))
	}

	return data, status, nil
}

func (s *requestChangeService) ValidatePhone(token, externalID string, req *dto.ValidatePhoneRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/request-change/validate-phone", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre validate phone failed: %s", string(data))
	}

	return data, status, nil
}

func (s *requestChangeService) RequestChangeEmail(token, externalID string, req *dto.RequestChangeEmailRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/request-change/email", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre request change email failed: %s", string(data))
	}

	return data, status, nil
}

func (s *requestChangeService) ResendEmail(token, externalID string, req *dto.ResendEmailRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/request-change/resend-email", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre resend email failed: %s", string(data))
	}

	return data, status, nil
}

func (s *requestChangeService) ValidateEmail(token, externalID string, req *dto.ValidateEmailRequest) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("POST", "/user/request-change/validate-email", req, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}

	if status >= 400 {
		return data, status, fmt.Errorf("psre validate email failed: %s", string(data))
	}

	return data, status, nil
}
