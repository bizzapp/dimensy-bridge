package service

import (
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"errors"
	"time"
)

type ClientIPWhitelistService interface {
	Create(req *dto.CreateClientIPWhitelistRequest) (*dto.ClientIPWhitelistResponse, error)
	GetByID(id int64) (*dto.ClientIPWhitelistResponse, error)
	GetByClientID(clientID int64, page, limit int) (*dto.ListClientIPWhitelistResponse, error)
	Update(id int64, req *dto.UpdateClientIPWhitelistRequest) (*dto.ClientIPWhitelistResponse, error)
	Delete(id int64) error
	IsIPWhitelisted(clientID int64, ipAddress string) (bool, error)
	GetActiveIPsByClientID(clientID int64) ([]string, error)
}

type clientIPWhitelistService struct {
	repo repository.ClientIPWhitelistRepository
}

func NewClientIPWhitelistService(repo repository.ClientIPWhitelistRepository) ClientIPWhitelistService {
	return &clientIPWhitelistService{repo: repo}
}

func (s *clientIPWhitelistService) Create(req *dto.CreateClientIPWhitelistRequest) (*dto.ClientIPWhitelistResponse, error) {
	ipWhitelist := &model.ClientIPWhitelist{
		ClientID:    req.ClientID,
		IPAddress:   req.IPAddress,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	result, err := s.repo.Create(ipWhitelist)
	if err != nil {
		return nil, err
	}

	return s.toDTO(result), nil
}

func (s *clientIPWhitelistService) GetByID(id int64) (*dto.ClientIPWhitelistResponse, error) {
	ipWhitelist, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("IP whitelist tidak ditemukan")
	}

	return s.toDTO(ipWhitelist), nil
}

func (s *clientIPWhitelistService) GetByClientID(clientID int64, page, limit int) (*dto.ListClientIPWhitelistResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	ipWhitelists, total, err := s.repo.GetByClientID(clientID, page, limit)
	if err != nil {
		return nil, err
	}

	data := make([]dto.ClientIPWhitelistResponse, len(ipWhitelists))
	for i, ip := range ipWhitelists {
		data[i] = *s.toDTO(&ip)
	}

	return &dto.ListClientIPWhitelistResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *clientIPWhitelistService) Update(id int64, req *dto.UpdateClientIPWhitelistRequest) (*dto.ClientIPWhitelistResponse, error) {
	ipWhitelist := &model.ClientIPWhitelist{
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	result, err := s.repo.Update(id, ipWhitelist)
	if err != nil {
		return nil, err
	}

	return s.toDTO(result), nil
}

func (s *clientIPWhitelistService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *clientIPWhitelistService) IsIPWhitelisted(clientID int64, ipAddress string) (bool, error) {
	return s.repo.IsIPWhitelisted(clientID, ipAddress)
}

func (s *clientIPWhitelistService) GetActiveIPsByClientID(clientID int64) ([]string, error) {
	return s.repo.GetActiveIPsByClientID(clientID)
}

func (s *clientIPWhitelistService) toDTO(ipWhitelist *model.ClientIPWhitelist) *dto.ClientIPWhitelistResponse {
	var createdAt, updatedAt string
	if ipWhitelist.CreatedAt != nil {
		createdAt = ipWhitelist.CreatedAt.Format(time.RFC3339)
	}
	if ipWhitelist.UpdatedAt != nil {
		updatedAt = ipWhitelist.UpdatedAt.Format(time.RFC3339)
	}

	return &dto.ClientIPWhitelistResponse{
		ID:          ipWhitelist.ID,
		ClientID:    ipWhitelist.ClientID,
		IPAddress:   ipWhitelist.IPAddress,
		Description: ipWhitelist.Description,
		IsActive:    ipWhitelist.IsActive,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
