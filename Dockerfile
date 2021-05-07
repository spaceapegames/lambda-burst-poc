############################
# STEP 1 build executable binary
############################
FROM golang:1.14.0 AS builder

WORKDIR /go/src

ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /lambda-burst cmd/server/main.go

############################
# STEP 2 build Lambda-compatible image
############################
FROM  amazon/aws-lambda-go
# Copy our static executable
COPY --from=builder /lambda-burst /var/task/lambda-burst
CMD ["lambda-burst"]
