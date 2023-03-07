FROM golang as build

COPY . /app
WORKDIR /app

ENV CGO_ENABLED=0
RUN go build -o userApp

FROM alpine:latest as product

WORKDIR /app
COPY --from=build /app/userApp /bin/userApp

RUN chmod 777 /bin/userApp

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

EXPOSE 80

CMD ["userApp"]
