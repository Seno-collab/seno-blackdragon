# seno-blackdragon

Challenge Project Booking

## 1. Authentication & Authorization Enterprise-grade

- JWT access + refresh token with Redis (revoke/rotation).
- Multi-device login(device + user_id), role-based + permission-based authorization.
- Middleware: logging TraceID, rate-limit, IP blocking.
- Challenge: Stateless performance vs. real-time revocation.
  Flow:
  Login: 1. user & password 2. ensure device 3. create session 4. issue access Token

## 2. Booking & Payment Consistency

- Flow: booking → hold → payment → confirmation.
- Ensure idempotency: prevent duplicate bookings.
- DB transaction + Redis lock for race conditions (double booking).
- Challenge: balance speed & data safety (overselling rooms).

## 3. Distributed Cache & Invalidation

- Redis cache for property/room list.
- Handle cache invalidation when landlord updates room info.
- Use versioning keys (room:v{n}:{id}).
- Challenge: stale cache leading to data mismatch.

## 4. CI/CD Zero Downtime

- GitHub Actions: build → push GHCR → deploy.
- Auto migration + rollback if failed.
- Hot-reload service / blue-green deployment.
- Challenge: heavy ALTER TABLE migration may lock DB.

## 5. Logging & Monitoring (Enterprise-grade)

- Zap structured JSON logs, integrate with ELK / Grafana Loki.
- TraceID from FE → API → DB query.
- Alerts when API 5xx exceed threshold.
- Challenge: avoid noisy logs while keeping debug info.

## 6. Microservice-ready Refactor

- Current: monolith (internal/api, internal/service…).
- Task: split into Auth, Booking, Payment services.
- Communication via gRPC or message queue (Kafka/NATS).
- Challenge: ensure data consistency (saga pattern, eventual consistency).

# 1) Xoá file cũ (nếu sai)

sudo rm -f /etc/apt/sources.list.d/docker.list

# 2) Đảm bảo keyrings tồn tại & key có sẵn

sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
 sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# 3) Ghi lại docker.list với đúng cú pháp, trên MỘT dòng

ARCH=$(dpkg --print-architecture)
CODENAME=focal   # hoặc: CODENAME=$(. /etc/os-release && echo $VERSION_CODENAME)
echo "deb [arch=${ARCH} signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu ${CODENAME} stable" | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null

# 4) (Tuỳ chọn) loại CRLF nếu bạn từng chỉnh file bằng editor Windows

sudo sed -i 's/\r$//' /etc/apt/sources.list.d/docker.list

# 5) Kiểm tra lại nội dung (phải là một dòng duy nhất, không thừa khoảng trắng lạ)

cat -n /etc/apt/sources.list.d/docker.list

# 6) Cập nhật & cài đặt

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
