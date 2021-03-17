provider "kubernetes" {
  config_path    = var.configpath[terraform.workspace]
  config_context = "default"
}
