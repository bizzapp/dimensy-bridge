package psreintegration

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type ClientService interface{

	Login(body []byte) ([]byte, error)
}

type clientService struct {
	logRepo           repository.ClientRequestLogRepository
	userRepo          repository.UserRepository
}

func NewClientService(logRepo repository.ClientRequestLogRepository, userRepo repository.UserRepository) ClientService {
	return &clientService{logRepo: logRepo, userRepo: userRepo}
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

	fmt.Println("Login attempt for:", req.Email, "ClientID:", clientID)

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
		URL:      url+path,
		Type:     "login",
	}
	_ = s.logRepo.Create(&log)

	return respBody, nil
}
