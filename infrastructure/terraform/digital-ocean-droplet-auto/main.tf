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
# Droplet Configuration with Provisioning
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

  # Wait for the droplet to be fully provisioned before attempting to connect
  provisioner "remote-exec" {
    inline = ["echo 'Droplet is ready'"]

    connection {
      type    = "ssh"
      user    = "root"
      host    = self.ipv4_address
      private_key = file(var.ssh_private_key_path)
    }
  }

  # Upload the setup script
  provisioner "file" {
    destination = "/root/setup.sh"
    content     = replace(file("${path.module}/setup.sh"), "\r\n", "\n")
    
    connection {
      type    = "ssh"
      user    = "root"
      host    = self.ipv4_address
      private_key = file(var.ssh_private_key_path)
    }
  }
  # Upload the environment config file
  provisioner "file" {
    content = templatefile("${path.module}/config.env.tpl", {
      # PostgreSQL configuration
      POSTGRES_HOST     = digitalocean_database_cluster.postgres.host
      POSTGRES_PORT     = digitalocean_database_cluster.postgres.port
      POSTGRES_USER     = digitalocean_database_cluster.postgres.user
      POSTGRES_PASSWORD = digitalocean_database_cluster.postgres.password
      POSTGRES_DB       = digitalocean_database_cluster.postgres.database
      
      # Redis configuration
      REDIS_USERNAME    = digitalocean_database_cluster.redis.user
      REDIS_HOST        = digitalocean_database_cluster.redis.host
      REDIS_PORT        = digitalocean_database_cluster.redis.port
      REDIS_PASSWORD    = digitalocean_database_cluster.redis.password
      
      # Domain configuration
      DOMAIN            = var.domains[0]
      
      # S3 configuration
      S3_ACCESS_KEY     = var.spaces_access_id
      S3_SECRET_KEY     = var.spaces_access_token
      S3_ENDPOINT       = digitalocean_spaces_bucket.bucket.endpoint
      S3_BUCKET_NAME    = digitalocean_spaces_bucket.bucket.name
      
      # JWT and other application secrets
      JWT_CUSTOMER_SECRET_KEY = var.jwt_customer_secret
      JWT_EMPLOYEE_SECRET_KEY = var.jwt_employee_secret
      PAYMENT_SECRET          = var.payment_secret
      
      # GitHub configuration
      GITHUB_TOKEN      = var.github_token
      GITHUB_REPO       = var.github_repo
      GITHUB_BRANCH     = var.github_branch
    })
    destination = "/root/config.env"

    connection {
      type    = "ssh"
      user    = "root"
      host    = self.ipv4_address
      private_key = file(var.ssh_private_key_path)
    }
  }

  # Execute the setup script
  provisioner "remote-exec" {
    inline = [
      "echo 'Waiting for apt to be available...'",
      "until [ ! -f /var/lib/apt/lists/lock ] && [ ! -f /var/lib/dpkg/lock-frontend ]; do sleep 2; done",
      "sudo apt-get update",
      "sudo apt-get install -y dos2unix",
      "sudo dos2unix /root/setup.sh || echo 'Failed to convert setup.sh'",
      "sudo dos2unix /root/config.env || echo 'Failed to convert config.env'",
      "chmod +x /root/setup.sh",
      "/root/setup.sh"
    ]
    
    connection {
      type    = "ssh"
      user    = "root"
      host    = self.ipv4_address
      private_key = file(var.ssh_private_key_path)
      timeout = "10m"
    }
  }

  depends_on = [
    digitalocean_database_cluster.postgres,
    digitalocean_database_cluster.redis,
    digitalocean_spaces_bucket.bucket
  ]
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
  name    = var.certificate_name
  type    = "lets_encrypt"
  domains = var.domains
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
