resource "google_artifact_registry_repository" "repo" {
  provider = google-beta
  location = var.region
  repository_id = "sketchnote-repo-${local.env}"
  format        = "DOCKER"
  description   = "Docker repository for Sketchnote Artist (${local.env})"

  vulnerability_scanning_config {
    enablement_config = "ENABLED"
  }
}
