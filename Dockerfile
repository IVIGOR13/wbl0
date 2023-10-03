# BUILD STEP
FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

# FINAL STEP
FROM alpine

ENV PORT 8080
EXPOSE $PORT

RUN apk update
RUN apk add postgresql-client

WORKDIR /app

COPY --from=build /app/l0 /app/

COPY ./wait-for-postgres.sh .
RUN chmod +x ./wait-for-postgres.sh

CMD ["./l0"]
