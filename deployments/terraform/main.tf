variable "bottoken" {
  description = "Bot token provided by Discord."
  type        = string
  sensitive   = true
}
variable "namespace" {
  description = "Kubernetes Namespace."
  type        = string
  sensitive   = false
  default     = "discordbot"
}

variable "configpath" {
  default     = "~/.kube/prod.config"
  description = "Path to Kubeconfig."
  type       = string
}

provider "kubernetes" {
  config_path    = var.configpath
  config_context = "default"
}

resource "kubernetes_namespace" "discordbot" {
  metadata {
    name = "discordbot"
  }
}

resource "kubernetes_deployment" "discordbot" {
  metadata {
    name      = "discordbot"
    namespace = var.namespace
    labels = {
      purpose = "discordbot"
    }
  }

  spec {
    replicas = "1"
    selector {
      match_labels = {
        purpose = "discordbot"
      }
    }
    template {
      metadata {
        name      = "discordbot"
        namespace = var.namespace
        labels = {
          purpose = "discordbot"
        }
      }
      spec {

        container {
          image = "aymon/hive.discordbot.slim:latest"
          name  = "discordbot"
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
        }
      }
    }
  }
}
