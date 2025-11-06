package psreservice

import (
	"dimensy-bridge/pkg/utils"
	"fmt"
)

type DocumentService interface {
	Upload(token, externalID string, file interface{}) ([]byte, int, error)
	Preview(token, externalID, documentID string) ([]byte, int, error)
}

type documentService struct {
	// Add necessary dependencies here
}

func NewDocumentService() DocumentService {
	return &documentService{}
}

func (s *documentService) Upload(token, externalID string, file interface{}) ([]byte, int, error) {
	// Implement the upload logic here, possibly calling an external API
	return nil, 0, nil
}

func (s *documentService) Preview(token, externalID, documentID string) ([]byte, int, error) {
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
