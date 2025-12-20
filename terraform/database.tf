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

  deletion_protection = false # Set to true for production
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

variable "db_password" {
  description = "Password for the database user"
  type        = string
  sensitive   = true
}
