HASH=$(git rev-parse --short HEAD)
BUILD="${HASH}-$(date '+%Y%m%d')"
VERSION=0.0.1
echo "Version: $VERSION Build: $BUILD"
env CGOENABLED=0 GOOS=linux GOARCH=amd64 go build -o daemon_linux -v -ldflags "-X main.version=$VERSION -X main.build=$BUILD" ../cmd/teü/main.go
gzip daemon_linux