# Outputs
output "instance_id" {
  description = "ID of the Linode instance"
  value       = linode_instance.planning_poker.id
}

output "instance_ip" {
  description = "Public IP address of the instance"
  value       = tolist(linode_instance.planning_poker.ipv4)[0]
}

output "instance_label" {
  description = "Label of the instance"
  value       = linode_instance.planning_poker.label
}

output "planning_poker_url" {
  description = "URL to access Planning Poker"
  value       = "http://${tolist(linode_instance.planning_poker.ipv4)[0]}:8080"
}

output "ssh_command" {
  description = "SSH command to connect to the instance"
  value       = "ssh root@${tolist(linode_instance.planning_poker.ipv4)[0]}"
}

output "firewall_id" {
  description = "ID of the firewall"
  value       = linode_firewall.planning_poker.id
}

output "domain_name" {
  description = "Domain name (if created)"
  value       = var.create_domain ? var.domain_name : null
}

output "status_check" {
  description = "Commands to check application status"
  value = {
    curl_test     = "curl -I http://${tolist(linode_instance.planning_poker.ipv4)[0]}:8080"
    ssh_logs      = "ssh root@${tolist(linode_instance.planning_poker.ipv4)[0]} 'journalctl -u planning-poker -f'"
    docker_status = "ssh root@${tolist(linode_instance.planning_poker.ipv4)[0]} 'docker ps'"
  }
}
