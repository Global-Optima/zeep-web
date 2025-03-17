###############################################################################
# Terraform Provider Configuration
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

provider "digitalocean" {
  token             = var.do_token
  spaces_access_id  = var.spaces_access_id
  spaces_secret_key = var.spaces_access_token
}

###############################################################################
# VPC Configuration
###############################################################################

resource "digitalocean_vpc" "main" {
  name        = var.vpc_name
  region      = var.region
  ip_range    = "10.10.0.0/16"
  description = "VPC for production environment"
}

###############################################################################
# Droplet Configuration (Application Servers)
###############################################################################

resource "digitalocean_droplet" "app_server" {
  count      = var.app_server_count
  name       = "app-server-${count.index + 1}"
  region     = var.region
  size       = var.droplet_size
  image      = var.droplet_image
  vpc_uuid   = digitalocean_vpc.main.id
  ssh_keys   = var.ssh_key_ids
  backups    = false
  monitoring = true
  tags       = concat(var.tags, ["app-server"])

  # Install Docker automatically
  user_data = <<-EOF
    #!/bin/bash
    set -eux

    # Remove existing Docker packages
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
      sudo apt-get remove -y $$pkg || true
    done

    # Install dependencies
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl

    # Add Docker GPG key and repository
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    echo "deb [arch=$$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $$(. /etc/os-release && echo "$${UBUNTU_CODENAME:-$$VERSION_CODENAME}") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    # Install Docker
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Start Docker and add user to docker group
    sudo systemctl enable --now docker
    sudo usermod -aG docker $$(whoami)
  EOF
}

###############################################################################
# Firewall Configuration
###############################################################################

resource "digitalocean_firewall" "app_firewall" {
  name        = "app-firewall"
  droplet_ids = [ for droplet in digitalocean_droplet.app_server : droplet.id ]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = [var.ssh_allowed_cidr]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["10.10.0.0/16", "0.0.0.0/0", "::/0"]
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

  outbound_rule {
    protocol              = "tcp"
    port_range            = "all"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  lifecycle {
    ignore_changes = [ droplet_ids ]
  }
}

###############################################################################
# Domain, DNS, and SSL Configuration
###############################################################################

resource "digitalocean_domain" "zeep_domain" {
  name = "zeep.kz"
}

resource "digitalocean_record" "root" {
  domain = digitalocean_domain.zeep_domain.name
  type   = "A"
  name   = "@"
  value  = digitalocean_loadbalancer.lb.ip
  ttl    = 300
}

resource "digitalocean_record" "www" {
  domain = digitalocean_domain.zeep_domain.name
  type   = "A"
  name   = "www"
  value  = digitalocean_loadbalancer.lb.ip
  ttl    = 300
}

resource "digitalocean_certificate" "zeep_certificate" {
  name    = "zeep-ssl-cert"
  type    = "lets_encrypt"
  domains = ["zeep.kz", "www.zeep.kz"]
}

###############################################################################
# Load Balancer Configuration
###############################################################################

resource "digitalocean_loadbalancer" "lb" {
  name     = "app-lb"
  region   = var.region
  vpc_uuid = digitalocean_vpc.main.id

  # Redirect HTTP to HTTPS
  forwarding_rule {
    entry_port      = 80
    entry_protocol  = "http"
    target_port     = 80
    target_protocol = "http"
  }

  # Secure HTTPS forwarding
  forwarding_rule {
    entry_port      = 443
    entry_protocol  = "https"
    target_port     = 80
    target_protocol = "http"
    certificate_name = digitalocean_certificate.zeep_certificate.name
  }

  healthcheck {
    protocol                 = "http"
    port                     = 80
    path                     = "/"
    check_interval_seconds   = 10
    response_timeout_seconds = 5
    unhealthy_threshold      = 3
    healthy_threshold        = 5
  }

  droplet_tag = "app-server"

  depends_on = [
    digitalocean_droplet.app_server,
    digitalocean_certificate.zeep_certificate
  ]
}

###############################################################################
# Managed Database Configuration
###############################################################################

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

###############################################################################
# Spaces (S3-Compatible) + CDN
###############################################################################

resource "digitalocean_spaces_bucket" "bucket" {
  name   = var.bucket_name
  region = var.region
}

resource "digitalocean_cdn" "bucket_cdn" {
  origin = digitalocean_spaces_bucket.bucket.bucket_domain_name
}
