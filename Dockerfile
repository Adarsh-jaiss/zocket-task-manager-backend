FROM golang:buster as build 

WORKDIR  /app

RUN useradd -u 1001 nonroot

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
 -ldflags="-linkmode external -extldflags -static" \
 -tags netgo \ 
 -o main 

# Multistage build
FROM scratch

# Copy the .env file into the image
COPY .env .env

# ENV GIN_MODE=release
ENV GOFIBER_MODE=release

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/main main

USER nonroot

EXPOSE 8000

CMD ["./main"]