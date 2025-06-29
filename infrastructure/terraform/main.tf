terraform {
  required_version = ">= 1.0"
  
  required_providers {
    linode = {
      source  = "linode/linode"
      version = "~> 2.0"
    }
  }
}

# Configure the Linode Provider
provider "linode" {
  token = var.linode_token
}

# Variables
variable "linode_token" {
  type        = string
  description = "Linode API token"
  sensitive   = true
}

variable "region" {
  type        = string
  description = "Linode region to deploy to"
  default     = "us-east"
}

variable "instance_type" {
  type        = string
  description = "Linode instance type"
  default     = "g6-nanode-1"  # 1GB RAM, 1 CPU, 25GB storage - $5/month
}

variable "instance_label" {
  type        = string
  description = "Label for the Linode instance"
  default     = "planning-poker"
}

variable "root_password" {
  type        = string
  description = "Root password for the instance"
  sensitive   = true
}

variable "ssh_public_key" {
  type        = string
  description = "SSH public key for access"
}

variable "planning_poker_image" {
  type        = string
  description = "Planning Poker Packer image ID"
}

variable "tags" {
  type        = list(string)
  description = "Tags to apply to resources"
  default     = ["planning-poker", "production"]
}

# Data source to get the custom image
data "linode_images" "planning_poker" {
  filter {
    name   = "label"
    values = [var.planning_poker_image]
  }
}

# Create Linode instance
resource "linode_instance" "planning_poker" {
  label           = var.instance_label
  image           = length(data.linode_images.planning_poker.images) > 0 ? data.linode_images.planning_poker.images[0].id : var.planning_poker_image
  region          = var.region
  type            = var.instance_type
  root_pass       = var.root_password
  authorized_keys = [var.ssh_public_key]
  tags            = var.tags

  # Firewall rules
  private_ip = false

  # Ensure the instance is ready before proceeding
  connection {
    type        = "ssh"
    host        = tolist(self.ipv4)[0]
    user        = "root"
    private_key = file("~/.ssh/id_rsa")  # Adjust path as needed
    timeout     = "5m"
  }

  provisioner "remote-exec" {
    inline = [
      "cloud-init status --wait",
      "systemctl start planning-poker",
      "systemctl status planning-poker --no-pager"
    ]
  }
}

# Create a Linode Firewall
resource "linode_firewall" "planning_poker" {
  label = "${var.instance_label}-firewall"
  tags  = var.tags

  # Firewall policies
  inbound_policy  = "DROP"
  outbound_policy = "ACCEPT"

  # Inbound rules
  inbound {
    label    = "allow-ssh"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "22"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  inbound {
    label    = "allow-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "80"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  inbound {
    label    = "allow-https"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "443"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  inbound {
    label    = "allow-planning-poker"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "8080"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  # Outbound rules (allow all)

  # Assign firewall to the instance
  linodes = [linode_instance.planning_poker.id]
}

# Optional: Create a Linode domain for easier access
resource "linode_domain" "planning_poker" {
  count       = var.create_domain ? 1 : 0
  domain      = var.domain_name
  type        = "master"
  description = "Domain for Planning Poker application"
  soa_email   = var.admin_email
  tags        = var.tags
}

resource "linode_domain_record" "planning_poker_a" {
  count     = var.create_domain ? 1 : 0
  domain_id = linode_domain.planning_poker[0].id
  name      = "@"
  record_type = "A"
  target    = tolist(linode_instance.planning_poker.ipv4)[0]
  ttl_sec   = 300
}

resource "linode_domain_record" "planning_poker_www" {
  count     = var.create_domain ? 1 : 0
  domain_id = linode_domain.planning_poker[0].id
  name      = "www"
  record_type = "A"
  target    = tolist(linode_instance.planning_poker.ipv4)[0]
  ttl_sec   = 300
}

# Additional variables for domain (optional)
variable "create_domain" {
  type        = bool
  description = "Whether to create a Linode domain"
  default     = false
}

variable "domain_name" {
  type        = string
  description = "Domain name for the application"
  default     = ""
}

variable "admin_email" {
  type        = string
  description = "Admin email for domain SOA record"
  default     = ""
}
