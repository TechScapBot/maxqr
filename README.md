# MaxQR API - High Performance QR Code Generator

API hiệu suất cao để tạo mã QR thanh toán theo chuẩn EMVCo cho các ngân hàng Việt Nam.

**Demo:** https://maxqr.scapbot.net

## Tính năng

- **Hiệu suất cao**: Xử lý hàng nghìn request/giây
- **Chuẩn EMVCo**: Tuân thủ 100% tiêu chuẩn EMVCo của NAPAS
- **60+ ngân hàng**: Hỗ trợ tất cả ngân hàng Việt Nam
- **Không cần API Key**: Sử dụng ngay, không cần đăng ký

## API Endpoints

### Danh sách ngân hàng

```bash
# Lấy tất cả ngân hàng
GET /api/v1/banks

# Tìm kiếm ngân hàng
GET /api/v1/banks/search?q=vietcom

# Lấy thông tin ngân hàng
GET /api/v1/banks/970436
GET /api/v1/banks/VIETCOMBANK
GET /api/v1/banks/Vietcombank
```

### Tạo mã QR

#### 1. Quick Generate (đơn giản nhất)

```bash
# Tạo QR ảnh PNG
GET /api/v1/quick?bank=970436&account=1234567890&amount=100000&message=Thanh%20toan

# Tham số:
# - bank: BIN code, short name hoặc bank code (bắt buộc)
# - account: Số tài khoản (bắt buộc)
# - amount: Số tiền VND (tùy chọn, mặc định 0)
# - message: Nội dung chuyển khoản (tùy chọn, tối đa 50 ký tự)
# - size: small/medium/large/xlarge (tùy chọn, mặc định medium)
# - format: png/json (tùy chọn, mặc định png)
```

#### 2. Generate Image (URL trực tiếp)

```bash
# Nhúng trực tiếp vào <img src="">
GET /api/v1/qr/970436/1234567890.png?amount=100000&message=Thanh%20toan
```

#### 3. Generate API (đầy đủ tùy chọn)

```bash
POST /api/v1/generate
Content-Type: application/json

{
  "bank_bin": "970436",        // hoặc "bank_code": "VIETCOMBANK"
  "account_number": "1234567890",
  "amount": 100000,
  "message": "Thanh toan don hang #123",
  "account_name": "NGUYEN VAN A",
  "size": "large",
  "format": "json"             // json, png, base64
}
```

**Response:**

```json
{
  "success": true,
  "qr_string": "00020101021138530010A000000727...",
  "bank": {
    "bin": "970436",
    "code": "VCB",
    "short_name": "Vietcombank",
    "name": "Ngân hàng TMCP Ngoại Thương Việt Nam"
  },
  "transfer": {
    "account_number": "1234567890",
    "account_name": "NGUYEN VAN A",
    "amount": 100000,
    "message": "Thanh toan don hang #123"
  },
  "base64_image": "data:image/png;base64,..."
}
```

### Giải mã QR

```bash
POST /api/v1/decode
Content-Type: application/json

{
  "qr_string": "00020101021138530010A000000727..."
}
```

## Ví dụ sử dụng

### HTML - Nhúng QR vào website

```html
<!-- Cách 1: Trực tiếp dùng img src -->
<img src="https://maxqr.scapbot.net/api/v1/qr/970436/1234567890.png?amount=100000&message=Thanh%20toan"
     alt="MaxQR Payment" />

<!-- Cách 2: Dùng URL động -->
<img id="qr-code" />
<script>
const qrUrl = new URL('https://maxqr.scapbot.net/api/v1/quick');
qrUrl.searchParams.set('bank', '970436');
qrUrl.searchParams.set('account', '1234567890');
qrUrl.searchParams.set('amount', '100000');
qrUrl.searchParams.set('message', 'Thanh toan');
qrUrl.searchParams.set('format', 'png');
document.getElementById('qr-code').src = qrUrl.toString();
</script>
```

### JavaScript/Node.js

```javascript
// Tạo QR với fetch (không cần API Key)
const response = await fetch('https://maxqr.scapbot.net/api/v1/generate', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    bank_code: 'VIETCOMBANK',
    account_number: '1234567890',
    amount: 100000,
    message: 'Thanh toan'
  })
});

const data = await response.json();
console.log(data.qr_string);
console.log(data.base64_image);
```

### Python

```python
import requests

# Tạo QR
response = requests.post('https://maxqr.scapbot.net/api/v1/generate', json={
    'bank_bin': '970436',
    'account_number': '1234567890',
    'amount': 100000,
    'message': 'Thanh toan'
})

data = response.json()
print(data['qr_string'])

# Tải ảnh QR
img_response = requests.get(
    'https://maxqr.scapbot.net/api/v1/quick',
    params={
        'bank': '970436',
        'account': '1234567890',
        'amount': 100000
    }
)
with open('qr.png', 'wb') as f:
    f.write(img_response.content)
```

### cURL

```bash
# Tạo QR và lấy JSON
curl -X POST https://maxqr.scapbot.net/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "bank_code": "VIETCOMBANK",
    "account_number": "1234567890",
    "amount": 100000,
    "message": "Thanh toan"
  }'

# Tải ảnh QR
curl -o qr.png "https://maxqr.scapbot.net/api/v1/quick?bank=970436&account=1234567890&amount=100000"

# Lấy danh sách ngân hàng
curl https://maxqr.scapbot.net/api/v1/banks
```

## Danh sách ngân hàng hỗ trợ

| BIN | Tên ngắn | Tên đầy đủ |
|-----|----------|------------|
| 970415 | VietinBank | Ngân hàng TMCP Công thương Việt Nam |
| 970436 | Vietcombank | Ngân hàng TMCP Ngoại Thương Việt Nam |
| 970418 | BIDV | Ngân hàng TMCP Đầu tư và Phát triển Việt Nam |
| 970405 | Agribank | Ngân hàng Nông nghiệp và Phát triển Nông thôn |
| 970422 | MBBank | Ngân hàng TMCP Quân đội |
| 970407 | Techcombank | Ngân hàng TMCP Kỹ thương Việt Nam |
| 970416 | ACB | Ngân hàng TMCP Á Châu |
| 970432 | VPBank | Ngân hàng TMCP Việt Nam Thịnh Vượng |
| 970423 | TPBank | Ngân hàng TMCP Tiên Phong |
| 970443 | SHB | Ngân hàng TMCP Sài Gòn - Hà Nội |
| ... | ... | Và 50+ ngân hàng khác |

Xem đầy đủ tại: `GET /api/v1/banks`

## License

MIT License
