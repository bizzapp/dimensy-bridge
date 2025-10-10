package psreservice

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
	// decode body → ambil email & password
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, errors.New("invalid request body")
	}

	// cek user di database
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email tidak terdaftar")
	}
	clientID := user.Client.ID

	// Kirim ke PSRE pakai utilitas
	url := os.Getenv("PSRE_BACKEND_URL")
	path := "/client/login"
	respBody, err := utils.DoPsreRequest("POST", path, req, nil)
	if err != nil {
		return nil, err
	}

	// Simpan log (opsional)
	log := model.ClientRequestLog{
		ClientID: &clientID,
		Body:     string(body),
		Response: string(respBody),
		URL:      url + path,
		Type:     "login",
	}
	_ = s.logRepo.Create(&log)

	return respBody, nil
}
func (s *clientService) Register(clientID int64) (*model.ClientPsre, error) {
	// cek client ada
	client, err := s.clientRepo.FindByID(clientID)
	if err != nil {
		return nil, errors.New("client tidak ditemukan")
	}
	email := client.User.Email

	// request ke sistem eksternal
	payload := map[string]interface{}{
		"clientName":  client.CompanyName,
		"picName":     client.PicName,
		"email":       email,
		"password":    utils.DefaultPassword(),
		"expiredDate": utils.ExpireDate(),
	}

	body, _ := json.Marshal(payload)

	token, err := utils.GetAdministratorToken()
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
	// cek client ada
	client, err := s.clientRepo.FindByID(clientID)
	if err != nil {
		return nil, errors.New("client tidak ditemukan")
	}

	req := map[string]interface{}{
		"email":    client.User.Email,
		"password": utils.DefaultPassword(),
	}
	path := "/client/login"
	respBody, err := utils.DoPsreRequest("POST", path, req, nil)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(respBody, &m); err != nil {
		panic(err)
	}

	// ambil accessToken
	var token string
	if data, ok := m["data"].(map[string]interface{}); ok {
		if t, ok := data["accessToken"].(string); ok {
			token = t
		}
	}

	// buat headers pakai token
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	// lakukan request berikutnya dengan header yang benar
	pathLogin := "/client/profile"
	respBody, err = utils.DoPsreRequest("GET", pathLogin, nil, headers)
	if err != nil {
		return nil, err
	}

	// parsing JSON profile
	var profile map[string]interface{}
	if err := json.Unmarshal(respBody, &profile); err != nil {
		return nil, err
	}

	// ambil id dari profile dan masukkan ke struct
	if id, ok := profile["id"].(string); ok {
		client.ClientPsre.ExternalID = id
	} else {
		return nil, fmt.Errorf("field id tidak ditemukan dalam response profile")
	}

	return client.ClientPsre, nil
}
