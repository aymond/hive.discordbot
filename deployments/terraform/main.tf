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

variable "containerimage" {
  description = "Discordbot container image"
  type        = string
  sensitive   = false
  default     = "aymon/hive.discordbot:dev"
}

variable "configpath" {
  default     = "~/.kube/kubeconfig.prod"
  description = "Path to Kubeconfig."
  type        = string
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
          app     = "discordbot"
        }
      }
      spec {

        container {
          image = var.containerimage
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

resource "kubernetes_service" "webbot" {
  metadata {
    name = "discordbot-web"
  }
  spec {
    selector = {
      app = kubernetes_deployment.discordbot.spec.0.template.0.metadata[0].labels.app
    }
    
    port {
      //node_port   = 30201
      port        = 3000
      target_port = 3000
    }

    type = "LoadBalancer"
    //session_affinity = "ClientIP"
  }
}
