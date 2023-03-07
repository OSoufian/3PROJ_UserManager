terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "1.12.0"
    }
  }
}

provider "kubernetes" {
  config_context_cluster = "minikube"
  config_context_user    = "minikube"
}

resource "kubernetes_service" "user_api" {
  metadata {
    name = "user-api"
    labels = {
      app = "user-api"
    }
  }

  spec {
    selector = {
      app = "user-api"
    }

    port {
      port        = 80
      target_port = 80
    }
  }
}

resource "kubernetes_deployment" "user_api" {
  metadata {
    name = "user-api"
    labels = {
      app = "user-api"
    }
  }

  spec {
    selector {
      match_labels = {
        app = "user-api"
      }
    }

    template {
      metadata {
        labels = {
          app = "user-api"
        }
      }

      spec {
        container {
          image = "user-api"
          name  = "user-api"
          port {
            container_port = 80
          }
        }
      }
    }
  }
}

resource "null_resource" "wait_for_user_api" {
  provisioner "local-exec" {
    command = "sleep 10"
  }

  depends_on = [
    kubernetes_deployment.user_api,
  ]
}

resource "kubernetes_service" "upload_api" {
  metadata {
    name = "upload-api"
    labels = {
      app = "upload-api"
    }
  }

  spec {
    selector = {
      app = "upload-api"
    }

    port {
      port        = 80
      target_port = 80
    }
  }
}

resource "kubernetes_deployment" "upload_api" {
  metadata {
    name = "upload-api"
    labels = {
      app = "upload-api"
    }
  }

  spec {
    selector {
      match_labels = {
        app = "upload-api"
      }
    }

    template {
      metadata {
        labels = {
          app = "upload-api"
        }
      }

      spec {
        container {
          image = "upload-api"
          name  = "upload-api"
          port {
            container_port = 80
          }
        }
      }
    }
  }
}

resource "null_resource" "wait_for_upload_api" {
  provisioner "local-exec" {
    command = "sleep 10"
  }

  depends_on = [
    kubernetes_deployment.upload_api,
  ]
}

resource "kubernetes_service" "chats_api" {
  metadata {
    name = "chats-api"
    labels = {
      app = "chats-api"
    }
  }

  spec {
    selector = {
      app = "chats-api"
    }

    port {
      port        = 80
      target_port = 80
    }
  }
}

resource "kubernetes_deployment" "chats_api" {
  metadata {
    name = "chats-api"
    labels = {
      app = "chats-api"
    }
  }

  spec {
    selector {
      match_labels = {
        app = "chats-api"
      }
    }

    template {
      metadata {
        labels = {
          app = "chats-api"
        }
      }

      spec {
        container {
          image = "chats-api"
          name  = "chats-api"
          port {
            container_port = 80
          }
        }
      }
    }
  }
}

resource "null_resource" "wait_for_chats_api" {
  provisioner "local-exec" {
    command = "sleep 10"
  }

  depends_on = [
    kubernetes_deployment.chats_api,
  ]
}

resource "kubernetes_service" "local_pgdb" {
  metadata {
    name = "local-pgdb"
    labels = {
      app = "local-pgdb"
    }
  }

  spec {
    selector = {
      app = "local-pgdb"
    }

    port {
      port        = 5432
      target_port = 5432
    }
  }
}

resource "kubernetes_deployment" "local_pgdb" {
  metadata {
    name = "local-pgdb"
    labels = {
      app = "local-pgdb"
    }
  }

  spec {
    selector {
      match_labels = {
        app = "local-pgdb"
      }
    }

    template {
      metadata {
        labels = {
          app = "local-pgdb"
        }
      }

      spec {
        container {
          image = "postgres"
          name  = "local-pgdb"
          port {
            container_port = 5432
          }
          env {
            name  = "POSTGRES_USER"
            value = "postgres"
          }
          env {
            name  = "POSTGRES_PASSWORD"
            value = "admin"
          }
        }
      }
    }
  }
}

resource "null_resource" "wait_for_local_pgdb" {
  provisioner "local-exec" {
    command = "sleep 10"
  }

  depends_on = [
    kubernetes_deployment.local_pgdb,
  ]
}

resource "kubernetes_service" "pgadmin" {
  metadata {
    name = "pgadmin"
    labels = {
      app = "pgadmin"
    }
  }

  spec {
    selector = {
      app = "pgadmin"
    }

    port {
      port        = 80
      target_port = 80
    }
  }
}

resource "kubernetes_deployment" "pgadmin" {
  metadata {
    name = "pgadmin"
    labels = {
      app = "pgadmin"
    }
  }

  spec {
    selector {
      match_labels = {
        app = "pgadmin"
      }
    }

    template {
      metadata {
        labels = {
          app = "pgadmin"
        }
      }

      spec {
        container {
          image = "dpage/pgadmin4"
          name  = "pgadmin"
          port {
            container_port = 80
          }
          env {
            name  = "PGADMIN_DEFAULT_EMAIL"
            value = "contact@cyprientaib.com"
          }
          env {
            name  = "PGADMIN_DEFAULT_PASSWORD"
            value = "admin"
          }
        }
      }
    }
  }
}

resource "null_resource" "wait_for_pgadmin" {
  provisioner "local-exec" {
    command = "sleep 10"
  }

  depends_on = [
    kubernetes_deployment.pgadmin,
  ]
}