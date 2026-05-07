package http

import (
	"fmt"
	"net/http"

	"AccountFlow/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	uc *usecase.TransactionUseCase
}

func NewTransactionHandler(uc *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{uc: uc}
}

type createTransactionRequest struct {
	AccountID       string  `json:"account_id" binding:"required"`
	OperationTypeID int     `json:"operation_type_id" binding:"required,min=1"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
}

type transferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	IdTransfer    string    `json:"id_transfer"`
}

// Create godoc
// POST /transactions
func (h *TransactionHandler) Create(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		bindError(c, err)
		return
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	tx, err := h.uc.CreateTransaction(c.Request.Context(), accountID, req.OperationTypeID, req.Amount)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) TransferBetweenAccounts(c *gin.Context) {
	var TransferRequest transferRequest
	if err := c.ShouldBindJSON(&TransferRequest); err != nil {
		bindError(c, err)
		return
	}

	err := h.uc.TransferAmount(c.Request.Context(), TransferRequest.FromAccountID, TransferRequest.ToAccountID, TransferRequest.Amount, TransferRequest.IdTransfer)

	if err != nil {
		writeError(c, err)
		return
	}

	msg := fmt.Sprintf(
		"transfer between accounts %s of %.2f to %s completed successfully",
		TransferRequest.FromAccountID.String(),
		TransferRequest.Amount,
		TransferRequest.ToAccountID.String(),
	)

	c.JSON(http.StatusCreated, gin.H{
		"transfer": msg,
	})
}
