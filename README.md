# Marketing Automation Backend

Dự án Backend cho hệ thống Marketing Automation, được xây dựng bằng ngôn ngữ **Go (Golang)** kết hợp với framework **Gin** và **PostgreSQL**.

## 🚀 Công nghệ sử dụng
- **Ngôn ngữ:** Go (Golang) 1.18+
- **Web Framework:** [Gin-Gonic](https://gin-gonic.com/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL
- **Bảo mật:** bcrypt (Mã hóa mật khẩu)
- **Quản lý cấu hình:** godotenv (Đọc file `.env`)

## 📁 Cấu trúc thư mục hiện tại
```text
.
├── config/
│   └── database.go         # Cấu hình và kết nối PostgreSQL
├── controllers/
│   └── user_controller.go  # Xử lý logic cho các API liên quan đến User
├── models/
│   └── user.go             # Định nghĩa cấu trúc bảng User trong database
├── routes/
│   └── user_routes.go      # Định tuyến (Routing) các API của User
├── .env                    # File chứa các biến môi trường (Database, Port...)
├── go.mod                  # Quản lý thư viện (Dependencies)
└── main.go                 # Điểm bắt đầu của ứng dụng, khởi tạo server
```

## ⚙️ Hướng dẫn cài đặt và chạy dự án

### 1. Yêu cầu hệ thống
- Đã cài đặt Go (phiên bản 1.18 trở lên).
- Đã cài đặt PostgreSQL và có sẵn một database (ví dụ: `marketing_db`).

### 2. Cấu hình biến môi trường
Tạo file `.env` ở thư mục gốc của dự án (nếu chưa có) và điền thông tin kết nối Database:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=marketing_db
```

### 3. Cài đặt thư viện
Mở terminal tại thư mục dự án và chạy lệnh:
```bash
go mod tidy
```

### 4. Khởi chạy Server (Hỗ trợ Hot-reload)
Dự án sử dụng **Air** để tự động reload server mỗi khi bạn lưu file (giống nodemon bên Node.js).

**Cài đặt Air (chỉ cần chạy 1 lần):**
```bash
go install github.com/air-verse/air@latest
```

**Chạy server:**
```bash
air
```

Server sẽ chạy ở địa chỉ: `http://localhost:8080`

---

## 🛠 Các chức năng đã hoàn thành (Tiến độ hiện tại)

### 1. Quản lý Người dùng (User)
- **Đăng ký tài khoản mới**
  - **Endpoint:** `POST /api/users/register`
  - **Mô tả:** Nhận thông tin `name`, `email`, `password`. Kiểm tra email trùng lặp, mã hóa mật khẩu bằng `bcrypt` và lưu vào database.
- **Đăng nhập (Login)**
  - **Endpoint:** `POST /api/users/login`
  - **Mô tả:** Nhận thông tin `email`, `password`. Kiểm tra thông tin đăng nhập, nếu đúng sẽ trả về một **JWT Token** để sử dụng cho các API cần xác thực.
- **Lấy danh sách người dùng**
  - **Endpoint:** `GET /api/users/listusers`
  - **Mô tả:** Lấy toàn bộ danh sách người dùng đã đăng ký từ database. Dữ liệu trả về đã được lọc bỏ trường `password` để đảm bảo bảo mật.

---
*README này sẽ được cập nhật liên tục theo tiến độ phát triển của dự án.*
