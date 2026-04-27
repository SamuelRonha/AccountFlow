package http

import (
	"net/http"

	"AccountFlow/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	uc *usecase.AccountUseCase
}

func NewAccountHandler(uc *usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{uc: uc}
}

type createAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}

// Create godoc
// POST /accounts
func (h *AccountHandler) Create(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		bindError(c, err)
		return
	}

	account, err := h.uc.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetByID godoc
// GET /accounts/:accountId
func (h *AccountHandler) GetByID(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	account, err := h.uc.GetByID(c.Request.Context(), accountID)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, account)
}
