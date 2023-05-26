# syntax=docker/dockerfile:1

FROM golang:1.20.4-alpine

WORKDIR /app

# RUN go install github.com/cosmtrek/air@latest

COPY . ./

RUN go mod tidy