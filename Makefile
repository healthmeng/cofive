export GOPATH=$(PWD)

all:cofived client

cofived:src/cofived/main.go
	go build -gcflags "-N"  cofived
# #go build -gcflags "-N" client
client:
	go install client
clean:
	rm -f cofived
#
#
