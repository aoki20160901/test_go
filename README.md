
###
cd /Users/user/go_test2

### 必要な時がある
go mod init something
go mod tidy
go run ./cmd/server

### postgreとMSSQLはdocker
docker compose up