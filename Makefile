BINARY_NAME=word-bo-piwo
TARGET_HOST=piotrek@xdr.com.pl
TARGET_DIR=~/
SSH_PORT=222

.PHONY: build deploy all

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME)

deploy: build
	scp -P $(SSH_PORT) $(BINARY_NAME) $(TARGET_HOST):$(TARGET_DIR)

all: build deploy
