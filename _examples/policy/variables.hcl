variable "enviroments" {
  description = "enviroment variables"
  type        = "list"
  default     = ["prod", "dev"]
}

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
  type        = "map"

  default = {
    # "x-gateway-jp-prod" = "gateway"
    # "x-gateway-jp-dev"  = "gateway"
    "gateway" = "x-gateway-jp-dev"
  }
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
