FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o syllabea ./

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/syllabea .
EXPOSE 9090
CMD ["./syllabea"]
