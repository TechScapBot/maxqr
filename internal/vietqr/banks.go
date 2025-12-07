package vietqr

// Bank represents a Vietnamese bank with its BIN code and information
type Bank struct {
	BIN       string `json:"bin"`
	Code      string `json:"code"`       // Short code (e.g., "VCB", "TCB")
	ShortName string `json:"short_name"` // Short name (e.g., "Vietcombank")
	Name      string `json:"name"`       // Full Vietnamese name
	Logo      string `json:"logo"`       // Logo URL
	SwiftCode string `json:"swift_code"` // SWIFT/BIC code
}

// Banks is a map of bank code to Bank information for O(1) lookup
var Banks = map[string]*Bank{
	"VIETINBANK": {BIN: "970415", Code: "CTG", ShortName: "VietinBank", Name: "Ngân hàng TMCP Công thương Việt Nam", SwiftCode: "ICBVVNVX"},
	"VIETCOMBANK": {BIN: "970436", Code: "VCB", ShortName: "Vietcombank", Name: "Ngân hàng TMCP Ngoại Thương Việt Nam", SwiftCode: "BFTVVNVX"},
	"BIDV": {BIN: "970418", Code: "BIDV", ShortName: "BIDV", Name: "Ngân hàng TMCP Đầu tư và Phát triển Việt Nam", SwiftCode: "BIDVVNVX"},
	"AGRIBANK": {BIN: "970405", Code: "VBA", ShortName: "Agribank", Name: "Ngân hàng Nông nghiệp và Phát triển Nông thôn Việt Nam", SwiftCode: "VBAAVNVX"},
	"OCB": {BIN: "970448", Code: "OCB", ShortName: "OCB", Name: "Ngân hàng TMCP Phương Đông", SwiftCode: "ORCOVNVX"},
	"MBBANK": {BIN: "970422", Code: "MB", ShortName: "MBBank", Name: "Ngân hàng TMCP Quân đội", SwiftCode: "MSCBVNVX"},
	"TECHCOMBANK": {BIN: "970407", Code: "TCB", ShortName: "Techcombank", Name: "Ngân hàng TMCP Kỹ thương Việt Nam", SwiftCode: "VTCBVNVX"},
	"ACB": {BIN: "970416", Code: "ACB", ShortName: "ACB", Name: "Ngân hàng TMCP Á Châu", SwiftCode: "ASCBVNVX"},
	"VPBANK": {BIN: "970432", Code: "VPB", ShortName: "VPBank", Name: "Ngân hàng TMCP Việt Nam Thịnh Vượng", SwiftCode: "VPBKVNVX"},
	"TPBANK": {BIN: "970423", Code: "TPB", ShortName: "TPBank", Name: "Ngân hàng TMCP Tiên Phong", SwiftCode: "TPBVVNVX"},
	"SACOMBANK": {BIN: "970403", Code: "STB", ShortName: "Sacombank", Name: "Ngân hàng TMCP Sài Gòn Thương Tín", SwiftCode: "SGTTVNVX"},
	"HDBANK": {BIN: "970437", Code: "HDB", ShortName: "HDBank", Name: "Ngân hàng TMCP Phát triển Thành phố Hồ Chí Minh", SwiftCode: "HABOREAD"},
	"VIETCAPITALBANK": {BIN: "970454", Code: "BVB", ShortName: "VietCapitalBank", Name: "Ngân hàng TMCP Bản Việt", SwiftCode: "VCBCVNVX"},
	"SCB": {BIN: "970429", Code: "SCB", ShortName: "SCB", Name: "Ngân hàng TMCP Sài Gòn", SwiftCode: "SACLVNVX"},
	"VIB": {BIN: "970441", Code: "VIB", ShortName: "VIB", Name: "Ngân hàng TMCP Quốc tế Việt Nam", SwiftCode: "VNIBVNVX"},
	"SHB": {BIN: "970443", Code: "SHB", ShortName: "SHB", Name: "Ngân hàng TMCP Sài Gòn - Hà Nội", SwiftCode: "SHBAVNVX"},
	"EXIMBANK": {BIN: "970431", Code: "EIB", ShortName: "Eximbank", Name: "Ngân hàng TMCP Xuất Nhập khẩu Việt Nam", SwiftCode: "EBVIVNVX"},
	"MSB": {BIN: "970426", Code: "MSB", ShortName: "MSB", Name: "Ngân hàng TMCP Hàng Hải Việt Nam", SwiftCode: "MCOBVNVX"},
	"CAKE": {BIN: "546034", Code: "CAKE", ShortName: "CAKE", Name: "TMCP Việt Nam Thịnh Vượng - Ngân hàng số CAKE", SwiftCode: ""},
	"UBANK": {BIN: "546035", Code: "Ubank", ShortName: "Ubank", Name: "TMCP Việt Nam Thịnh Vượng - Ubank", SwiftCode: ""},
	"VIETTELMONEY": {BIN: "971005", Code: "VTLMONEY", ShortName: "ViettelMoney", Name: "Tổng Công ty Dịch vụ số Viettel", SwiftCode: ""},
	"TIMO": {BIN: "963388", Code: "TIMO", ShortName: "Timo", Name: "Ngân hàng số Timo by Ban Viet Bank", SwiftCode: ""},
	"VNPTMONEY": {BIN: "971011", Code: "VNPTMONEY", ShortName: "VNPTMoney", Name: "VNPT Money", SwiftCode: ""},
	"SAIGONBANK": {BIN: "970400", Code: "SGB", ShortName: "SaigonBank", Name: "Ngân hàng TMCP Sài Gòn Công Thương", SwiftCode: "SBITVNVX"},
	"BACABANK": {BIN: "970409", Code: "BAB", ShortName: "BacABank", Name: "Ngân hàng TMCP Bắc Á", SwiftCode: "NASCVNVX"},
	"MOMO": {BIN: "971025", Code: "MOMO", ShortName: "MoMo", Name: "CTCP Dịch Vụ Di Động Trực Tuyến", SwiftCode: ""},
	"PVCOMBANK": {BIN: "970412", Code: "PVC", ShortName: "PVcomBank", Name: "Ngân hàng TMCP Đại Chúng Việt Nam", SwiftCode: "WLORVNVX"},
	"NCB": {BIN: "970419", Code: "NCB", ShortName: "NCB", Name: "Ngân hàng TMCP Quốc Dân", SwiftCode: "NVBAVNVX"},
	"SHINHANBANK": {BIN: "970424", Code: "SHBVN", ShortName: "ShinhanBank", Name: "Ngân hàng TNHH MTV Shinhan Việt Nam", SwiftCode: "SHBKVNVX"},
	"ABBANK": {BIN: "970425", Code: "ABB", ShortName: "ABBANK", Name: "Ngân hàng TMCP An Bình", SwiftCode: "ABBKVNVX"},
	"VIETABANK": {BIN: "970427", Code: "VAB", ShortName: "VietABank", Name: "Ngân hàng TMCP Việt Á", SwiftCode: "VNACVNVX"},
	"NAMABANK": {BIN: "970428", Code: "NAB", ShortName: "NamABank", Name: "Ngân hàng TMCP Nam Á", SwiftCode: "NAMAVNVX"},
	"PGBANK": {BIN: "970430", Code: "PGB", ShortName: "PGBank", Name: "Ngân hàng TMCP Thịnh vượng và Phát triển", SwiftCode: "PGBLVNVX"},
	"VIETBANK": {BIN: "970433", Code: "VBB", ShortName: "VietBank", Name: "Ngân hàng TMCP Việt Nam Thương Tín", SwiftCode: "VNTTVNVX"},
	"BAOVIETBANK": {BIN: "970438", Code: "BVB", ShortName: "BaoVietBank", Name: "Ngân hàng TMCP Bảo Việt", SwiftCode: "BVBVVNVX"},
	"SEABANK": {BIN: "970440", Code: "SSB", ShortName: "SeABank", Name: "Ngân hàng TMCP Đông Nam Á", SwiftCode: "SEAVVNVX"},
	"COOPBANK": {BIN: "970446", Code: "COOPBANK", ShortName: "COOPBANK", Name: "Ngân hàng Hợp tác xã Việt Nam", SwiftCode: ""},
	"LPBANK": {BIN: "970449", Code: "LPB", ShortName: "LPBank", Name: "Ngân hàng TMCP Lộc Phát Việt Nam", SwiftCode: "LVBKVNVX"},
	"KIENLONGBANK": {BIN: "970452", Code: "KLB", ShortName: "KienLongBank", Name: "Ngân hàng TMCP Kiên Long", SwiftCode: "KLBKVNVX"},
	"KBANK": {BIN: "668888", Code: "KBank", ShortName: "KBank", Name: "Ngân hàng Đại chúng TNHH Kasikornbank", SwiftCode: "KASIVNVX"},
	"HONGLEONG": {BIN: "970442", Code: "HLBVN", ShortName: "HongLeong", Name: "Ngân hàng TNHH MTV Hong Leong Việt Nam", SwiftCode: "HLBBVNVX"},
	"KEBHANAHN": {BIN: "970467", Code: "KEBHN", ShortName: "KEBHanaHN", Name: "Ngân hàng KEB Hana - Chi nhánh Hà Nội", SwiftCode: "KOEXVNVX"},
	"KEBHANAHCM": {BIN: "970466", Code: "KEBHCM", ShortName: "KEBHanaHCM", Name: "Ngân hàng KEB Hana - Chi nhánh TP. Hồ Chí Minh", SwiftCode: "KOEXVNVX"},
	"CITIBANK": {BIN: "533948", Code: "CITI", ShortName: "Citibank", Name: "Ngân hàng Citibank, N.A. - Chi nhánh Hà Nội", SwiftCode: "CABOREAD"},
	"CBBANK": {BIN: "970444", Code: "CBB", ShortName: "CBBank", Name: "Ngân hàng Thương mại TNHH MTV Xây dựng Việt Nam", SwiftCode: "GTBAVNVX"},
	"CIMB": {BIN: "422589", Code: "CIMB", ShortName: "CIMB", Name: "Ngân hàng TNHH MTV CIMB Việt Nam", SwiftCode: "CIABOREAD"},
	"DBSBANK": {BIN: "796500", Code: "DBS", ShortName: "DBSBank", Name: "DBS Bank Ltd - Chi nhánh TP. Hồ Chí Minh", SwiftCode: "DBSSVNVX"},
	"VBSP": {BIN: "999888", Code: "VBSP", ShortName: "VBSP", Name: "Ngân hàng Chính sách Xã hội", SwiftCode: "VABOREAD"},
	"GPBANK": {BIN: "970408", Code: "GPB", ShortName: "GPBank", Name: "Ngân hàng Thương mại TNHH MTV Dầu Khí Toàn Cầu", SwiftCode: "GABOREAD"},
	"KOOKMINHCM": {BIN: "970463", Code: "KBHCM", ShortName: "KookminHCM", Name: "Ngân hàng Kookmin - Chi nhánh TP. Hồ Chí Minh", SwiftCode: "CZNBVNVX"},
	"KOOKMINHN": {BIN: "970462", Code: "KBHN", ShortName: "KookminHN", Name: "Ngân hàng Kookmin - Chi nhánh Hà Nội", SwiftCode: "CZNBVNVX"},
	"WOORI": {BIN: "970457", Code: "WVN", ShortName: "Woori", Name: "Ngân hàng TNHH MTV Woori Việt Nam", SwiftCode: "HVBKVNVX"},
	"VRB": {BIN: "970421", Code: "VRB", ShortName: "VRB", Name: "Ngân hàng Liên doanh Việt - Nga", SwiftCode: "VABOREAD"},
	"HSBC": {BIN: "458761", Code: "HSBC", ShortName: "HSBC", Name: "Ngân hàng TNHH MTV HSBC (Việt Nam)", SwiftCode: "HSBCVNVX"},
	"IBKHN": {BIN: "970455", Code: "IBKHN", ShortName: "IBKHN", Name: "Ngân hàng Công nghiệp Hàn Quốc - Hà Nội", SwiftCode: "IBKOVNVX"},
	"IBKHCM": {BIN: "970456", Code: "IBKHCM", ShortName: "IBKHCM", Name: "Ngân hàng Công nghiệp Hàn Quốc - TP. Hồ Chí Minh", SwiftCode: "IBKOVNVX"},
	"INDOVINABANK": {BIN: "970434", Code: "IVB", ShortName: "IndovinaBank", Name: "Ngân hàng TNHH Indovina", SwiftCode: "IABOREAD"},
	"UNITEDOVERSEAS": {BIN: "970458", Code: "UOB", ShortName: "UnitedOverseas", Name: "Ngân hàng United Overseas - TP. Hồ Chí Minh", SwiftCode: "UOVBVNVX"},
	"NONGHYUP": {BIN: "801011", Code: "NHB", ShortName: "Nonghyup", Name: "Ngân hàng Nonghyup - Chi nhánh Hà Nội", SwiftCode: "NHBKVNVX"},
	"STANDARDCHARTERED": {BIN: "970410", Code: "SCVN", ShortName: "StandardChartered", Name: "Ngân hàng TNHH MTV Standard Chartered Việt Nam", SwiftCode: "SCBLVNVX"},
	"PUBLICBANK": {BIN: "970439", Code: "PBVN", ShortName: "PublicBank", Name: "Ngân hàng TNHH MTV Public Việt Nam", SwiftCode: "VIDPVNVX"},
}

// BanksByBIN is a map for quick lookup by BIN code
var BanksByBIN = make(map[string]*Bank)

// BanksByShortName is a map for quick lookup by short name
var BanksByShortName = make(map[string]*Bank)

func init() {
	// Build lookup maps
	for _, bank := range Banks {
		BanksByBIN[bank.BIN] = bank
		BanksByShortName[bank.ShortName] = bank
	}
}

// GetBankByBIN returns bank information by BIN code
func GetBankByBIN(bin string) *Bank {
	return BanksByBIN[bin]
}

// GetBankByCode returns bank information by bank code (e.g., "VIETCOMBANK")
func GetBankByCode(code string) *Bank {
	return Banks[code]
}

// GetBankByShortName returns bank information by short name (e.g., "Vietcombank")
func GetBankByShortName(shortName string) *Bank {
	return BanksByShortName[shortName]
}

// GetAllBanks returns a slice of all banks
func GetAllBanks() []*Bank {
	banks := make([]*Bank, 0, len(Banks))
	for _, bank := range Banks {
		banks = append(banks, bank)
	}
	return banks
}

// IsValidBIN checks if a BIN code is valid
func IsValidBIN(bin string) bool {
	_, exists := BanksByBIN[bin]
	return exists
}

// IsValidBankCode checks if a bank code is valid
func IsValidBankCode(code string) bool {
	_, exists := Banks[code]
	return exists
}
