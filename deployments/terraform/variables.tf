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

variable "name" {
  description = "Default label for all resources."
  type        = string
  sensitive   = false
  default     = "discordbot"
}

variable "containerimage" {
  description = "Discordbot container image"
  type        = string
  sensitive   = false
  default     = "aymon/hive.discordbot"
}

variable "containerimageversion" {
  description = "Discordbot container image version"
  type        = string
  sensitive   = false
  default     = "dev"
}

variable "configpath" {
  description = "Path to Kubeconfig."
  type        = map(string)
  default = {
    prod    = "~/.kube/kubeconfig.prod"
    dev     = "~/.kube/kubeconfig.dev"
    default = "~/.kube/kubeconfig.prod"
  }
}
