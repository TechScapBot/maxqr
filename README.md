# MaxQR API - High Performance QR Code Generator

API hiệu suất cao để tạo mã QR thanh toán theo chuẩn EMVCo cho các ngân hàng Việt Nam.

## Tính năng

- **Hiệu suất cao**: Xử lý hàng nghìn request/giây
- **Chuẩn EMVCo**: Tuân thủ 100% tiêu chuẩn EMVCo của NAPAS
- **60+ ngân hàng**: Hỗ trợ tất cả ngân hàng Việt Nam
- **Không cần API Key**: Sử dụng ngay, không cần đăng ký
- **Caching thông minh**: In-memory cache giảm latency
- **Rate limiting**: Bảo vệ khỏi DDoS
- **Docker ready**: Deploy nhanh chóng

## Cài đặt

### Chạy với Docker (Khuyến nghị)

```bash
# Clone repository
git clone https://github.com/your-repo/maxqr-api
cd maxqr-api

# Chạy với docker-compose
docker-compose up -d

# API sẽ chạy tại http://localhost:8080
```

### Chạy từ source

```bash
# Cài đặt dependencies
go mod download

# Chạy server
go run cmd/server/main.go

# Hoặc build binary
go build -o maxqr-api cmd/server/main.go
./maxqr-api
```

## API Key

**Không bắt buộc!** API hoạt động public mặc định. Bạn có thể sử dụng ngay mà không cần API Key.

Nếu muốn bảo mật, cấu hình biến môi trường `API_KEY` và gửi kèm header `X-API-Key` trong request.

## API Endpoints

### Health Check

```bash
GET /health
GET /ready
GET /stats
```

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
<img src="http://localhost:8080/api/v1/qr/970436/1234567890.png?amount=100000&message=Thanh%20toan"
     alt="MaxQR Payment" />

<!-- Cách 2: Dùng URL động -->
<img id="qr-code" />
<script>
const qrUrl = new URL('http://localhost:8080/api/v1/quick');
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
const response = await fetch('http://localhost:8080/api/v1/generate', {
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
response = requests.post('http://localhost:8080/api/v1/generate', json={
    'bank_bin': '970436',
    'account_number': '1234567890',
    'amount': 100000,
    'message': 'Thanh toan'
})

data = response.json()
print(data['qr_string'])

# Tải ảnh QR
img_response = requests.get(
    'http://localhost:8080/api/v1/quick',
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
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "bank_code": "VIETCOMBANK",
    "account_number": "1234567890",
    "amount": 100000,
    "message": "Thanh toan"
  }'

# Tải ảnh QR
curl -o qr.png "http://localhost:8080/api/v1/quick?bank=970436&account=1234567890&amount=100000"

# Lấy danh sách ngân hàng
curl http://localhost:8080/api/v1/banks
```

## Cấu hình

Sử dụng biến môi trường hoặc file `.env`:

```bash
# Server
PORT=8080
HOST=0.0.0.0

# Security (để trống = public API, không cần API Key)
API_KEY=
ALLOWED_ORIGINS=*

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_PER_SECOND=100
RATE_LIMIT_BURST=200

# Cache
CACHE_ENABLED=true
CACHE_MAX_SIZE_MB=100
CACHE_DEFAULT_TTL=5m

# Logging
LOG_LEVEL=info
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

## Performance Benchmarks

Kết quả benchmark trên máy 4 cores, 8GB RAM:

```
Benchmark                    Requests/sec    Avg Latency
----------------------------------------------------------
QR String Generation         50,000+         < 1ms
QR Image (cache hit)         30,000+         < 2ms
QR Image (cache miss)        10,000+         < 5ms
Bank Lookup                  100,000+        < 0.1ms
```

## Kiến trúc

```
┌─────────────────────────────────────────────────────────┐
│                     MaxQR API                           │
├─────────────────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │  Gin    │─▶│ Rate Limit  │─▶│  Security Headers   │ │
│  │ Router  │  │ Middleware  │  │     Middleware      │ │
│  └─────────┘  └─────────────┘  └─────────────────────┘ │
│        │                                                │
│        ▼                                                │
│  ┌─────────────────────────────────────────────────┐   │
│  │                   Handlers                       │   │
│  │  ┌──────────┐  ┌───────────┐  ┌──────────────┐  │   │
│  │  │ QR Gen   │  │   Banks   │  │    Health    │  │   │
│  │  └──────────┘  └───────────┘  └──────────────┘  │   │
│  └─────────────────────────────────────────────────┘   │
│        │                                                │
│        ▼                                                │
│  ┌─────────────────────────────────────────────────┐   │
│  │                  Core Engine                     │   │
│  │  ┌──────────┐  ┌───────────┐  ┌──────────────┐  │   │
│  │  │ EMVCo    │  │ QR Image  │  │  In-Memory   │  │   │
│  │  │ Encoder  │  │ Generator │  │    Cache     │  │   │
│  │  └──────────┘  └───────────┘  └──────────────┘  │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## License

MIT License
