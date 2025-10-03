package service

import (
	"bytes"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PsreService interface {
	Login(body []byte, headers map[string]string) ([]byte, error)
}

type psreService struct {
	logRepo  repository.ClientRequestLogRepository
	userRepo repository.UserRepository
}

func NewPsreService(logRepo repository.ClientRequestLogRepository, userRepo repository.UserRepository) PsreService {
	return &psreService{logRepo: logRepo, userRepo: userRepo}
}

func (s *psreService) Login(body []byte, headers map[string]string) ([]byte, error) {
	// decode body → ambil email
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, errors.New("invalid request body")
	}

	// cek apakah email ada di table users
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email tidak terdaftar")
	}
	fmt.Println("User found:", user.Email, "Role:", user.Role)
	clientID := user.Client.ID
	fmt.Println("Login attempt for user:", req.Email, "ClientID:", clientID)

	// forward ke PSRE backend
	url := "http://10.100.20.14:2000/client/login"
	reqHttp, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		reqHttp.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(reqHttp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	// simpan log
	log := model.ClientRequestLog{
		ClientID: &clientID,
		Body:     string(body),
		Header:   "", // bisa json.Marshal(headers) kalau mau
		Response: string(resBody),
		URL:      url,
		Type:     "login",
	}
	_ = s.logRepo.Create(&log)

	return resBody, nil
}
