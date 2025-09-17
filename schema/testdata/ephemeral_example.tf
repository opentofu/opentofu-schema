# Example OpenTofu configuration demonstrating ephemeral resources in v1.11

# Ephemeral variable for sensitive data
variable "api_key" {
  type        = string
  description = "API key for external service"
  ephemeral   = true
  sensitive   = true
}

# Ephemeral resource to fetch temporary credentials
ephemeral "aws_secretsmanager_secret_version" "temp_credentials" {
  secret_id = "arn:aws:secretsmanager:us-west-2:123456789012:secret:temp-creds"

  lifecycle {
    precondition {
      condition     = var.api_key != ""
      error_message = "API key must be provided"
    }

    postcondition {
      condition     = ephemeral.aws_secretsmanager_secret_version.temp_credentials.secret_string != ""
      error_message = "Secret must not be empty"
    }
  }

  depends_on = [aws_secretsmanager_secret.example]
}

# Provider configuration using ephemeral values
provider "aws" {
  alias      = "temp"
  access_key = jsondecode(ephemeral.aws_secretsmanager_secret_version.temp_credentials.secret_string)["access_key"]
  secret_key = jsondecode(ephemeral.aws_secretsmanager_secret_version.temp_credentials.secret_string)["secret_key"]
}

# Resource using write-only attributes with ephemeral data
resource "aws_ssm_parameter" "secret_param" {
  provider = aws.temp
  name     = "/myapp/secret"
  type     = "SecureString"
  value_wo = jsonencode({
    api_key    = var.api_key
    temp_token = jsondecode(ephemeral.aws_secretsmanager_secret_version.temp_credentials.secret_string)["token"]
  })
  value_wo_version = 1
}

# Ephemeral output for module composition
output "temp_endpoint" {
  value       = jsondecode(ephemeral.aws_secretsmanager_secret_version.temp_credentials.secret_string)["endpoint"]
  description = "Temporary endpoint URL"
  ephemeral   = true
}

# Local value using ephemeral data
locals {
  # This local becomes ephemeral because it references ephemeral values
  connection_string = "${jsondecode(ephemeral.aws_secretsmanager_secret_version.temp_credentials.secret_string)["endpoint"]}?key=${var.api_key}"
}

# Supporting resource (regular resource)
resource "aws_secretsmanager_secret" "example" {
  name = "temp-creds"
}

# Resource with provisioner using ephemeral values
resource "null_resource" "deploy" {
  # This provisioner output will be suppressed due to ephemeral values
  provisioner "local-exec" {
    command = "deploy.sh --endpoint=${local.connection_string}"
  }

  # This provisioner output will be shown normally
  provisioner "local-exec" {
    command = "echo 'Deployment started'"
  }

  connection {
    type = "ssh"
    host = "example.com"
    user = "deploy"
    # This uses ephemeral data, so connection details won't be logged
    private_key = jsondecode(ephemeral.aws_secretsmanager_secret_version.temp_credentials.secret_string)["ssh_key"]
  }
}
