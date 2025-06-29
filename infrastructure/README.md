# Planning Poker Infrastructure

This directory contains Infrastructure as Code (IaC) for deploying Planning Poker to Linode using Packer and Terraform.

## Overview

The infrastructure consists of two main components:

1. **Packer**: Creates a custom Linode image with Docker and deployment scripts pre-installed
2. **Terraform**: Deploys the Packer image and configures the infrastructure

## Prerequisites

- [Packer](https://www.packer.io/downloads) installed
- [Terraform](https://www.terraform.io/downloads) installed
- Linode API token with read/write permissions
- SSH key pair for server access

## Quick Start

### 1. Build the Packer Image

```bash
cd infrastructure/packer

# Set your Linode API token
export LINODE_TOKEN="your-linode-api-token"

# Build the image
packer build planning-poker.pkr.hcl
```

This creates a custom Linode image with:
- Ubuntu 22.04 base
- Docker and Docker Compose
- GitHub CLI for downloading releases
- Planning Poker deployment script
- Systemd service configuration

### 2. Deploy with Terraform

```bash
cd infrastructure/terraform

# Copy and configure variables
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values

# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Deploy the infrastructure
terraform apply
```

### 3. Access Your Application

After deployment, Terraform will output the instance IP and access URLs:

```bash
# Get the Planning Poker URL
terraform output planning_poker_url

# SSH into the instance
terraform output ssh_command
```

## Configuration

### Packer Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `linode_token` | Linode API token | `$LINODE_TOKEN` |
| `image_label` | Label for the created image | `planning-poker-docker` |
| `region` | Linode region | `us-east` |

### Terraform Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `linode_token` | Linode API token | Yes | - |
| `root_password` | Root password for instance | Yes | - |
| `ssh_public_key` | SSH public key | Yes | - |
| `planning_poker_image` | Packer image label | Yes | - |
| `region` | Linode region | No | `us-east` |
| `instance_type` | Instance type | No | `g6-nanode-1` |
| `instance_label` | Instance label | No | `planning-poker` |
| `create_domain` | Create Linode domain | No | `false` |
| `domain_name` | Domain name | No | `""` |
| `admin_email` | Admin email for domain | No | `""` |

## Instance Types and Pricing

| Instance Type | RAM | CPU | Storage | Monthly Cost |
|---------------|-----|-----|---------|--------------|
| `g6-nanode-1` | 1GB | 1 | 25GB | $5 |
| `g6-standard-1` | 2GB | 1 | 50GB | $10 |
| `g6-standard-2` | 4GB | 2 | 80GB | $20 |

The `g6-nanode-1` is sufficient for small to medium Planning Poker sessions.

## Firewall Configuration

The Terraform configuration creates a firewall with the following rules:

**Inbound (allowed):**
- SSH (port 22)
- HTTP (port 80)
- HTTPS (port 443)
- Planning Poker (port 8080)

**Outbound:** All traffic allowed

## Deployment Process

1. **Packer Image Build:**
   - Provisions Ubuntu 22.04 instance
   - Installs Docker, Docker Compose, GitHub CLI
   - Copies deployment script
   - Creates systemd service
   - Saves as custom image

2. **Terraform Deployment:**
   - Creates Linode instance from Packer image
   - Configures firewall rules
   - Runs startup script
   - Starts Planning Poker service
   - Optionally creates DNS records

3. **Application Startup:**
   - Downloads latest Docker image from GitHub releases
   - Starts Planning Poker container
   - Configures automatic restart

## Monitoring and Management

### Check Application Status

```bash
# Test application availability
curl -I http://YOUR_INSTANCE_IP:8080

# SSH and check Docker status
ssh root@YOUR_INSTANCE_IP
docker ps
docker logs planning-poker

# Check systemd service
systemctl status planning-poker
```

### Update to Latest Release

The deployment script automatically downloads the latest release. To manually update:

```bash
ssh root@YOUR_INSTANCE_IP
systemctl restart planning-poker
```

### View Logs

```bash
# Application deployment logs
tail -f /var/log/planning-poker-deploy.log

# Systemd service logs
journalctl -u planning-poker -f

# Docker container logs
docker logs -f planning-poker
```

## Scaling and High Availability

For production deployments, consider:

1. **Load Balancer**: Use Linode NodeBalancer for multiple instances
2. **Database**: External Redis for session storage
3. **CDN**: Linode Object Storage for static assets
4. **Monitoring**: Prometheus + Grafana
5. **Backup**: Automated instance snapshots

## Security Considerations

- Change default SSH port (edit firewall rules)
- Use strong root password
- Consider SSH key-only authentication
- Regular security updates
- Enable fail2ban for brute force protection

## Troubleshooting

### Common Issues

1. **Packer build fails:**
   - Check Linode API token permissions
   - Verify region availability
   - Check instance type quotas

2. **Terraform deployment fails:**
   - Ensure Packer image exists
   - Verify SSH key format
   - Check variable values

3. **Application not responding:**
   - Check Docker daemon status
   - Verify GitHub release download
   - Review deployment logs

### Debug Commands

```bash
# Check Packer build
packer validate planning-poker.pkr.hcl

# Debug Terraform
terraform plan -detailed-exitcode
terraform apply -auto-approve

# Instance diagnostics
ssh root@INSTANCE_IP
systemctl status docker
systemctl status planning-poker
docker ps -a
```

## Cost Optimization

- Use `g6-nanode-1` for development/testing
- Enable auto-shutdown during off-hours
- Use Linode's billing alerts
- Consider reserved instances for production

## Contributing

When modifying infrastructure:

1. Test Packer builds in development region
2. Validate Terraform plans before applying
3. Update documentation for new variables
4. Test deployment end-to-end

## Support

For infrastructure issues:
- Check Linode status page
- Review Terraform and Packer documentation
- Create GitHub issues for application problems
