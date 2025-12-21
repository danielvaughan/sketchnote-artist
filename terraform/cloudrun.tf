resource "google_cloud_run_v2_service" "default" {
  name     = local.service_name
  location = var.region
  ingress  = local.is_prod ? "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER" : "INGRESS_TRAFFIC_ALL"

  template {
    timeout = "300s"
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/sketchnote-repo-${local.env}/sketchnote-artist:${var.image_tag}"
      ports {
        container_port = 8080
      }
      resources {
        limits = {
          cpu    = "2"
          memory = "2Gi"
        }
      }
      env {
        name = "GOOGLE_API_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.google_api_key.secret_id
            version = "latest"
          }
        }
      }
      env {
        name  = "DEPLOYMENT_MODE"
        value = "cloud_run"
      }
      env {
        name  = "GCS_BUCKET_BRIEFS"
        value = google_storage_bucket.visual_briefs.name
      }
      env {
        name  = "GCS_BUCKET_IMAGES"
        value = google_storage_bucket.images.name
      }
      env {
        name  = "GOOGLE_CLOUD_PROJECT"
        value = var.project_id
      }
      env {
        name  = "GOOGLE_CLOUD_LOCATION"
        value = var.region
      }
      env {
        name  = "DB_USER"
        value = google_sql_user.default.name
      }
      env {
        name = "DB_PASS"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_password.secret_id
            version = "latest"
          }
        }
      }
      env {
        name  = "DB_NAME"
        value = google_sql_database.default.name
      }
      env {
        name  = "DB_CONNECTION_NAME"
        value = google_sql_database_instance.default.connection_name
      }
      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }
    }
    service_account = google_service_account.run_sa.email
    vpc_access {
      connector = null # Cloud SQL Auth Proxy doesn't strictly need a connector if ipv4 is enabled
    }
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.default.connection_name]
      }
    }
  }
}

resource "google_service_account" "run_sa" {
  account_id   = "sketchnote-run-sa-${local.env}"
  display_name = "Cloud Run Service Account - ${local.env}"
}

resource "google_storage_bucket_iam_member" "run_sa_briefs_creator" {
  bucket = google_storage_bucket.visual_briefs.name
  role   = "roles/storage.objectCreator"
  member = "serviceAccount:${google_service_account.run_sa.email}"
}

resource "google_storage_bucket_iam_member" "run_sa_images_creator" {
  bucket = google_storage_bucket.images.name
  role   = "roles/storage.objectCreator"
  member = "serviceAccount:${google_service_account.run_sa.email}"
}

resource "google_project_iam_member" "run_sa_cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.run_sa.email}"
}



resource "google_cloud_run_v2_service_iam_member" "noauth" {
  count    = local.is_prod ? 0 : 1
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_secret_manager_secret_iam_member" "run_sa_db_pass_accessor" {
  secret_id = google_secret_manager_secret.db_password.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.run_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "run_sa_google_api_key_accessor" {
  secret_id = google_secret_manager_secret.google_api_key.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.run_sa.email}"
}
