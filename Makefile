
OUT_DIR = ./build

all: linux windows

linux:
	GOOS=linux GOARCH=amd64 go build -o "$(OUT_DIR)/sol" .

windows:
	GOOS=windows GOARCH=amd64 go build -o "$(OUT_DIR)/sol.exe" .
