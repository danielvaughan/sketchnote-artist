locals {
  # Use workspace as environment, mapping default to dev for safety
  env = terraform.workspace == "default" ? "dev" : terraform.workspace
  
  # Determine if production
  is_prod = local.env == "prod"

  # Naming conventions
  service_name = "sketchnote-artist-${local.env}"
  
  # Cloud Build Trigger
  # Prod triggers on main, Dev triggers on dev
  branch_pattern = local.is_prod ? "^main$" : "^dev$"

  # Load Balancer is only for production
  create_lb = local.is_prod
}
