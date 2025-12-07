package vietqr

import (
	"errors"
	"strconv"
	"strings"
)

// Common errors
var (
	ErrInvalidQRString   = errors.New("invalid VietQR string")
	ErrInvalidCRC        = errors.New("CRC validation failed")
	ErrInvalidLength     = errors.New("invalid field length")
	ErrMissingField      = errors.New("missing required field")
	ErrUnsupportedFormat = errors.New("unsupported QR format")
)

// DecodedInfo contains parsed VietQR information
type DecodedInfo struct {
	PayloadFormat    string // Should be "01"
	InitiationMethod string // "11" (static) or "12" (dynamic)
	BankBin          string // 6-digit bank BIN
	AccountNumber    string // Recipient's account number
	Amount           int64  // Transfer amount (0 if not specified)
	Currency         string // Currency code (usually "704" for VND)
	Country          string // Country code (usually "VN")
	MerchantName     string // Optional merchant name
	MerchantCity     string // Optional city
	Message          string // Transfer description
	CRC              string // CRC checksum
	IsValid          bool   // Whether CRC validation passed
}

// Decoder parses VietQR strings
type Decoder struct {
	data   string
	pos    int
	length int
}

// NewDecoder creates a new VietQR decoder
func NewDecoder(qrString string) *Decoder {
	return &Decoder{
		data:   qrString,
		pos:    0,
		length: len(qrString),
	}
}

// Decode parses a VietQR string and returns decoded information
func (d *Decoder) Decode() (*DecodedInfo, error) {
	if d.length < 8 {
		return nil, ErrInvalidQRString
	}

	info := &DecodedInfo{}

	// Validate CRC first
	if d.length >= 8 {
		crcPos := d.length - 4
		if d.data[crcPos-4:crcPos] == "6304" {
			expectedCRC := CRC16String(d.data[:crcPos])
			actualCRC := d.data[crcPos:]
			info.CRC = actualCRC
			info.IsValid = strings.EqualFold(expectedCRC, actualCRC)
		}
	}

	// Parse all TLV fields
	for d.pos < d.length {
		tag, value, err := d.readTLV()
		if err != nil {
			break
		}

		switch tag {
		case TagPayloadFormat:
			info.PayloadFormat = value
		case TagInitiationMethod:
			info.InitiationMethod = value
		case TagMerchantAccount:
			d.parseMerchantAccount(value, info)
		case TagCurrency:
			info.Currency = value
		case TagAmount:
			if amt, err := strconv.ParseInt(value, 10, 64); err == nil {
				info.Amount = amt
			}
		case TagCountry:
			info.Country = value
		case TagMerchantName:
			info.MerchantName = value
		case TagMerchantCity:
			info.MerchantCity = value
		case TagAdditionalData:
			d.parseAdditionalData(value, info)
		case TagCRC:
			info.CRC = value
		}
	}

	return info, nil
}

// readTLV reads a single Tag-Length-Value triplet
func (d *Decoder) readTLV() (tag, value string, err error) {
	if d.pos+4 > d.length {
		return "", "", ErrInvalidQRString
	}

	tag = d.data[d.pos : d.pos+2]
	d.pos += 2

	lengthStr := d.data[d.pos : d.pos+2]
	d.pos += 2

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", "", ErrInvalidLength
	}

	if d.pos+length > d.length {
		return "", "", ErrInvalidLength
	}

	value = d.data[d.pos : d.pos+length]
	d.pos += length

	return tag, value, nil
}

// parseMerchantAccount parses nested TLV in tag 38
func (d *Decoder) parseMerchantAccount(data string, info *DecodedInfo) {
	decoder := NewDecoder(data)
	decoder.length = len(data)

	for decoder.pos < decoder.length {
		tag, value, err := decoder.readTLV()
		if err != nil {
			break
		}

		switch tag {
		case "01": // Beneficiary Organization
			d.parseBeneficiaryOrg(value, info)
		}
	}
}

// parseBeneficiaryOrg parses the beneficiary organization structure
func (d *Decoder) parseBeneficiaryOrg(data string, info *DecodedInfo) {
	decoder := NewDecoder(data)
	decoder.length = len(data)

	for decoder.pos < decoder.length {
		tag, value, err := decoder.readTLV()
		if err != nil {
			break
		}

		switch tag {
		case "00": // Bank BIN
			info.BankBin = value
		case "01": // Account Number
			info.AccountNumber = value
		}
	}
}

// parseAdditionalData parses nested TLV in tag 62
func (d *Decoder) parseAdditionalData(data string, info *DecodedInfo) {
	decoder := NewDecoder(data)
	decoder.length = len(data)

	for decoder.pos < decoder.length {
		tag, value, err := decoder.readTLV()
		if err != nil {
			break
		}

		switch tag {
		case SubTagPurpose: // Purpose of Transaction (Message)
			info.Message = value
		}
	}
}

// Decode is a package-level convenience function
func Decode(qrString string) (*DecodedInfo, error) {
	return NewDecoder(qrString).Decode()
}

// ValidateCRC checks if the QR string has a valid CRC
func ValidateCRC(qrString string) bool {
	if len(qrString) < 8 {
		return false
	}

	// CRC position is at the end, 4 characters
	crcPos := len(qrString) - 4

	// Check that "6304" exists before the CRC value
	if crcPos < 4 || qrString[crcPos-4:crcPos] != "6304" {
		return false
	}

	// CRC is calculated over the entire string including "6304" tag
	dataForCRC := qrString[:crcPos] // Everything up to (but not including) the CRC value
	expectedCRC := CalculateCRC16([]byte(dataForCRC))
	actualCRC := qrString[crcPos:]

	return strings.EqualFold(sprintf("%04X", expectedCRC), actualCRC)
}
