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
        name  = "GOOGLE_API_KEY"
        value = var.google_api_key
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
    }
    service_account = google_service_account.run_sa.email
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



resource "google_cloud_run_v2_service_iam_member" "noauth" {
  count    = local.is_prod ? 0 : 1
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
