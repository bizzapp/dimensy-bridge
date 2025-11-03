package psrehandler

import "github.com/gin-gonic/gin"

type PsreUserHandler struct {
}

func NewPsreUserHandler() *PsreUserHandler {
	return &PsreUserHandler{}
}

func (h *PsreUserHandler) Register(c *gin.Context) {
	// Implementasi pendaftaran user PSRE di sini
}
