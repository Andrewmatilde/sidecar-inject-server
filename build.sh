go mod download
go mod vendor
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sidecar-inject-server .
docker build --no-cache -t yisaer/sidecar-inject-server:v1 .
rm -rf sidecar-inject-server