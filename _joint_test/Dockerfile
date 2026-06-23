FROM golang:1.24.0-alpine

RUN go install github.com/air-verse/air@latest

WORKDIR /src
COPY go.mod go.sum ./
COPY _joint_test/go.mod _joint_test/go.sum ./_joint_test/
RUN cd _joint_test && go mod download

COPY . .

WORKDIR /src/_joint_test
CMD ["air", "-c", ".air.toml"]
