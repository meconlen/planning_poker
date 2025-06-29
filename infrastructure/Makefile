# Planning Poker Infrastructure Makefile

.PHONY: help packer-validate packer-build terraform-init terraform-plan terraform-apply terraform-destroy clean

# Default target
help:
	@echo "Planning Poker Infrastructure Commands:"
	@echo ""
	@echo "Packer Commands:"
	@echo "  packer-validate    Validate Packer configuration"
	@echo "  packer-build      Build Packer image"
	@echo ""
	@echo "Terraform Commands:"
	@echo "  terraform-init    Initialize Terraform"
	@echo "  terraform-plan    Plan Terraform deployment"
	@echo "  terraform-apply   Apply Terraform configuration"
	@echo "  terraform-destroy Destroy Terraform infrastructure"
	@echo ""
	@echo "Utility Commands:"
	@echo "  clean            Clean up temporary files"
	@echo "  check-deps       Check required dependencies"
	@echo ""
	@echo "Environment Variables Required:"
	@echo "  LINODE_TOKEN     Your Linode API token"

# Check if required tools are installed
check-deps:
	@echo "Checking dependencies..."
	@command -v packer >/dev/null 2>&1 || { echo "❌ Packer is required but not installed."; exit 1; }
	@command -v terraform >/dev/null 2>&1 || { echo "❌ Terraform is required but not installed."; exit 1; }
	@command -v curl >/dev/null 2>&1 || { echo "❌ curl is required but not installed."; exit 1; }
	@echo "✅ All dependencies are installed"

# Packer targets
packer-validate: check-deps
	@echo "Validating Packer configuration..."
	cd packer && packer validate planning-poker.pkr.hcl

packer-build: check-deps
	@echo "Building Packer image..."
	@if [ -z "$$LINODE_TOKEN" ]; then \
		echo "❌ LINODE_TOKEN environment variable is required"; \
		exit 1; \
	fi
	cd packer && packer build planning-poker.pkr.hcl

# Terraform targets
terraform-init: check-deps
	@echo "Initializing Terraform..."
	cd terraform && terraform init

terraform-plan: terraform-init
	@echo "Planning Terraform deployment..."
	cd terraform && terraform plan

terraform-apply: terraform-init
	@echo "Applying Terraform configuration..."
	cd terraform && terraform apply

terraform-destroy: terraform-init
	@echo "⚠️  Destroying Terraform infrastructure..."
	@echo "This will DELETE all resources. Are you sure? (yes/no)"
	@read answer && [ "$$answer" = "yes" ]
	cd terraform && terraform destroy

# Utility targets
clean:
	@echo "Cleaning up temporary files..."
	cd packer && rm -f manifest.json packer_cache
	cd terraform && rm -f terraform.tfstate.backup .terraform.lock.hcl
	@echo "✅ Cleanup completed"

# Development targets
dev-setup:
	@echo "Setting up development environment..."
	@if [ ! -f terraform/terraform.tfvars ]; then \
		echo "Creating terraform.tfvars from example..."; \
		cp terraform/terraform.tfvars.example terraform/terraform.tfvars; \
		echo "⚠️  Please edit terraform/terraform.tfvars with your values"; \
	fi

# Quick deployment (build image and deploy)
deploy: packer-build terraform-apply
	@echo "✅ Full deployment completed!"
	@echo ""
	@echo "Access your Planning Poker application:"
	@cd terraform && terraform output planning_poker_url

# Status check
status:
	@echo "Checking infrastructure status..."
	@cd terraform && terraform output -json | jq -r '.planning_poker_url.value' | xargs -I {} curl -I {}

# Show outputs
show:
	@cd terraform && terraform output
