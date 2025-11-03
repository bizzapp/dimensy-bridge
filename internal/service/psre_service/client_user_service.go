package psreservice

type ClientUserService interface {
	Register(body []byte) ([]byte, error)
}

type clientUserService struct {
	// tambahkan dependency repository jika diperlukan
}

func NewClientUserService() ClientUserService {
	return &clientUserService{
		// inisialisasi dependency repository jika diperlukan
	}
}

func (s *clientUserService) Register(body []byte) ([]byte, error) {
	// implementasi logika pendaftaran user PSRE di sini
	return nil, nil
}
