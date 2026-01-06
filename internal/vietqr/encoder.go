package vietqr

import (
	"strconv"
	"strings"
)

// EMVCo Tag IDs
const (
	TagPayloadFormat      = "00" // Payload Format Indicator
	TagInitiationMethod   = "01" // Point of Initiation Method
	TagMerchantAccount    = "38" // Merchant Account Information (VietQR specific)
	TagMerchantCategory   = "52" // Merchant Category Code
	TagCurrency           = "53" // Transaction Currency
	TagAmount             = "54" // Transaction Amount
	TagTipIndicator       = "55" // Tip or Convenience Indicator
	TagCountry            = "58" // Country Code
	TagMerchantName       = "59" // Merchant Name
	TagMerchantCity       = "60" // Merchant City
	TagPostalCode         = "61" // Postal Code
	TagAdditionalData     = "62" // Additional Data Field Template
	TagCRC                = "63" // CRC
)

// Sub-tags within Tag 38 (Merchant Account Information)
const (
	SubTagGUID          = "00" // Globally Unique Identifier
	SubTagBankBin       = "01" // Bank BIN
	SubTagAccountNumber = "02" // Account Number
	SubTagTransferMethod = "03" // Transfer Method (QRIBFTTA/QRIBFTTC - optional, defaults to QRIBFTTA)
)

// Sub-tags within Tag 62 (Additional Data)
const (
	SubTagBillNumber     = "01" // Bill Number
	SubTagMobileNumber   = "02" // Mobile Number
	SubTagStoreLabel     = "03" // Store Label
	SubTagLoyaltyNumber  = "04" // Loyalty Number
	SubTagReferenceLabel = "05" // Reference Label
	SubTagCustomerLabel  = "06" // Customer Label
	SubTagTerminalLabel  = "07" // Terminal Label
	SubTagPurpose        = "08" // Purpose of Transaction (Description/Message)
)

// Constants
const (
	PayloadFormatEMV    = "01"          // EMV QR Code version
	InitiationStatic    = "11"          // Static QR (reusable)
	InitiationDynamic   = "12"          // Dynamic QR (single use)
	CurrencyVND         = "704"         // ISO 4217 for Vietnamese Dong
	CountryVN           = "VN"          // ISO 3166-1 alpha-2 for Vietnam
	NAPAS_GUID          = "A000000727"  // NAPAS provider identifier
	TransferToAccount   = "QRIBFTTA"    // Transfer to Account
	TransferToCard      = "QRIBFTTC"    // Transfer to Card
)

// TransferInfo contains all information needed to generate a VietQR code
type TransferInfo struct {
	BankBin       string // 6-digit bank BIN code (e.g., "970436" for Vietcombank)
	AccountNumber string // Recipient's bank account number
	Amount        int64  // Transfer amount in VND (0 for static QR)
	Message       string // Transfer description/message (max 50 chars)
	MerchantName  string // Optional: Merchant/recipient name (max 25 chars)
	MerchantCity  string // Optional: City (default: "Ha Noi")
	IsDynamic     bool   // true for single-use QR, false for reusable
	Editable      bool   // true to allow user to edit amount/message when scanning
}

// Encoder generates VietQR strings following EMVCo standard
type Encoder struct {
	builder strings.Builder
}

// NewEncoder creates a new VietQR encoder
func NewEncoder() *Encoder {
	return &Encoder{
		builder: strings.Builder{},
	}
}

// Encode generates a complete VietQR string from transfer information
func (e *Encoder) Encode(info TransferInfo) string {
	e.builder.Reset()
	e.builder.Grow(256) // Pre-allocate for performance

	// 00 - Payload Format Indicator (mandatory)
	e.appendTLV(TagPayloadFormat, PayloadFormatEMV)

	// 01 - Point of Initiation Method (mandatory)
	// Static (11): Allows user to edit amount/message when scanning
	// Dynamic (12): Fixed amount/message, cannot be edited
	if info.Editable {
		// Force static QR to allow editing even with preset amount/message
		e.appendTLV(TagInitiationMethod, InitiationStatic)
	} else if info.IsDynamic || info.Amount > 0 {
		e.appendTLV(TagInitiationMethod, InitiationDynamic)
	} else {
		e.appendTLV(TagInitiationMethod, InitiationStatic)
	}

	// 38 - Merchant Account Information (VietQR/NAPAS specific)
	merchantAccount := e.buildMerchantAccount(info.BankBin, info.AccountNumber)
	e.appendTLV(TagMerchantAccount, merchantAccount)

	// 52 - Merchant Category Code (optional, using default)
	// Not included for basic transfer QR

	// 53 - Transaction Currency (mandatory)
	e.appendTLV(TagCurrency, CurrencyVND)

	// 54 - Transaction Amount (conditional)
	if info.Amount > 0 {
		e.appendTLV(TagAmount, strconv.FormatInt(info.Amount, 10))
	}

	// 58 - Country Code (mandatory)
	e.appendTLV(TagCountry, CountryVN)

	// 59 - Merchant Name (optional)
	if info.MerchantName != "" {
		name := truncateString(removeVietnameseDiacritics(info.MerchantName), 25)
		e.appendTLV(TagMerchantName, name)
	}

	// 60 - Merchant City (optional)
	city := info.MerchantCity
	if city == "" {
		city = "Ha Noi"
	}
	city = truncateString(removeVietnameseDiacritics(city), 15)
	e.appendTLV(TagMerchantCity, city)

	// 62 - Additional Data Field (conditional - for message/description)
	if info.Message != "" {
		additionalData := e.buildAdditionalData(info.Message)
		e.appendTLV(TagAdditionalData, additionalData)
	}

	// 63 - CRC (mandatory - calculated last)
	qrWithoutCRC := e.builder.String()
	crc := CRC16StringWithTag(qrWithoutCRC)
	e.builder.WriteString(TagCRC)
	e.builder.WriteString("04")
	e.builder.WriteString(crc)

	return e.builder.String()
}

// buildMerchantAccount creates the nested TLV structure for tag 38
func (e *Encoder) buildMerchantAccount(bankBin, accountNumber string) string {
	var sb strings.Builder
	sb.Grow(64)

	// 00 - GUID (NAPAS identifier)
	appendTLVTo(&sb, SubTagGUID, NAPAS_GUID)

	// 01 - Beneficiary Organization (contains bank BIN and account)
	beneficiary := e.buildBeneficiaryOrg(bankBin, accountNumber)
	appendTLVTo(&sb, "01", beneficiary)

	// 02 - Service Code (QRIBFTTA - Transfer to Account)
	appendTLVTo(&sb, "02", TransferToAccount)

	return sb.String()
}

// buildBeneficiaryOrg creates the nested structure for beneficiary info
func (e *Encoder) buildBeneficiaryOrg(bankBin, accountNumber string) string {
	var sb strings.Builder
	sb.Grow(32)

	// 00 - Bank BIN (6 digits)
	appendTLVTo(&sb, "00", bankBin)

	// 01 - Account Number
	appendTLVTo(&sb, "01", accountNumber)

	return sb.String()
}

// buildAdditionalData creates the nested TLV structure for tag 62
func (e *Encoder) buildAdditionalData(message string) string {
	var sb strings.Builder
	sb.Grow(64)

	// 08 - Purpose of Transaction (Description)
	msg := truncateString(removeVietnameseDiacritics(message), 50)
	appendTLVTo(&sb, SubTagPurpose, msg)

	return sb.String()
}

// appendTLV adds a Tag-Length-Value triplet to the encoder's builder
func (e *Encoder) appendTLV(tag, value string) {
	appendTLVTo(&e.builder, tag, value)
}

// appendTLVTo adds a TLV triplet to any strings.Builder
func appendTLVTo(sb *strings.Builder, tag, value string) {
	sb.WriteString(tag)
	sb.WriteString(formatLength(len(value)))
	sb.WriteString(value)
}

// formatLength formats length as 2-digit string
func formatLength(length int) string {
	if length < 10 {
		return "0" + strconv.Itoa(length)
	}
	return strconv.Itoa(length)
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}

// removeVietnameseDiacritics converts Vietnamese characters to ASCII
// This is required because QR code readers may not support Unicode properly
func removeVietnameseDiacritics(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if mapped, ok := vietnameseToASCII[r]; ok {
			result.WriteRune(mapped)
		} else if r < 128 {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ') // Replace unknown Unicode with space
		}
	}

	return result.String()
}

// vietnameseToASCII maps Vietnamese characters to their ASCII equivalents
var vietnameseToASCII = map[rune]rune{
	// Lowercase
	'á': 'a', 'à': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a',
	'ă': 'a', 'ắ': 'a', 'ằ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
	'â': 'a', 'ấ': 'a', 'ầ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
	'é': 'e', 'è': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e',
	'ê': 'e', 'ế': 'e', 'ề': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
	'í': 'i', 'ì': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
	'ó': 'o', 'ò': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o',
	'ô': 'o', 'ố': 'o', 'ồ': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
	'ơ': 'o', 'ớ': 'o', 'ờ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
	'ú': 'u', 'ù': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u',
	'ư': 'u', 'ứ': 'u', 'ừ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
	'ý': 'y', 'ỳ': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
	'đ': 'd',
	// Uppercase
	'Á': 'A', 'À': 'A', 'Ả': 'A', 'Ã': 'A', 'Ạ': 'A',
	'Ă': 'A', 'Ắ': 'A', 'Ằ': 'A', 'Ẳ': 'A', 'Ẵ': 'A', 'Ặ': 'A',
	'Â': 'A', 'Ấ': 'A', 'Ầ': 'A', 'Ẩ': 'A', 'Ẫ': 'A', 'Ậ': 'A',
	'É': 'E', 'È': 'E', 'Ẻ': 'E', 'Ẽ': 'E', 'Ẹ': 'E',
	'Ê': 'E', 'Ế': 'E', 'Ề': 'E', 'Ể': 'E', 'Ễ': 'E', 'Ệ': 'E',
	'Í': 'I', 'Ì': 'I', 'Ỉ': 'I', 'Ĩ': 'I', 'Ị': 'I',
	'Ó': 'O', 'Ò': 'O', 'Ỏ': 'O', 'Õ': 'O', 'Ọ': 'O',
	'Ô': 'O', 'Ố': 'O', 'Ồ': 'O', 'Ổ': 'O', 'Ỗ': 'O', 'Ộ': 'O',
	'Ơ': 'O', 'Ớ': 'O', 'Ờ': 'O', 'Ở': 'O', 'Ỡ': 'O', 'Ợ': 'O',
	'Ú': 'U', 'Ù': 'U', 'Ủ': 'U', 'Ũ': 'U', 'Ụ': 'U',
	'Ư': 'U', 'Ứ': 'U', 'Ừ': 'U', 'Ử': 'U', 'Ữ': 'U', 'Ự': 'U',
	'Ý': 'Y', 'Ỳ': 'Y', 'Ỷ': 'Y', 'Ỹ': 'Y', 'Ỵ': 'Y',
	'Đ': 'D',
}

// Encode is a package-level convenience function
func Encode(info TransferInfo) string {
	return NewEncoder().Encode(info)
}
