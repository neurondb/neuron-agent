variable "environment" {
  type        = string
  default     = "prod"
  description = "Environment name (e.g. prod, staging)"
}

variable "db_host" {
  type        = string
  description = "PostgreSQL host"
}

variable "db_port" {
  type        = number
  default     = 5432
  description = "PostgreSQL port"
}

variable "db_name" {
  type        = string
  description = "PostgreSQL database name"
}

variable "db_user" {
  type        = string
  description = "PostgreSQL user"
}

variable "db_password" {
  type        = string
  sensitive   = true
  description = "PostgreSQL password"
}

variable "agent_image" {
  type        = string
  default     = "neurondb/neuron-agent:latest"
  description = "NeuronAgent container image"
}

variable "agent_replicas" {
  type        = number
  default     = 1
  description = "Number of agent replicas"
}

variable "enable_lb" {
  type        = bool
  default     = true
  description = "Create load balancer for the service"
}
