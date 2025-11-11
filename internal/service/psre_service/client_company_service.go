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
	InviteClientCompany(authData interface{}, token string, body dto.PsreInviteClientCompanyRequest) ([]byte, int, error)
	AcceptInvitationClientUser(authData interface{}, token string, body dto.PsreAcceptInvitationClientUserRequest) ([]byte, int, error)
}

type clientCompanyService struct {
	db                      *gorm.DB
	clientSvc               service.ClientService
	clientCompanyRepo       repository.ClientCompanyRepository
	quotaClientSvc          service.QuotaClientService
	clientCompanyInviteSvc  service.ClientCompanyInviteService
	clientCompanyInviteRepo repository.ClientCompanyInviteRepository
	clientUserRepo          repository.ClientUserRepository
}

func NewClientCompanyService(
	db *gorm.DB,
	clientSvc service.ClientService,
	clientCompanyRepo repository.ClientCompanyRepository,
	quotaClientSvc service.QuotaClientService,
	clientCompanyInviteSvc service.ClientCompanyInviteService,
	clientCompanyInviteRepo repository.ClientCompanyInviteRepository,
	clientUserRepo repository.ClientUserRepository,
) ClientCompanyService {
	return &clientCompanyService{
		db:                      db,
		clientSvc:               clientSvc,
		clientCompanyRepo:       clientCompanyRepo,
		quotaClientSvc:          quotaClientSvc,
		clientCompanyInviteSvc:  clientCompanyInviteSvc,
		clientCompanyInviteRepo: clientCompanyInviteRepo,
		clientUserRepo:          clientUserRepo,
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
func (s *clientCompanyService) InviteClientCompany(authData interface{}, token string, body dto.PsreInviteClientCompanyRequest) ([]byte, int, error) {
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
			message := fmt.Sprintf("Unauthorized: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("unauthorized client: %w", err)
		}
		findClientCompanyInvite, err := s.clientCompanyInviteRepo.FindByExternal(client.ID, body.UserID, body.CompanyID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			message := fmt.Sprintf("Failed to check existing invitation: %v", err)
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("failed to check existing invitation: %w", err)
		}

		// Check if invitation exists and is already verified
		if findClientCompanyInvite != nil && findClientCompanyInvite.IsVerify {
			message := "User already member of this company"
			respBody = utils.ResponseError(message, 400)
			status = 400
			return fmt.Errorf("invitation already exists")
		}

		// Use quota only if no existing invitation
		if findClientCompanyInvite == nil {
			quantity := 1
			_, err = utils.NewQuotaUtils().UseQuota(tx, dto.UseQuotaClientRequest{
				MasterProductID: seeder.ID_PRODUCT_USER_PERSONAL_COMPANY,
				ClientID:        client.ID,
				Quantity:        int64(quantity),
			})
			if err != nil {
				message := fmt.Sprintf("Failed to use quota: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to use quota: %w", err)
			}

			clientUser, err := s.clientUserRepo.FindByExternalID(body.UserID)
			if err != nil {
				message := fmt.Sprintf("Failed to find client user: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				// fmt.Println("kesini")
				return fmt.Errorf("failed to find client user: %w", err)
			}
			clientCompany, err := s.clientCompanyRepo.FindByExternalID(body.CompanyID)
			if err != nil {
				message := fmt.Sprintf("Failed to find client company: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to find client company: %w", err)
			}
			newInvite := &model.ClientCompanyInvite{

				ClientID:          client.ID,
				ExternalUserID:    body.UserID,
				ExternalCompanyID: body.CompanyID,
				ClientUserID:      clientUser.ID,
				ClientCompanyID:   clientCompany.ID,
			}
			if err := s.clientCompanyInviteRepo.CreateTx(tx, newInvite); err != nil {
				message := fmt.Sprintf("Failed to create invitation record: %v", err)
				respBody = utils.ResponseError(message, 400)
				status = 400
				return fmt.Errorf("failed to create invitation record: %w", err)
			}
		}
		data, psreStatus, err := utils.PsreRequest("POST", "/client/company/invite", body, token, nil)
		respBody, status = data, psreStatus

		if err != nil {
			// ini trigger rollback
			return fmt.Errorf("PSrE request failed: %w", err)
		}
		if psreStatus >= 400 {
			// ini juga trigger rollback
			return fmt.Errorf("PSrE returned HTTP %d: %s", psreStatus, string(data))
		}
		return nil
	})
	if txErr != nil {
		if respBody != nil {
			return respBody, status, txErr
		}
		return nil, http.StatusBadRequest, txErr
	}

	return respBody, status, nil
}

func (s *clientCompanyService) AcceptInvitationClientUser(authData interface{}, token string, req dto.PsreAcceptInvitationClientUserRequest) ([]byte, int, error) {
	// return utils.PsreRequest("POST", "/client/users/accept-invitation", body, token, nil)
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

		data, psreStatus, err := utils.PsreRequest("POST", "/client/users/accept-invitation", req, token, nil)
		respBody, status = data, psreStatus
		if err != nil {
			return fmt.Errorf("PSrE request failed: %w", err)
		}
		if psreStatus >= 400 {
			return fmt.Errorf("PSrE returned HTTP %d: %s", psreStatus, string(data))
		}

		var resp dto.PsreAcceptInvitationResponse

		if err := json.Unmarshal(respBody, &resp); err != nil {
			return fmt.Errorf("failed to parse PSrE response: %w", err)
		}
		if resp.Code == 0 {
			clientCompanyInvite, err := s.clientCompanyInviteRepo.FindByExternal(client.ID, resp.Data.UserID, resp.Data.CompanyID)
			if err != nil {
				return fmt.Errorf("failed to find invitation record: %w", err)
			}
			if clientCompanyInvite == nil {
				return fmt.Errorf("invitation record not found")
			}
			if err := s.clientCompanyInviteRepo.VerifyInvite(clientCompanyInvite.ID); err != nil {
				return fmt.Errorf("failed to verify invitation record: %w", err)
			}
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
			MasterProductID: seeder.ID_PRODUCT_COMPANY,
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
