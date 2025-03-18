####################################################################
# variables.tf
####################################################################

variable "do_token" {
  type        = string
  description = "DigitalOcean API Token"
  sensitive   = true
}

variable "spaces_access_id" {
  type        = string
  description = "Spaces Access Key ID"
  sensitive   = true
}

variable "spaces_access_token" {
  type        = string
  description = "Spaces Secret Access Token"
  sensitive   = true
}

variable "region" {
  type        = string
  default     = "fra1"
  description = "Region to deploy all resources in."
}

variable "ssh_key_ids" {
  type        = list(string)
  default     = []
  description = "List of DigitalOcean SSH key IDs for Droplet access."
}

variable "droplet_size" {
  type        = string
  default     = "s-2vcpu-4gb"
  description = "Droplet size slug for the application servers."
}

variable "droplet_image" {
  type        = string
  default     = "ubuntu-24-10-x64"
  description = "Droplet image for the application servers."
}

variable "db_node_size" {
  type        = string
  default     = "db-s-2vcpu-4gb"
  description = "Node size for the PostgreSQL DB cluster."
}

variable "redis_node_size" {
  type        = string
  default     = "db-s-1vcpu-1gb"
  description = "Node size for the Redis cluster."
}

variable "vpc_name" {
  type        = string
  default     = "zeep-vps"
  description = "Name of the VPC."
}

variable "tags" {
  type        = list(string)
  default     = ["terraform-managed", "production"]
  description = "List of tags to apply to resources."
}

variable "bucket_name" {
  type        = string
  default     = "zeep-media"
  description = "Name of the S3-compatible Spaces bucket."
}

variable "app_server_count" {
  type        = number
  default     = 2
  description = "Number of Droplets (app servers) to create."
}

variable "ssh_allowed_cidr" {
  type        = string
  default     = "0.0.0.0/0"
  description = "CIDR block allowed to SSH into the Droplets."
}

variable "domains" {
  type        = list(string)
  default     = ["zeep.kz", "www.zeep.kz"]
  sensitive   = true
  description = "The domains for the load balancer (e.g. example.com)."
}

variable "certificate_name" {
  type        = string
  default = "zeep-cert"
  description = "The name of the DigitalOcean certificate"
}

variable "ssh_private_key_path" {
  type        = string
  description = "Path to the SSH private key file for accessing Droplets"
  sensitive   = true
}

variable "github_token" {
  type        = string
  description = "GitHub Personal Access Token for private repository access"
  sensitive   = true
}

variable "github_repo" {
  type        = string
  description = "GitHub repository in format username/repository"
  default     = "your-org/zeep-app"
}

variable "github_branch" {
  type        = string
  description = "GitHub branch to deploy"
  default     = "main"
}

variable "jwt_customer_secret" {
  type        = string
  description = "Secret key for customer JWT tokens"
  sensitive   = true
}

variable "jwt_employee_secret" {
  type        = string
  description = "Secret key for employee JWT tokens"
  sensitive   = true
}

variable "payment_secret" {
  type        = string
  description = "Secret key for payment processing"
  sensitive   = true
}