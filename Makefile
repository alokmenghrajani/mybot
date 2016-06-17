SOURCE_FILES := $(shell find . \( -name '*.go' -not -path './vendor*' \))

all:	mybot

mybot:	vendor $(SOURCE_FILES)
	go build -o mybot .

vendor:
	glide install

mybot_linux:	vendor $(SOURC_FILES)
	GOOS=linux GOARCH=amd64 go build -o mybot_linux .
	scp mybot_linux ctf.quaxio.com:~/
