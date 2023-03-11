ARG  DISTROLESS_IMAGE=gcr.io/distroless/static:nonroot

FROM golang:alpine3.17 as build

COPY . /app
WORKDIR /app

ENV CGO_ENABLED=0
RUN go build -o /go/bin/user-manager

FROM ${DISTROLESS_IMAGE}

USER 65532:65532

# Copy the binary from the previous stage
COPY --from=build /go/bin/user-manager /go/bin/user-manger 
COPY --from=build ssl.so.3     /lib/libssl.so.3
COPY --from=build crypto.so.3     /lib/libcrypto.so.3
COPY --from=build musl-x86_64.so.1  /lib/ld-musl-x86_64.so.1



ENV PostgresHost=local_pgdb
ENV PostgresUser=postgres
ENV PostgresPassword=admin
ENV PostgresDatabase=postgres
ENV PostgresPort=5432
ENV ChatsAPI=http://127.0.0.1:3002
ENV Origins=http://localhost:8080
ENV FilesApi=http://127.0.0.1:3001
ENV RPDisplayName=LocalTest
ENV RPID=localhost
ENV RPOrigin=http://localhost:80
ENV RPIcon=https://duo.com/logo.png
ENV AppListen=":80"

# Expose the port that the application will listen on
EXPOSE 80

# Run the binary
ENTRYPOINT ["/go/bin/chatsapi"]
