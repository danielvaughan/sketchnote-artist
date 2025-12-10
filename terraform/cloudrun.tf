resource "google_cloud_run_v2_service" "default" {
  name     = "sketchnote-artist"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"

  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/sketchnote-repo/sketchnote-artist:${var.image_tag}"
      ports {
        container_port = 8080
      }
      env {
        name  = "GOOGLE_API_KEY"
        value = var.google_api_key
      }
    }
  }
}


