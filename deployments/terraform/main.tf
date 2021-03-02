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
        }
      }
    }
  }
}
