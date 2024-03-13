# FROM golang:1.22

# COPY ./ /app
# WORKDIR /app
# RUN make build

# ENTRYPOINT  ["./dist/server"]
FROM gcr.io/distroless/static-debian12

COPY ./dist/server /app/server

ENTRYPOINT ["./app/server"]
