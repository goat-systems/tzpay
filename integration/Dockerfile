FROM golang:1.12
RUN go get -u honnef.co/go/tools/cmd/staticcheck
RUN go build -o staticcheck $GOPATH/src/honnef.co/go/tools/cmd/staticcheck/staticcheck.go
RUN mv staticcheck /usr/local/bin/