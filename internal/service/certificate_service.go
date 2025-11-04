package service

import (
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
	"time"
)

type CertificateService interface {
	CreateCertificate(cert *model.Certificate) error
	GetCertificateByID(id uint64) (*model.Certificate, error)
	GetCertificatesByClientID(clientID uint64) ([]model.Certificate, error)
	RevokeCertificate(id uint64) error
}

type certificateService struct {
	repo repository.CertificateRepository
}

func NewCertificateService(repo repository.CertificateRepository) CertificateService {
	return &certificateService{repo: repo}
}

func (s *certificateService) CreateCertificate(cert *model.Certificate) error {
	return s.repo.Create(cert)
}

func (s *certificateService) GetCertificateByID(id uint64) (*model.Certificate, error) {
	return s.repo.FindByID(id)
}

func (s *certificateService) GetCertificatesByClientID(clientID uint64) ([]model.Certificate, error) {
	return s.repo.FindAllByClientID(clientID)
}

func (s *certificateService) RevokeCertificate(id uint64) error {
	cert, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	cert.Status = "revoked"
	cert.RevokeAt = &now

	return s.repo.Update(cert)
}
