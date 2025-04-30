terraform {
  required_providers {
    hashicups = {
      source  = "hashicorp/hashicups"
      version = "0.0.0"
    }
  }
}

module "example" {
  source = "./source"
}
