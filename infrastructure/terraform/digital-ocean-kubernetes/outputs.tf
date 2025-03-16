output "kubernetes_endpoint" {
  description = "Kubernetes cluster endpoint"
  value       = digitalocean_kubernetes_cluster.main.endpoint
}

output "kube_config" {
  description = "Kubernetes cluster kubeconfig"
  value       = digitalocean_kubernetes_cluster.main.kube_config
  sensitive   = true
}

output "postgres_connection" {
  description = "Postgres connection details"
  value = {
    host     = digitalocean_database_cluster.postgres.host
    port     = digitalocean_database_cluster.postgres.port
    database = digitalocean_database_cluster.postgres.database
    user     = digitalocean_database_cluster.postgres.user
    password = digitalocean_database_cluster.postgres.password
  }
  sensitive   = true
}

output "redis_connection" {
  description = "Redis connection details"
  value = {
    host     = digitalocean_database_cluster.redis.host
    port     = digitalocean_database_cluster.redis.port
    password = digitalocean_database_cluster.redis.password
  }
  sensitive   = true
}

output "spaces_bucket" {
  description = "DigitalOcean Spaces bucket name"
  value       = digitalocean_spaces_bucket.bucket.name
}
