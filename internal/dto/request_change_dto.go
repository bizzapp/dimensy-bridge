package dto

type RequestChangePhoneRequest struct {
	UserID string `json:"userId" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
}

type ResendPhoneRequest struct {
	ChangeDataUserID string `json:"changeDataUserId" binding:"required"`
}

type ValidatePhoneRequest struct {
	ChangeDataUserID string `json:"changeDataUserId" binding:"required"`
	OTP              string `json:"otp" binding:"required"`
}

type RequestChangeEmailRequest struct {
	UserID string `json:"userId" binding:"required"`
	Email  string `json:"email" binding:"required"`
}

type ResendEmailRequest struct {
	ChangeDataUserID string `json:"changeDataUserId" binding:"required"`
}

type ValidateEmailRequest struct {
	ChangeDataUserID string `json:"changeDataUserId" binding:"required"`
	OTP              string `json:"otp" binding:"required"`
}
