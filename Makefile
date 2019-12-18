export GOPATH=$(PWD)

all:cofived client

cofived:
	go install cofived
# #go build -gcflags "-N" client
client:
	go install client
#
#
