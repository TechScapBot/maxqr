package vietqr

import (
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		info     TransferInfo
		wantLen  int // Minimum length
		contains []string
	}{
		{
			name: "Basic transfer without amount",
			info: TransferInfo{
				BankBin:       "970436", // Vietcombank
				AccountNumber: "1234567890",
			},
			wantLen: 80,
			contains: []string{
				"000201",            // Payload format
				"010211",            // Static QR
				"A000000727",        // NAPAS GUID
				"970436",            // Bank BIN
				"1234567890",        // Account number
				"5303704",           // Currency VND
				"5802VN",            // Country
				"6304",              // CRC tag
			},
		},
		{
			name: "Transfer with amount",
			info: TransferInfo{
				BankBin:       "970422", // MBBank
				AccountNumber: "9876543210",
				Amount:        100000,
			},
			wantLen: 90,
			contains: []string{
				"000201",
				"010212",            // Dynamic QR
				"970422",
				"9876543210",
				"54",                // Amount tag
				"100000",
				"6304",
			},
		},
		{
			name: "Transfer with message",
			info: TransferInfo{
				BankBin:       "970407", // Techcombank
				AccountNumber: "123456789",
				Amount:        50000,
				Message:       "Thanh toan",
			},
			wantLen: 100,
			contains: []string{
				"970407",
				"62",                // Additional data tag
				"Thanh toan",
			},
		},
		{
			name: "Transfer with Vietnamese message (should be converted)",
			info: TransferInfo{
				BankBin:       "970436",
				AccountNumber: "1234567890",
				Amount:        100000,
				Message:       "Cảm ơn",
			},
			wantLen: 90,
			contains: []string{
				"Cam on", // Diacritics removed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.info)

			if len(result) < tt.wantLen {
				t.Errorf("Encode() length = %d, want at least %d", len(result), tt.wantLen)
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Encode() = %s, should contain %s", result, substr)
				}
			}

			// Verify CRC
			if !ValidateCRC(result) {
				t.Errorf("Encode() generated invalid CRC")
			}
		})
	}
}

func TestCRC16(t *testing.T) {
	// Test with known valid VietQR string
	qrWithoutCRC := "00020101021138530010A00000072701210006970407010797968680208QRIBFTTA53037045802VN62240820gen by sunary/vietqr"
	result := CRC16StringWithTag(qrWithoutCRC)

	// Verify that the result is a valid 4-character hex string
	if len(result) != 4 {
		t.Errorf("CRC16StringWithTag() length = %d, want 4", len(result))
	}

	// Verify CRC validation works for generated QR
	fullQR := qrWithoutCRC + "6304" + result
	if !ValidateCRC(fullQR) {
		t.Errorf("Generated CRC should be valid, got %s", result)
	}
}

func TestValidateCRC(t *testing.T) {
	// Valid QR string
	validQR := Encode(TransferInfo{
		BankBin:       "970436",
		AccountNumber: "1234567890",
		Amount:        100000,
	})

	if !ValidateCRC(validQR) {
		t.Error("ValidateCRC() should return true for valid QR")
	}

	// Invalid QR (tampered)
	invalidQR := validQR[:len(validQR)-1] + "X"
	if ValidateCRC(invalidQR) {
		t.Error("ValidateCRC() should return false for invalid QR")
	}
}

func TestRemoveVietnameseDiacritics(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Cảm ơn", "Cam on"},
		{"Nguyễn Văn A", "Nguyen Van A"},
		{"Đà Nẵng", "Da Nang"},
		{"Hồ Chí Minh", "Ho Chi Minh"},
		{"Thanh toán đơn hàng", "Thanh toan don hang"},
	}

	for _, tt := range tests {
		result := removeVietnameseDiacritics(tt.input)
		if result != tt.expected {
			t.Errorf("removeVietnameseDiacritics(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	info := TransferInfo{
		BankBin:       "970436",
		AccountNumber: "1234567890",
		Amount:        100000,
		Message:       "Thanh toan don hang",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(info)
	}
}

func BenchmarkEncodeParallel(b *testing.B) {
	info := TransferInfo{
		BankBin:       "970436",
		AccountNumber: "1234567890",
		Amount:        100000,
		Message:       "Thanh toan don hang",
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Encode(info)
		}
	})
}
