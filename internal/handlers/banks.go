package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maxqr-api/internal/vietqr"
)

// BankHandler handles bank-related requests
type BankHandler struct{}

// NewBankHandler creates a new bank handler
func NewBankHandler() *BankHandler {
	return &BankHandler{}
}

// ListBanks handles GET /api/v1/banks
func (h *BankHandler) ListBanks(c *gin.Context) {
	banks := vietqr.GetAllBanks()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(banks),
		"data":    banks,
	})
}

// GetBank handles GET /api/v1/banks/:identifier
func (h *BankHandler) GetBank(c *gin.Context) {
	identifier := c.Param("identifier")

	// Try to find by BIN first
	bank := vietqr.GetBankByBIN(identifier)

	// Try by code
	if bank == nil {
		bank = vietqr.GetBankByCode(strings.ToUpper(identifier))
	}

	// Try by short name
	if bank == nil {
		bank = vietqr.GetBankByShortName(identifier)
	}

	if bank == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "bank_not_found",
			"message": "Bank not found. Use BIN code, bank code, or short name.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bank,
	})
}

// SearchBanks handles GET /api/v1/banks/search
func (h *BankHandler) SearchBanks(c *gin.Context) {
	query := strings.ToLower(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "missing_query",
			"message": "Query parameter 'q' is required",
		})
		return
	}

	var results []*vietqr.Bank
	for _, bank := range vietqr.Banks {
		if strings.Contains(strings.ToLower(bank.ShortName), query) ||
			strings.Contains(strings.ToLower(bank.Name), query) ||
			strings.Contains(bank.BIN, query) ||
			strings.Contains(strings.ToLower(bank.Code), query) {
			results = append(results, bank)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"query":   query,
		"count":   len(results),
		"data":    results,
	})
}
