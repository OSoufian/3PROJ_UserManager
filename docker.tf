resource "docker_container" "local_pgdb" {
  name  = "local_pgdb"
  image = "postgres"
  restart = "always"
  ports {
    internal = 5432
    external = 5432
  }
  env = {
    POSTGRES_USER = "postgres"
    POSTGRES_PASSWORD = "admin"
  }
  volumes {
    volume_name = "local_pgdata"
    container_path = "/var/lib/postgresql/data"
  }
}

resource "docker_container" "pgadmin" {
  name  = "pgadmin4_container"
  image = "dpage/pgadmin4"
  restart = "always"
  ports {
    internal = 80
    external = 5050
  }
  env = {
    PGADMIN_DEFAULT_EMAIL = "contact@cyprientaib.com"
    PGADMIN_DEFAULT_PASSWORD = "admin"
  }
  volumes {
    volume_name = "pgadmin-data"
    container_path = "/var/lib/pgadmin"
  }
  links = ["local_pgdb"]
}

resource "docker_container" "user_api" {
  name  = "user_api"
  image = "user-api"
  env = {
    PostgresHost = docker_container.local_pgdb.ip_address
    PostgresUser = "postgres"
    PostgresPassword = "admin"
    PostgresDatabase = "postgres"
    PostgresPort = "5432"
    FilesApi="http://127.0.0.1:3001"
    RPDisplayName = "LocalTest"
    RPID = "localhost"
    RPOrigin = "http://localhost:80"
    RPIcon = "https://duo.com/logo.png"
    AppListen = ":80"
  }
  ports {
    internal = 80
    external = 3000
  }
  links = ["local_pgdb"]
}

resource "docker_volume" "local_pgdata" {}

resource "docker_volume" "pgadmin-data" {}
