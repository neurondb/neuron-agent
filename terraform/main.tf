# NeuronAgent Terraform module (example structure).
# Adapt to your cloud provider (e.g. AWS EKS, GKE, or VM + Docker).
# This file documents expected resources; actual provider resources
# depend on your target (Kubernetes, ECS, Cloud Run, etc.).

terraform {
  required_version = ">= 1.0"
  required_providers {
    null = { source = "hashicorp/null", version = "~> 3.0" }
  }
  # Uncomment and set backend for state
  # backend "s3" { bucket = "my-tf-state"; key = "neurondb/agent.tfstate"; region = "us-east-1" }
}

# Example: local-exec to run Helm (when Kubernetes provider is used)
# resource "helm_release" "neuron_agent" {
#   name       = "neuron-agent"
#   chart      = "../helm"
#   namespace  = "default"
#   set_sensitive {
#     name  = "database.password"
#     value = var.db_password
#   }
#   set {
#     name  = "database.host"
#     value = var.db_host
#   }
#   set {
#     name  = "database.name"
#     value = var.db_name
#   }
#   set {
#     name  = "database.user"
#     value = var.db_user
#   }
#   set {
#     name  = "image.repository"
#     value = split(":", var.agent_image)[0]
#   }
#   set {
#     name  = "image.tag"
#     value = try(split(":", var.agent_image)[1], "latest")
#   }
#   set {
#     name  = "replicaCount"
#     value = var.agent_replicas
#   }
# }

# Placeholder so terraform validate succeeds; replace with real resources for your platform.
resource "null_resource" "neuron_agent_placeholder" {
  triggers = {
    env   = var.environment
    image = var.agent_image
    db    = var.db_host
  }
}
