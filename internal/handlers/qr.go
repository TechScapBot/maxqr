package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maxqr-api/internal/cache"
	"github.com/maxqr-api/internal/qrgen"
	"github.com/maxqr-api/internal/vietqr"
)

// QRHandler handles QR code generation requests
type QRHandler struct {
	generator *qrgen.Generator
	cache     *cache.Cache
	useCache  bool
}

// NewQRHandler creates a new QR handler
func NewQRHandler(generator *qrgen.Generator, qrCache *cache.Cache, useCache bool) *QRHandler {
	return &QRHandler{
		generator: generator,
		cache:     qrCache,
		useCache:  useCache,
	}
}

// GenerateRequest represents a QR code generation request
type GenerateRequest struct {
	BankBin       string `json:"bank_bin" form:"bank_bin"`             // Bank BIN code (6 digits)
	BankCode      string `json:"bank_code" form:"bank_code"`           // Bank code (e.g., "VIETCOMBANK")
	AccountNumber string `json:"account_number" form:"account_number"` // Account number
	Amount        int64  `json:"amount" form:"amount"`                 // Amount in VND
	Message       string `json:"message" form:"message"`               // Transfer description
	AccountName   string `json:"account_name" form:"account_name"`     // Account holder name
	Size          string `json:"size" form:"size"`                     // QR size: small, medium, large, xlarge
	Format        string `json:"format" form:"format"`                 // Output format: png, base64, json
}

// GenerateResponse represents the API response
type GenerateResponse struct {
	Success     bool                `json:"success"`
	QRString    string              `json:"qr_string,omitempty"`
	QRImageURL  string              `json:"qr_image_url,omitempty"`
	Base64Image string              `json:"base64_image,omitempty"`
	Bank        *vietqr.Bank        `json:"bank,omitempty"`
	Transfer    *TransferDetails    `json:"transfer,omitempty"`
	Error       string              `json:"error,omitempty"`
}

// TransferDetails contains transfer information
type TransferDetails struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name,omitempty"`
	Amount        int64  `json:"amount"`
	Message       string `json:"message,omitempty"`
}

// Generate handles POST /api/v1/generate
func (h *QRHandler) Generate(c *gin.Context) {
	var req GenerateRequest

	// Bind request (supports both JSON and form data)
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenerateResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Validate and get bank info
	bank, err := h.resolveBank(req.BankBin, req.BankCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenerateResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate account number
	if req.AccountNumber == "" {
		c.JSON(http.StatusBadRequest, GenerateResponse{
			Success: false,
			Error:   "Account number is required",
		})
		return
	}

	// Validate amount (must be non-negative)
	if req.Amount < 0 {
		c.JSON(http.StatusBadRequest, GenerateResponse{
			Success: false,
			Error:   "Amount must be non-negative",
		})
		return
	}

	// Validate message length
	if len(req.Message) > 50 {
		c.JSON(http.StatusBadRequest, GenerateResponse{
			Success: false,
			Error:   "Message must not exceed 50 characters",
		})
		return
	}

	// Generate VietQR string
	qrString := vietqr.Encode(vietqr.TransferInfo{
		BankBin:       bank.BIN,
		AccountNumber: req.AccountNumber,
		Amount:        req.Amount,
		Message:       req.Message,
		MerchantName:  req.AccountName,
		IsDynamic:     req.Amount > 0,
	})

	// Determine output format
	format := strings.ToLower(req.Format)
	if format == "" {
		format = "json"
	}

	// Parse size
	size := qrgen.ParseSize(req.Size)

	response := GenerateResponse{
		Success:  true,
		QRString: qrString,
		Bank:     bank,
		Transfer: &TransferDetails{
			AccountNumber: req.AccountNumber,
			AccountName:   req.AccountName,
			Amount:        req.Amount,
			Message:       req.Message,
		},
	}

	switch format {
	case "png":
		h.servePNG(c, qrString, size)
		return
	case "base64":
		imgData, err := h.getOrGenerateQR(qrString, size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GenerateResponse{
				Success: false,
				Error:   "Failed to generate QR image",
			})
			return
		}
		response.Base64Image = "data:image/png;base64," + encodeBase64(imgData)
	}

	c.JSON(http.StatusOK, response)
}

// GenerateImage handles GET /api/v1/qr/:bank_bin/:account_number.png
func (h *QRHandler) GenerateImage(c *gin.Context) {
	bankBin := c.Param("bank_bin")
	accountNumber := c.Param("account_number")

	// Remove .png extension if present
	accountNumber = strings.TrimSuffix(accountNumber, ".png")

	// Get query parameters
	amountStr := c.DefaultQuery("amount", "0")
	message := c.DefaultQuery("message", "")
	sizeStr := c.DefaultQuery("size", "medium")

	// Parse amount
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		amount = 0
	}

	// Validate bank
	bank := vietqr.GetBankByBIN(bankBin)
	if bank == nil {
		// Try by short name
		bank = vietqr.GetBankByShortName(bankBin)
	}
	if bank == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_bank",
			"message": "Invalid bank BIN or code",
		})
		return
	}

	// Generate QR string
	qrString := vietqr.Encode(vietqr.TransferInfo{
		BankBin:       bank.BIN,
		AccountNumber: accountNumber,
		Amount:        amount,
		Message:       message,
		IsDynamic:     amount > 0,
	})

	// Parse size
	size := qrgen.ParseSize(sizeStr)

	h.servePNG(c, qrString, size)
}

// QuickGenerate handles GET /api/v1/quick with query params
func (h *QRHandler) QuickGenerate(c *gin.Context) {
	bankBin := c.Query("bank")
	accountNumber := c.Query("account")
	amountStr := c.DefaultQuery("amount", "0")
	message := c.DefaultQuery("message", "")
	format := c.DefaultQuery("format", "png")
	sizeStr := c.DefaultQuery("size", "medium")

	if bankBin == "" || accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_params",
			"message": "bank and account parameters are required",
		})
		return
	}

	// Parse amount
	amount, _ := strconv.ParseInt(amountStr, 10, 64)

	// Resolve bank
	bank := vietqr.GetBankByBIN(bankBin)
	if bank == nil {
		bank = vietqr.GetBankByShortName(bankBin)
	}
	if bank == nil {
		bank = vietqr.GetBankByCode(strings.ToUpper(bankBin))
	}
	if bank == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_bank",
			"message": "Invalid bank identifier. Use BIN code, short name, or bank code.",
		})
		return
	}

	// Generate QR string
	qrString := vietqr.Encode(vietqr.TransferInfo{
		BankBin:       bank.BIN,
		AccountNumber: accountNumber,
		Amount:        amount,
		Message:       message,
		IsDynamic:     amount > 0,
	})

	size := qrgen.ParseSize(sizeStr)

	if format == "png" || format == "image" {
		h.servePNG(c, qrString, size)
		return
	}

	// Return JSON with base64
	imgData, err := h.getOrGenerateQR(qrString, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "generation_failed",
			"message": "Failed to generate QR code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"qr_string": qrString,
		"base64":    "data:image/png;base64," + encodeBase64(imgData),
		"bank":      bank,
	})
}

// Decode handles POST /api/v1/decode
func (h *QRHandler) Decode(c *gin.Context) {
	var req struct {
		QRString string `json:"qr_string" form:"qr_string"`
	}

	if err := c.ShouldBind(&req); err != nil || req.QRString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "qr_string is required",
		})
		return
	}

	decoded, err := vietqr.Decode(req.QRString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "decode_failed",
			"message": err.Error(),
		})
		return
	}

	// Get bank info
	bank := vietqr.GetBankByBIN(decoded.BankBin)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"bank_bin":       decoded.BankBin,
			"bank":           bank,
			"account_number": decoded.AccountNumber,
			"amount":         decoded.Amount,
			"message":        decoded.Message,
			"currency":       decoded.Currency,
			"country":        decoded.Country,
			"merchant_name":  decoded.MerchantName,
			"merchant_city":  decoded.MerchantCity,
			"is_valid":       decoded.IsValid,
		},
	})
}

// servePNG serves a PNG image response
func (h *QRHandler) servePNG(c *gin.Context, qrString string, size qrgen.QRSize) {
	imgData, err := h.getOrGenerateQR(qrString, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "generation_failed",
			"message": "Failed to generate QR code",
		})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Length", strconv.Itoa(len(imgData)))
	c.Header("Cache-Control", "public, max-age=300") // 5 minutes browser cache
	c.Data(http.StatusOK, "image/png", imgData)
}

// getOrGenerateQR gets from cache or generates new QR
func (h *QRHandler) getOrGenerateQR(qrString string, size qrgen.QRSize) ([]byte, error) {
	cacheKey := qrgen.ContentHash(qrString, size)

	// Try cache first
	if h.useCache {
		if data, found := h.cache.Get(cacheKey); found {
			return data, nil
		}
	}

	// Generate new QR
	data, err := h.generator.GeneratePNG(qrString, size)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if h.useCache {
		h.cache.Set(cacheKey, data)
	}

	return data, nil
}

// resolveBank resolves bank from BIN or code
func (h *QRHandler) resolveBank(bankBin, bankCode string) (*vietqr.Bank, error) {
	var bank *vietqr.Bank

	if bankBin != "" {
		bank = vietqr.GetBankByBIN(bankBin)
		if bank == nil {
			return nil, fmt.Errorf("invalid bank BIN: %s", bankBin)
		}
	} else if bankCode != "" {
		bank = vietqr.GetBankByCode(strings.ToUpper(bankCode))
		if bank == nil {
			bank = vietqr.GetBankByShortName(bankCode)
		}
		if bank == nil {
			return nil, fmt.Errorf("invalid bank code: %s", bankCode)
		}
	} else {
		return nil, fmt.Errorf("bank_bin or bank_code is required")
	}

	return bank, nil
}

// encodeBase64 encodes bytes to base64 string using optimized standard library
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
