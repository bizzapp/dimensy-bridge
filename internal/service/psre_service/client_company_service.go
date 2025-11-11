package psreservice

import (
	"context"
	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/model/seeder"
	"dimensy-bridge/internal/repository"
	"dimensy-bridge/internal/service"
	"dimensy-bridge/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type ClientCompanyService interface {
	CreateClientCompany(ctx context.Context, authData interface{}, token string, req dto.PsreCreateClientCompanyRequest) ([]byte, int, error)
	GetCompany(token string, params map[string]string) ([]byte, int, error)
	DetailClientCompany(token string, id string) ([]byte, int, error)
	InviteClientCompany(token string, body interface{}) ([]byte, int, error)
	AcceptInvitationClientUser(token string, body dto.PsreAcceptInvitationClientUserRequest) ([]byte, int, error)
}

type clientCompanyService struct {
	db                *gorm.DB
	clientSvc         service.ClientService
	clientCompanyRepo repository.ClientCompanyRepository
	quotaClientSvc    service.QuotaClientService
}

func NewClientCompanyService(
	db *gorm.DB,
	clientSvc service.ClientService,
	clientCompanyRepo repository.ClientCompanyRepository,
	quotaClientSvc service.QuotaClientService,
) ClientCompanyService {
	return &clientCompanyService{
		db:                db,
		clientSvc:         clientSvc,
		clientCompanyRepo: clientCompanyRepo,
		quotaClientSvc:    quotaClientSvc,
	}
}

// Get company list from PSrE
func (s *clientCompanyService) GetCompany(token string, params map[string]string) ([]byte, int, error) {
	return utils.PsreRequest("GET", "/client/company", nil, token, params)
}

// Get company detail from PSrE
func (s *clientCompanyService) DetailClientCompany(token string, id string) ([]byte, int, error) {
	return utils.PsreRequest("GET", "/client/company/detail/"+id, nil, token, nil)
}

// Invite company via PSrE
func (s *clientCompanyService) InviteClientCompany(token string, body interface{}) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/client/company/invite", body, token, nil)
}

func (s *clientCompanyService) AcceptInvitationClientUser(token string, body dto.PsreAcceptInvitationClientUserRequest) ([]byte, int, error) {
	return utils.PsreRequest("POST", "/client/users/accept-invitation", body, token, nil)
}

// Create company both in local DB and PSrE
func (s *clientCompanyService) CreateClientCompany(
	ctx context.Context,
	authData interface{},
	token string,
	req dto.PsreCreateClientCompanyRequest,
) ([]byte, int, error) {

	var (
		respBody []byte
		status   int
	)
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		externalID, err := utils.ExtractExternalID(authData)
		if err != nil {
			return fmt.Errorf("unauthorized: %w", err)
		}

		client, err := s.clientSvc.GetClientByExternalId(externalID)
		if err != nil {
			return fmt.Errorf("unauthorized client: %w", err)
		}

		quantity := 1
		_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
			MasterProductID: seeder.ID_PRODUCT_USER_PERSONAL,
			ClientID:        client.ID,
			Quantity:        int64(quantity),
		})
		if err != nil {
			dataResp := map[string]interface{}{
				"error": err.Error(),
			}
			respBody, _ = json.Marshal(dataResp)
			status = http.StatusBadRequest

			return fmt.Errorf("failed to use quota: %w", err)
		}
		// fmt.Println("Use Quota Success:", useQuota)

		clientCompany := &model.ClientCompany{
			ClientID: client.ID,
			Name:     req.CompanyName,
			Address:  req.CompanyAddress,
			Industry: req.CompanyIndustry,
			NPWP:     req.NPWP,
			NIB:      req.NIB,
			PICName:  req.PICName,
			PICEmail: req.PICEmail,
			// Status:   "PENDING",
		}

		if err := s.clientCompanyRepo.CreateTx(tx, clientCompany); err != nil {
			return fmt.Errorf("failed to create local client company: %w", err)
		}

		// 🔹 External API Call (critical section)
		data, psreStatus, err := utils.PsreRequest("POST", "/client/company/create", req, token, nil)
		respBody, status = data, psreStatus

		if err != nil {
			// ini trigger rollback
			return fmt.Errorf("PSrE request failed: %w", err)
		}
		if psreStatus >= 400 {
			// ini juga trigger rollback
			return fmt.Errorf("PSrE returned HTTP %d: %s", psreStatus, string(data))
		}

		var resp dto.PsreRegisterCompanyResponse
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return fmt.Errorf("failed to parse PSrE response: %w", err)
		}
		if resp.CompanyID == "" {
			return errors.New("invalid PSrE response: missing CompanyID")
		}

		// 🔹 update external id

		if err := s.clientCompanyRepo.UpdateExternalIDTx(tx, clientCompany.ID, resp.CompanyID); err != nil {
			return fmt.Errorf("failed to update external id: %w", err)
		}
		return nil
	})

	// rollback handled automatically by GORM if txErr != nil
	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		return nil, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}
