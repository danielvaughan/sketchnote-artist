resource "google_secret_manager_secret" "google_api_key" {
  secret_id = "GOOGLE_API_KEY"

  replication {
    auto {}
  }

  depends_on = [google_project_service.secret_manager_api]
}

resource "google_secret_manager_secret_iam_member" "cloudbuild_secret_accessor" {
  project   = var.project_id
  secret_id = google_secret_manager_secret.google_api_key.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloudbuild_sa.email}"
}
