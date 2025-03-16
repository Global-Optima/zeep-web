###############################################################################
# main.tf
###############################################################################

terraform {
  required_version = ">= 1.11.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.49.0"
    }
  }
}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = var.do_token
  spaces_access_id = var.spaces_access_id
  spaces_secret_key = var.spaces_access_token
}

###############################################################################
# VPC
###############################################################################

resource "digitalocean_vpc" "main" {
  name      = var.vpc_name
  region    = var.region
  ip_range  = "10.10.0.0/16"
  description = "VPC for production environment"
}

###############################################################################
# Droplets (Application Servers) using count
###############################################################################

resource "digitalocean_droplet" "app_server" {
  count             = var.app_server_count
  name              = "app-server-${count.index + 1}"
  region            = var.region
  size              = var.droplet_size
  image             = var.droplet_image
  vpc_uuid          = digitalocean_vpc.main.id
  ssh_keys          = var.ssh_key_ids
  backups           = false
  monitoring        = true
  tags              = concat(var.tags, ["app-server"])
}

###############################################################################
# Firewall
###############################################################################
# This firewall allows:
#  - SSH from a trusted IP (configure via var.ssh_allowed_cidr).
#  - HTTP/HTTPS from anywhere.
#  - ICMP from anywhere.
#  - Full TCP traffic within the VPC.
#  - Outbound to anywhere for any TCP.

resource "digitalocean_firewall" "app_firewall" {
  name        = "app-firewall"
  droplet_ids = [
    for droplet in digitalocean_droplet.app_server : droplet.id
  ]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = [var.ssh_allowed_cidr]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "icmp"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "1-65535"
    source_addresses = [digitalocean_vpc.main.ip_range]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "all"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  lifecycle {
    ignore_changes = [
      droplet_ids
    ]
  }
}

###############################################################################
# Load Balancer
###############################################################################

resource "digitalocean_loadbalancer" "lb" {
  name     = "app-lb"
  region   = var.region
  vpc_uuid = digitalocean_vpc.main.id

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"
    target_port     = 80
    target_protocol = "http"
  }

  forwarding_rule {
    entry_port       = 443
    entry_protocol   = "https"
    target_port      = 80
    target_protocol  = "http"
    tls_passthrough  = false
  }

  healthcheck {
    protocol                 = "http"
    port                     = 80
    path                     = "/healthz"
    check_interval_seconds   = 10
    response_timeout_seconds = 5
    unhealthy_threshold      = 3
    healthy_threshold        = 5
  }

  droplet_tag = "app-server"

  depends_on = [
    digitalocean_droplet.app_server
  ]
}

###############################################################################
# Managed Databases
###############################################################################

resource "digitalocean_database_cluster" "postgres" {
  name                = "app-postgres-cluster"
  engine              = "pg"
  version             = "17"
  size                = var.db_node_size
  region              = var.region
  node_count          = 2
  private_network_uuid = digitalocean_vpc.main.id
  tags               = concat(var.tags, ["postgres-db"])
}

resource "digitalocean_database_cluster" "redis" {
  name                = "app-redis-cache"
  engine              = "redis"
  version             = "7"
  size                = var.redis_node_size
  region              = var.region
  node_count          = 1
  private_network_uuid = digitalocean_vpc.main.id
  tags               = concat(var.tags, ["redis-cache"])
}

###############################################################################
# Spaces (S3-Compatible)
###############################################################################

resource "digitalocean_spaces_bucket" "bucket" {
  name   = var.bucket_name
  region = var.region
}

resource "digitalocean_cdn" "bucket_cdn" {
  origin = digitalocean_spaces_bucket.bucket.bucket_domain_name
}