resource "google_sql_database_instance" "default" {
  name             = "sketchnote-db-${local.env}"
  database_version = "POSTGRES_15"
  region           = var.region

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled = true
    }
  }

  deletion_protection = local.is_prod
}

resource "google_sql_database" "default" {
  name     = "sketchnote"
  instance = google_sql_database_instance.default.name
}

resource "google_sql_user" "default" {
  name     = "sketchnote-user"
  instance = google_sql_database_instance.default.name
  password = var.db_password
}

# Store database password in Secret Manager
resource "google_secret_manager_secret" "db_password" {
  secret_id = "DB_PASS-${local.env}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = var.db_password
}

variable "db_password" {
  description = "Password for the database user"
  type        = string
  sensitive   = true
}
