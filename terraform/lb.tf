# Reserve a static IP for the Load Balancer
resource "google_compute_global_address" "default" {
  count = local.create_lb ? 1 : 0
  name  = "sketchnote-artist-ip-${local.env}"
}

# Create a Serverless Network Endpoint Group (NEG) for Cloud Run
resource "google_compute_region_network_endpoint_group" "serverless_neg" {
  count                 = local.create_lb ? 1 : 0
  name                  = "sketchnote-neg-${local.env}"
  network_endpoint_type = "SERVERLESS"
  region                = var.region
  cloud_run {
    service = google_cloud_run_v2_service.default.name
  }
}

# Create a managed SSL certificate
resource "google_compute_managed_ssl_certificate" "default" {
  count = local.create_lb ? 1 : 0
  name  = "sketchnote-cert-${local.env}"
  managed {
    domains = [var.domain]
  }
}

# Backend Service with IAP enabled
resource "google_compute_backend_service" "default" {
  count       = local.create_lb ? 1 : 0
  name        = "sketchnote-backend-${local.env}"
  port_name   = "http"
  protocol    = "HTTPS"
  timeout_sec = 300

  backend {
    group = google_compute_region_network_endpoint_group.serverless_neg[0].id
  }

  iap {
    oauth2_client_id     = var.iap_client_id
    oauth2_client_secret = var.iap_client_secret
  }
}

# Backend Bucket for Images (with CDN)
resource "google_compute_backend_bucket" "images" {
  count       = local.create_lb ? 1 : 0
  name        = "sketchnote-images-bucket-${local.env}"
  bucket_name = google_storage_bucket.images.name
  enable_cdn  = true
}

# URL Map
resource "google_compute_url_map" "default" {
  count           = local.create_lb ? 1 : 0
  name            = "sketchnote-url-map-${local.env}"
  default_service = google_compute_backend_service.default[0].id

  host_rule {
    hosts        = ["*"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = google_compute_backend_service.default[0].id

    path_rule {
      paths   = ["/images/*"]
      service = google_compute_backend_bucket.images[0].id
    }
  }
}

# Target HTTPS Proxy
resource "google_compute_target_https_proxy" "default" {
  count            = local.create_lb ? 1 : 0
  name             = "sketchnote-https-proxy-${local.env}"
  url_map          = google_compute_url_map.default[0].id
  ssl_certificates = [google_compute_managed_ssl_certificate.default[0].id]
}

# Global Forwarding Rule
resource "google_compute_global_forwarding_rule" "default" {
  count      = local.create_lb ? 1 : 0
  name       = "sketchnote-lb-${local.env}"
  target     = google_compute_target_https_proxy.default[0].id
  port_range = "443"
  ip_address = google_compute_global_address.default[0].address
}

# IAP Access Control
resource "google_iap_web_backend_service_iam_binding" "binding" {
  count               = local.create_lb ? 1 : 0
  project             = var.project_id
  web_backend_service = google_compute_backend_service.default[0].name
  role                = "roles/iap.httpsResourceAccessor"
  members = [
    for email in jsondecode(file("${path.module}/allowed_users.json")) : "user:${email}"
  ]
}
