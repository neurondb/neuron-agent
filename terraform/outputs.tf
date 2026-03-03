# Outputs for NeuronAgent deployment.
# When using Helm provider, expose service URL and DB connection info (redacted).

output "environment" {
  value     = var.environment
  description = "Environment name"
}

output "service_url" {
  value       = "http://neuron-agent:8080"
  description = "NeuronAgent service URL (cluster-internal); use ingress or LB for external"
}

output "db_connection_info" {
  value = {
    host = var.db_host
    port = var.db_port
    name = var.db_name
    user = var.db_user
    # password intentionally omitted
  }
  description = "Database connection info (password from secret)"
  sensitive   = false
}
