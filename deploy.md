# Hướng dẫn triển khai MaxQR API

## Triển khai mới

```bash
cd /www/wwwroot && git clone https://github.com/TechScapBot/maxqr.git && cd maxqr && docker-compose up -d
```

## Cập nhật

```bash
cd /www/wwwroot/maxqr && docker-compose down && git pull && docker-compose up -d --build
```

## Kiểm tra

```bash
curl http://localhost:8080/health
```

## Đổi port (ví dụ 9090)

```bash
cd /www/wwwroot/maxqr && PORT=9090 docker-compose up -d
```
