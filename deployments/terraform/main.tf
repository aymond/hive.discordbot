data local_file config {
  filename = "../../configs/config.yaml"
}
resource "kubernetes_namespace" "discordbot" {
  metadata {
    name = var.name
  }
}

resource "kubernetes_deployment" "discordbot" {
  metadata {
    name      = var.name
    namespace = kubernetes_namespace.discordbot.id
    labels = {
      purpose = var.name
    }
  }

  spec {
    replicas = "1"
    selector {
      match_labels = {
        purpose = var.name
      }
    }
    template {
      metadata {
        name      = var.name
        namespace = kubernetes_namespace.discordbot.id
        labels = {
          purpose = var.name
          app     = var.name
        }
      }
      spec {

        container {
          image = "${var.containerimage}:${var.containerimageversion}"
          name  = var.name
          resources {
            limits = {
              cpu    = "100m"
              memory = "20Mi"
            }
            requests = {
              cpu    = "10m"
              memory = "8Mi"
            }
          }
          env {
            name  = "TOKEN"
            value = var.bottoken
          }

          volume_mount {
            mount_path        = "/data/discordbot-config"
            mount_propagation = "None"
            name              = "discordbot-config"
            read_only         = true
          }
        }

        volume {
          name = "discordbot-config"
          config_map {
            name         = "discordbot-config"
            default_mode = "0644"
            optional     = true
          }
        }
      }
    }
  }
}

resource "kubernetes_config_map" "discordbot" {
  metadata {
    name      = "discordbot-config"
    namespace = kubernetes_namespace.discordbot.id
    labels = {
      purpose = var.name
    }
  }

  data = {
    "commands.yaml" = data.local_file.config.content
    "version.yaml"  = <<EOF
version: ${var.containerimageversion}
\n
EOF
  }
}