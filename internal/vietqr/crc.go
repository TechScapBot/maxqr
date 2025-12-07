package vietqr

// CRC-16/CCITT-FALSE algorithm for EMVCo QR Code
// Polynomial: 0x1021, Initial: 0xFFFF, XorOut: 0x0000

const (
	crcPoly   uint16 = 0x1021
	crcInit   uint16 = 0xFFFF
	crcXorOut uint16 = 0x0000
)

// CalculateCRC16 computes CRC-16/CCITT-FALSE checksum
// This is the standard used by EMVCo for QR code validation
func CalculateCRC16(data []byte) uint16 {
	crc := crcInit

	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ crcPoly
			} else {
				crc <<= 1
			}
		}
	}

	return crc ^ crcXorOut
}

// CRC16String returns the CRC as a 4-character uppercase hex string
// The input should already include "6304" at the end
func CRC16String(data string) string {
	crc := CalculateCRC16([]byte(data))
	return sprintf("%04X", crc)
}

// CRC16StringWithTag appends "6304" and calculates CRC
func CRC16StringWithTag(data string) string {
	dataWithCRCTag := data + "6304"
	crc := CalculateCRC16([]byte(dataWithCRCTag))
	return sprintf("%04X", crc)
}

// sprintf is a minimal implementation to avoid fmt package overhead
func sprintf(format string, v uint16) string {
	const hexDigits = "0123456789ABCDEF"
	result := make([]byte, 4)
	result[0] = hexDigits[(v>>12)&0xF]
	result[1] = hexDigits[(v>>8)&0xF]
	result[2] = hexDigits[(v>>4)&0xF]
	result[3] = hexDigits[v&0xF]
	return string(result)
}
