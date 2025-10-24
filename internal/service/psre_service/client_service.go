package psreservice

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type ClientService interface {
	Login(body []byte) ([]byte, error)
	Register(clientID int64) (*model.ClientPsre, error)
	FillExternalId(clientID int64) (*model.ClientPsre, error)
	Profile(clientID int64) (*model.ClientPsre, error)
}

type clientService struct {
	logRepo        repository.ClientRequestLogRepository
	userRepo       repository.UserRepository
	clientRepo     repository.ClientRepository
	clientPsreRepo repository.ClientPsreRepository
}

func NewClientService(logRepo repository.ClientRequestLogRepository, userRepo repository.UserRepository, clientRepo repository.ClientRepository, clientPsreRepo repository.ClientPsreRepository) ClientService {
	return &clientService{
		logRepo:        logRepo,
		userRepo:       userRepo,
		clientRepo:     clientRepo,
		clientPsreRepo: clientPsreRepo}
}
func (s *clientService) Login(body []byte) ([]byte, error) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, errors.New("invalid request body")
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email tidak terdaftar")
	}
	clientID := user.Client.ID

	respBody, status, err := utils.PsreRequest("POST", "/client/login", req, "", nil)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v (status: %d)", err, status)
	}

	// Simpan log
	s.logRepo.Create(&model.ClientRequestLog{
		ClientID: &clientID,
		Body:     string(body),
		Response: string(respBody),
		URL:      os.Getenv("PSRE_BACKEND_URL") + "/client/login",
		Type:     "login",
	})

	return respBody, nil
}
func (s *clientService) Register(clientID int64) (*model.ClientPsre, error) {
	client, err := s.clientRepo.FindByID(clientID)
	if err != nil {
		return nil, errors.New("client tidak ditemukan")
	}

	token, err := utils.GetAdministratorToken()
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"clientName":  client.CompanyName,
		"picName":     client.PicName,
		"email":       client.User.Email,
		"password":    utils.DefaultPassword(),
		"expiredDate": utils.ExpireDate(),
	}

	respBody, status, err := utils.PsreRequest("POST", "/backend/client/create", payload, "Bearer "+token, nil)
	if err != nil {
		return nil, fmt.Errorf("register PSRE failed: %v (status %d)", err, status)
	}

	var result struct {
		ExternalID string `json:"externalId"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("invalid PSRE response: %w", err)
	}

	psre := model.ClientPsre{
		ClientID:   clientID,
		ExternalID: result.ExternalID,
		ExpireDate: utils.ExpireDate(),
	}
	if err := s.clientPsreRepo.Create(&psre); err != nil {
		return nil, err
	}

	return &psre, nil
}

func (s *clientService) FillExternalId(clientID int64) (*model.ClientPsre, error) {

	client, err := s.ProfilePsre(clientID)
	if err != nil {
		return nil, err
	}

	if err := s.clientPsreRepo.Update(client); err != nil {
		return nil, err
	}

	return client, nil
}

func (s *clientService) Profile(clientID int64) (*model.ClientPsre, error) {

	client, err := s.clientRepo.FindByID(clientID)
	if err != nil {
		return nil, errors.New("client tidak ditemukan")
	}

	return client.ClientPsre, nil
}
func (s *clientService) ProfilePsre(clientID int64) (*model.ClientPsre, error) {
	client, err := s.clientRepo.FindByID(clientID)
	if err != nil {
		return nil, errors.New("client tidak ditemukan")
	}

	// 2️⃣ Login ke PSRE untuk dapatkan token
	loginPayload := map[string]interface{}{
		"email":    client.User.Email,
		"password": utils.DefaultPassword(),
	}

	loginPath := "/client/login"
	respBody, _, err := utils.PsreRequest("POST", loginPath, loginPayload, "", nil)
	if err != nil {
		return nil, fmt.Errorf("gagal login ke PSRE: %w", err)
	}

	var loginResp map[string]interface{}
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("gagal decode response login: %v", err)
	}

	// 3️⃣ Ambil accessToken dari response
	var token string
	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		if t, ok := data["accessToken"].(string); ok && t != "" {
			token = t
		}
	}
	if token == "" {
		return nil, errors.New("accessToken tidak ditemukan pada response login PSRE")
	}

	// 4️⃣ Ambil profile dari PSRE
	profilePath := "/client/profile"
	respBody, _, err = utils.PsreRequest("GET", profilePath, nil, "Bearer "+token, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil profile dari PSRE: %w", err)
	}

	var profile map[string]interface{}
	if err := json.Unmarshal(respBody, &profile); err != nil {
		return nil, fmt.Errorf("gagal decode response profile: %v", err)
	}

	// 5️⃣ Ambil id / external_id dari profile
	psreID, ok := profile["id"].(string)
	if !ok || psreID == "" {
		return nil, fmt.Errorf("field id tidak ditemukan pada response profile: %s", string(respBody))
	}

	// 6️⃣ Update ke struct ClientPsre
	if client.ClientPsre == nil {
		client.ClientPsre = &model.ClientPsre{}
	}
	client.ClientPsre.ExternalID = psreID

	return client.ClientPsre, nil
}
