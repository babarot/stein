variable "shortened_environment" {
  description = "Shortened environment, such as prod, dev"
  type        = "map"

  default = {
    production  = "prod"
    development = "dev"
    laboratory  = "lab"
  }
}

variable "special_cases" {
  description = ""
  type        = "list"

  default = [
    "x-gateway-jp-dev",
    "x-gateway-jp-prod",
  ]
}

variable "namespace_name_map" {
  type = "map"

  default = {
    "gateway" = [
      "x-gateway-jp-dev",
      "x-gateway-jp-prod",
    ]
  }
}
