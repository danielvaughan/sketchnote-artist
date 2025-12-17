resource "google_project_service" "cloudbuild_api" {
  service            = "cloudbuild.googleapis.com"
  disable_on_destroy = false
}

resource "google_cloudbuild_trigger" "push_to_main" {
  name        = "deploy-${local.env}"
  description = "Triggered by push to ${local.env} branch"
  location    = var.region

  repository_event_config {
    repository = "projects/${var.project_id}/locations/${var.region}/connections/${var.github_connection_name}/repositories/${var.github_repository_id}"
    push {
      branch = local.branch_pattern
    }
  }

  filename = "cloudbuild.yaml"

  substitutions = {
    _REGION       = var.region
    _SERVICE_NAME = local.service_name
    _REPO_NAME    = "sketchnote-repo-${local.env}"
    _GOOGLE_API_KEY = "sm://projects/$${var.project_id}/secrets/GOOGLE_API_KEY-$${local.env}/versions/latest"
  }

  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"

  service_account = google_service_account.cloudbuild_sa.id

  depends_on = [google_project_service.cloudbuild_api]
}

resource "google_service_account" "cloudbuild_sa" {
  account_id   = "sketchnote-builder-${local.env}"
  display_name = "Cloud Build Service Account for Sketchnote Artist (${local.env})"
}


resource "google_project_iam_member" "cloudbuild_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.cloudbuild_sa.email}"
}

resource "google_project_iam_member" "cloudbuild_sa_user" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.cloudbuild_sa.email}"
}

resource "google_project_iam_member" "cloudbuild_artifact_registry_writer" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.cloudbuild_sa.email}"
}

resource "google_project_iam_member" "cloudbuild_logging_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloudbuild_sa.email}"
}


