resource "google_artifact_registry_repository" "repo" {
  location      = var.region
  repository_id = "sketchnote-repo-${local.env}"
  format        = "DOCKER"
  description   = "Docker repository for Sketchnote Artist (${local.env})"
}
