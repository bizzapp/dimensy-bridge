package psreservice

import "dimensy-bridge/pkg/utils"

type ClientUserService interface {
	Register(token string, body interface{}) ([]byte, int, error)
}

type clientUserService struct {
	// tambahkan dependency repository jika diperlukan
}

func NewClientUserService() ClientUserService {
	return &clientUserService{}
}

func (s *clientUserService) Register(token string, body interface{}) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/user/register", body, token, nil)
}
