package psreservice

import (
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type DashboardService interface {
	GetCertificateDashboard(token, externalID string) ([]byte, int, error)
	GetDocumentDashboard(token, externalID string) ([]byte, int, error)
}

type dashboardService struct {
	// Add necessary dependencies here
}

func NewDashboardService() DashboardService {
	return &dashboardService{}
}
func (s *dashboardService) GetCertificateDashboard(token, externalID string) ([]byte, int, error) {
	path := "/dashboard/certificates"
	data, status, err := utils.PsreRequest("GET", path, nil, token, map[string]string{
		"external_id": externalID,
	})
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre get certificate dashboard failed: %s", string(data))
	}
	return data, status, nil
}

func (s *dashboardService) GetDocumentDashboard(token, externalID string) ([]byte, int, error) {
	path := "/dashboard/documents"
	data, status, err := utils.PsreRequest("GET", path, nil, token, map[string]string{
		"external_id": externalID,
	})
	if err != nil {
		return data, status, fmt.Errorf("failed call psre api: %w", err)
	}
	if status >= 400 {
		return data, status, fmt.Errorf("psre get document dashboard failed: %s", string(data))
	}
	return data, status, nil
}
