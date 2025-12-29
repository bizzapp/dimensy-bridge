package dto

type CreateClientIPWhitelistRequest struct {
	ClientID    int64  `json:"client_id" binding:"required"`
	IPAddress   string `json:"ip_address" binding:"required,ip"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" binding:"required"`
}

type UpdateClientIPWhitelistRequest struct {
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type ClientIPWhitelistResponse struct {
	ID          int64  `json:"id"`
	ClientID    int64  `json:"client_id"`
	IPAddress   string `json:"ip_address"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ListClientIPWhitelistResponse struct {
	Data  []ClientIPWhitelistResponse `json:"data"`
	Total int64                       `json:"total"`
	Page  int                         `json:"page"`
	Limit int                         `json:"limit"`
}
