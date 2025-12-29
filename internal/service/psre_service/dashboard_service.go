package psreservice

import (
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type DashboardService interface {
	GetCertificateDashboard(token string) ([]byte, int, error)
	GetDocumentDashboard(token string) ([]byte, int, error)
}

type dashboardService struct {
	// Add necessary dependencies here
}

func NewDashboardService() DashboardService {
	return &dashboardService{}
}
func (s *dashboardService) GetCertificateDashboard(token string) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("GET", "/backend/dashboard/certificate", nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre login failed: %s", string(data))
	}
	return data, status, nil
}

func (s *dashboardService) GetDocumentDashboard(token string) ([]byte, int, error) {
	data, status, err := utils.PsreRequest("GET", "/backend/dashboard/document", nil, token, nil)
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre login failed: %s", string(data))
	}
	return data, status, nil
}
