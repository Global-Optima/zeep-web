###############################################################################
# main.tf (Kubernetes-Based Deployment)
###############################################################################

terraform {
  required_version = ">= 1.11.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.49.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.24.0"
    }
  }
}

##############################
# Providers Configuration
##############################

# DigitalOcean Provider
provider "digitalocean" {
  token             = var.do_token
  spaces_access_id  = var.spaces_access_id
  spaces_secret_key = var.spaces_access_token
}

# Create a VPC
resource "digitalocean_vpc" "main" {
  name     = var.vpc_name
  region   = var.region
  ip_range = "10.10.0.0/16"
}

# Create a Kubernetes Cluster in the VPC
resource "digitalocean_kubernetes_cluster" "main" {
  name     = "zeep-cluster"
  region   = var.region
  version  = "1.30.3-do.0"  # Use a supported DO Kubernetes version
  vpc_uuid = digitalocean_vpc.main.id

  node_pool {
    name       = "default-pool"
    size       = var.k8s_node_size
    node_count = var.k8s_node_count
    tags       = concat(var.tags, ["k8s-node"])
  }
}

# Configure the Kubernetes Provider using the cluster's kube_config
provider "kubernetes" {
  host = digitalocean_kubernetes_cluster.main.endpoint

  token = var.do_token

  cluster_ca_certificate = base64decode(
    digitalocean_kubernetes_cluster.main.kube_config[0].cluster_ca_certificate
  )
}

# Create a Firewall for Kubernetes
resource "digitalocean_firewall" "k8s_firewall" {
  name = "k8s-firewall"

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = [var.ssh_allowed_cidr]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80,443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "icmp"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  # Allow all TCP within the VPC
  inbound_rule {
    protocol         = "tcp"
    port_range       = "1-65535"
    source_addresses = [digitalocean_vpc.main.ip_range]
  }

  # Allow Kubernetes API traffic
  inbound_rule {
    protocol         = "tcp"
    port_range       = "6443"
    source_addresses = [digitalocean_vpc.main.ip_range]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "all"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

##############################
# Managed Databases
##############################

# PostgreSQL Database Cluster
resource "digitalocean_database_cluster" "postgres" {
  name                 = "app-postgres-cluster"
  engine               = "pg"
  version              = "17"
  size                 = var.db_node_size
  region               = var.region
  node_count           = 2
  private_network_uuid = digitalocean_vpc.main.id
  tags                 = concat(var.tags, ["postgres-db"])
}

# Redis Database Cluster
resource "digitalocean_database_cluster" "redis" {
  name                 = "app-redis-cache"
  engine               = "redis"
  version              = "7"
  size                 = var.redis_node_size
  region               = var.region
  node_count           = 1
  private_network_uuid = digitalocean_vpc.main.id
  tags                 = concat(var.tags, ["redis-cache"])
}

##############################
# Spaces Bucket & CDN
##############################

resource "digitalocean_spaces_bucket" "bucket" {
  name   = var.bucket_name
  region = var.region
}

resource "digitalocean_cdn" "bucket_cdn" {
  origin = digitalocean_spaces_bucket.bucket.bucket_domain_name
}

##############################
# Domain and Certificate (Conditional)
##############################

resource "digitalocean_domain" "zeep_domain" {
  count = var.use_domain ? 1 : 0
  name  = var.domain_name
}

resource "digitalocean_certificate" "zeep_certificate" {
  count   = var.use_domain ? 1 : 0
  name    = var.certificate_name
  type    = "lets_encrypt"
  domains = [var.domain_name]
}

##############################
# Kubernetes Secret (app-secrets)
##############################

resource "kubernetes_secret" "app_secrets" {
  metadata {
    name      = "app-secrets"
    namespace = "default"
  }
  data = {
    # Application Configuration
    ENV               = base64encode("production")
    SERVER_PORT       = base64encode("8080")
    # If using a domain, CLIENT_URL is built from it; otherwise, you can adjust.
    CLIENT_URL        = base64encode(var.use_domain ? "https://${var.domain_name}" : "http://<LB_IP>") 
    DEFAULT_PAGE      = base64encode("1")
    DEFAULT_PAGE_SIZE = base64encode("10")
    MAX_PAGE_SIZE     = base64encode("100")

    # PostgreSQL values from managed resource outputs
    DB_HOST           = base64encode(digitalocean_database_cluster.postgres.host)
    DB_PORT           = base64encode(tostring(digitalocean_database_cluster.postgres.port))
    DB_USER           = base64encode(digitalocean_database_cluster.postgres.user)
    DB_PASSWORD       = base64encode(digitalocean_database_cluster.postgres.password)
    DB_NAME           = base64encode(digitalocean_database_cluster.postgres.database)

    # Redis values from managed resource outputs
    REDIS_HOST        = base64encode(digitalocean_database_cluster.redis.host)
    REDIS_PORT        = base64encode(tostring(digitalocean_database_cluster.redis.port))
    REDIS_PASSWORD    = base64encode(digitalocean_database_cluster.redis.password)

    # Spaces configuration
    S3_ACCESS_KEY     = base64encode(var.spaces_access_id)
    S3_SECRET_KEY     = base64encode(var.spaces_access_token)
    S3_ENDPOINT       = base64encode("https://${var.region}.digitaloceanspaces.com")
    S3_BUCKET_NAME    = base64encode(digitalocean_spaces_bucket.bucket.name)

    # Frontend configuration
    NGINX_SERVER_NAME = base64encode(var.domain_name)
    VITE_API_URL      = base64encode(var.use_domain ? "https://${var.domain_name}/api" : "http://<LB_IP>/api")
    VITE_WS_URL       = base64encode(var.use_domain ? "wss://${var.domain_name}/api" : "ws://<LB_IP>/api")
    VITE_PAYMENT_SECRET = base64encode(var.payment_secret)

    # JWT secrets
    JWT_CUSTOMER_SECRET_KEY = base64encode(var.jwt_customer_secret_key)
    JWT_EMPLOYEE_SECRET_KEY = base64encode(var.jwt_employee_secret_key)
  }
  depends_on = [
    digitalocean_database_cluster.postgres,
    digitalocean_database_cluster.redis,
    digitalocean_spaces_bucket.bucket
  ]
}
