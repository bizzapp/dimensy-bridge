package service

import (
	"bytes"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type PsreService interface {
	// Profile() ([]byte, error)
	CreateClientCompany(idStr string) ([]byte, error)
}

type psreService struct {
	logRepo           repository.ClientRequestLogRepository
	userRepo          repository.UserRepository
	clientCompanyRepo repository.ClientCompanyRepository
}

func NewPsreService(logRepo repository.ClientRequestLogRepository, userRepo repository.UserRepository, clientCompanyRepo repository.ClientCompanyRepository) PsreService {
	return &psreService{logRepo: logRepo, userRepo: userRepo, clientCompanyRepo: clientCompanyRepo}
}

func (s *psreService) CreateClientCompany(idStr string) ([]byte, error) {
	// parse id
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	// ambil company dari DB
	company, err := s.clientCompanyRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	// build payload sesuai PSRE API
	payload := map[string]interface{}{
		"companyName":     company.Name,
		"companyAddress":  company.Address,
		"companyIndustry": company.Industry,
		"npwp":            company.NPWP,
		"nib":             company.NIB,
		"picName":         company.PICName,
		"picEmail":        company.PICEmail,
	}
	body, _ := json.Marshal(payload)

	// request ke PSRE
	url := "http://10.100.20.14:2000/client/company/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// misal pakai token dari GetPsreToken()
	token, _ := utils.GetAdministratorToken()
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	// simpan log
	log := model.ClientRequestLog{
		ClientID: &company.ClientID,
		Body:     string(body),
		Header:   "Authorization: Bearer <hidden>",
		Response: string(resBody),
		URL:      url,
		Type:     "client_company_create",
	}
	_ = s.logRepo.Create(&log)

	return resBody, nil
}
