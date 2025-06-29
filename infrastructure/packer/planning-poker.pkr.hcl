packer {
  required_plugins {
    linode = {
      version = ">= 1.0.0"
      source  = "github.com/linode/linode"
    }
  }
}

variable "linode_token" {
  type        = string
  description = "Linode API token"
  sensitive   = true
  default     = env("LINODE_TOKEN")
}

variable "image_label" {
  type        = string
  description = "Label for the created image"
  default     = "planning-poker-docker"
}

variable "region" {
  type        = string
  description = "Linode region"
  default     = "us-east"
}

source "linode" "planning_poker" {
  linode_token    = var.linode_token
  image           = "linode/ubuntu22.04"
  region          = var.region
  instance_type   = "g6-nanode-1"
  instance_label  = "packer-planning-poker-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  image_label     = "${var.image_label}-${formatdate("YYYY-MM-DD", timestamp())}"
  image_description = "Planning Poker application image with Docker pre-installed"
  
  ssh_username = "root"
  
  # Tags for organization
  tags = ["planning-poker", "docker", "packer"]
}

build {
  name = "planning-poker-image"
  sources = ["source.linode.planning_poker"]

  # Update system and install Docker
  provisioner "shell" {
    inline = [
      "apt-get update",
      "apt-get upgrade -y",
      
      # Install Docker
      "apt-get install -y ca-certificates curl gnupg lsb-release",
      "mkdir -p /etc/apt/keyrings",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg",
      "echo \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\" | tee /etc/apt/sources.list.d/docker.list > /dev/null",
      "apt-get update",
      "apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
      
      # Start and enable Docker
      "systemctl start docker",
      "systemctl enable docker",
      
      # Install useful tools
      "apt-get install -y htop curl wget unzip jq",
      
      # Install GitHub CLI for downloading releases
      "curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg",
      "chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg",
      "echo \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main\" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null",
      "apt-get update",
      "apt-get install -y gh",
      
      # Clean up
      "apt-get autoremove -y",
      "apt-get autoclean",
      
      # Verify installations
      "docker --version",
      "docker-compose --version",
      "gh --version"
    ]
  }

  # Create planning poker service directory
  provisioner "shell" {
    inline = [
      "mkdir -p /opt/planning-poker",
      "chown root:root /opt/planning-poker",
      "chmod 755 /opt/planning-poker"
    ]
  }

  # Copy deployment script
  provisioner "file" {
    source      = "scripts/deploy-planning-poker.sh"
    destination = "/opt/planning-poker/deploy.sh"
  }

  # Make deployment script executable and create systemd service
  provisioner "shell" {
    inline = [
      "chmod +x /opt/planning-poker/deploy.sh",
      
      # Create systemd service for planning poker
      "cat > /etc/systemd/system/planning-poker.service << 'EOF'",
      "[Unit]",
      "Description=Planning Poker Application",
      "Requires=docker.service",
      "After=docker.service",
      "",
      "[Service]",
      "Type=oneshot",
      "RemainAfterExit=yes",
      "ExecStart=/opt/planning-poker/deploy.sh",
      "ExecStop=/usr/bin/docker stop planning-poker",
      "TimeoutStartSec=300",
      "",
      "[Install]",
      "WantedBy=multi-user.target",
      "EOF",
      
      # Enable the service (but don't start it yet)
      "systemctl daemon-reload",
      "systemctl enable planning-poker.service"
    ]
  }

  # Final system preparation
  provisioner "shell" {
    inline = [
      # Clear logs and history for cleaner image
      "truncate -s 0 /var/log/*log",
      "history -c",
      "cat /dev/null > ~/.bash_history",
      
      # Ensure Docker daemon is running for next boot
      "systemctl is-enabled docker"
    ]
  }

  post-processor "manifest" {
    output = "manifest.json"
    strip_path = true
  }
}
