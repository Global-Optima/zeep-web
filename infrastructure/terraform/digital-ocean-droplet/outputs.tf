###############################################################################
# outputs.tf
###############################################################################

output "load_balancer_ip" {
  description = "Public IP of the Load Balancer"
  value       = digitalocean_loadbalancer.lb.ip
}

output "app_server_ips" {
  description = "Public IP addresses of application servers"
  value = [
    for droplet in digitalocean_droplet.app_server : droplet.ipv4_address
  ]
}

output "postgres_connection_info" {
  description = "Connection details for the Postgres database"
  value = {
    host     = digitalocean_database_cluster.postgres.host
    port     = digitalocean_database_cluster.postgres.port
    database = digitalocean_database_cluster.postgres.database
    user     = digitalocean_database_cluster.postgres.user
    password = digitalocean_database_cluster.postgres.password
  }
  sensitive = true
}

output "redis_connection_info" {
  description = "Connection details for the Redis instance"
  value = {
    host     = digitalocean_database_cluster.redis.host
    port     = digitalocean_database_cluster.redis.port
    password = digitalocean_database_cluster.redis.password
  }
  sensitive = true
}

output "spaces_bucket_name" {
  description = "The name of the S3-compatible Spaces bucket"
  value       = digitalocean_spaces_bucket.bucket.name
}

output "bucket_cdn_endpoint" {
  description = "The domain name for the CDN distribution"
  value       = digitalocean_cdn.bucket_cdn.endpoint
}