# Planning Poker Linode Deployment - Quick Reference

## Prerequisites
- Linode account with API token
- Packer and Terraform installed
- SSH key pair

## One-Command Deployment

```bash
# Set environment
export LINODE_TOKEN="your-token-here"

# Deploy everything
cd infrastructure
make deploy
```

## Step-by-Step Deployment

### 1. Build Custom Image
```bash
cd infrastructure/packer
export LINODE_TOKEN="your-token"
packer build planning-poker.pkr.hcl
```

### 2. Configure Terraform
```bash
cd infrastructure/terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
```

### 3. Deploy Infrastructure
```bash
terraform init
terraform plan
terraform apply
```

## Key Configuration

**terraform.tfvars:**
```hcl
linode_token         = "your-linode-api-token"
root_password        = "secure-password-123"
ssh_public_key       = "ssh-rsa AAAAB3NzaC1yc2E..."
planning_poker_image = "planning-poker-docker-2025-06-29"
instance_type        = "g6-nanode-1"  # $5/month
```

## Post-Deployment

**Get access information:**
```bash
terraform output planning_poker_url  # http://IP:8080
terraform output ssh_command         # ssh root@IP
```

**Check status:**
```bash
curl -I $(terraform output -raw planning_poker_url)
```

## Management Commands

```bash
# SSH into server
ssh root@$(terraform output -raw instance_ip)

# Check application
docker ps
systemctl status planning-poker

# View logs
journalctl -u planning-poker -f
docker logs planning-poker

# Restart application
systemctl restart planning-poker
```

## Cost: ~$5/month for g6-nanode-1 instance
