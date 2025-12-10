resource "google_artifact_registry_repository" "repo" {
  location      = var.region
  repository_id = "sketchnote-repo"
  format        = "DOCKER"
  description   = "Docker repository for Sketchnote Artist"
}
