package service

import (
	"bytes"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ClientPsreService interface {
	RegisterPsre(clientID int64, expiredDate time.Time, email, password string) (*model.ClientPsre, error)
	GetByID(id int64) (*model.ClientPsre, error)
	GetByClientID(clientID int64) (*model.ClientPsre, error)
	UpdatePsre(psre *model.ClientPsre) error
	DeletePsre(id int64) error
}

type clientPsreService struct {
	psreRepo   repository.ClientPsreRepository
	clientRepo repository.ClientRepository
}

func NewClientPsreService(psreRepo repository.ClientPsreRepository, clientRepo repository.ClientRepository) ClientPsreService {
	return &clientPsreService{psreRepo, clientRepo}
}

func (s *clientPsreService) RegisterPsre(clientID int64, expiredDate time.Time, email, password string) (*model.ClientPsre, error) {
	// cek client ada
	client, err := s.clientRepo.FindByID(clientID)
	if err != nil {
		return nil, errors.New("client tidak ditemukan")
	}

	// jika email kosong -> ambil dari client.User.Email
	if email == "" && client.User.Email != nil {
		email = *client.User.Email
	}

	// cek psre sudah ada?
	// if _, err := s.psreRepo.FindByClientID(clientID); err == nil {
	// 	return nil, errors.New("client sudah punya psre")
	// }

	// request ke sistem eksternal
	payload := map[string]interface{}{
		"clientName":  client.CompanyName,
		"picName":     client.PicName,
		"email":       email,
		"password":    password,
		"expiredDate": expiredDate.Format("2006-01-02"),
	}

	body, _ := json.Marshal(payload)

	token, err := utils.GetPsreToken()
	if err != nil {
		return nil, err
	}

	psreUrl := os.Getenv("PSRE_BACKEND_URL")
	if psreUrl == "" {
		psreUrl = "http://10.100.20.14:2000" // fallback default
	}
	req, err := http.NewRequest("POST", psreUrl+"/backend/client/create", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	fmt.Println("Request to PSRE:", req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	clientHttp := &http.Client{}
	resp, err := clientHttp.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, errors.New("gagal register ke PSRE system")
	}

	// parse response (anggap dapat external_id)
	var result struct {
		ExternalID string `json:"externalId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// simpan ke DB
	psre := model.ClientPsre{
		ClientID:   clientID,
		ExternalID: result.ExternalID,
		ExpireDate: expiredDate,
	}
	if err := s.psreRepo.Create(&psre); err != nil {
		return nil, err
	}

	return &psre, nil
}

func (s *clientPsreService) GetByID(id int64) (*model.ClientPsre, error) {
	return s.psreRepo.FindByID(id)
}

func (s *clientPsreService) GetByClientID(clientID int64) (*model.ClientPsre, error) {
	return s.psreRepo.FindByClientID(clientID)
}

func (s *clientPsreService) UpdatePsre(psre *model.ClientPsre) error {
	return s.psreRepo.Update(psre)
}

func (s *clientPsreService) DeletePsre(id int64) error {
	return s.psreRepo.Delete(id)
}
