# Reserve a static IP for the Load Balancer
resource "google_compute_global_address" "default" {
  name = "sketchnote-ip"
}

# Create a Serverless Network Endpoint Group (NEG) for Cloud Run
resource "google_compute_region_network_endpoint_group" "serverless_neg" {
  name                  = "sketchnote-neg"
  network_endpoint_type = "SERVERLESS"
  region                = var.region
  cloud_run {
    service = google_cloud_run_v2_service.default.name
  }
}

# Create a managed SSL certificate
resource "google_compute_managed_ssl_certificate" "default" {
  name = "sketchnote-cert"
  managed {
    domains = [var.domain]
  }
}

# Backend Service with IAP enabled
resource "google_compute_backend_service" "default" {
  name        = "sketchnote-backend"
  port_name   = "http"
  protocol    = "HTTPS"
  timeout_sec = 30

  backend {
    group = google_compute_region_network_endpoint_group.serverless_neg.id
  }

  iap {
    oauth2_client_id     = var.iap_client_id
    oauth2_client_secret = var.iap_client_secret
  }
}

# URL Map
resource "google_compute_url_map" "default" {
  name            = "sketchnote-url-map"
  default_service = google_compute_backend_service.default.id
}

# Target HTTPS Proxy
resource "google_compute_target_https_proxy" "default" {
  name             = "sketchnote-https-proxy"
  url_map          = google_compute_url_map.default.id
  ssl_certificates = [google_compute_managed_ssl_certificate.default.id]
}

# Global Forwarding Rule
resource "google_compute_global_forwarding_rule" "default" {
  name       = "sketchnote-lb"
  target     = google_compute_target_https_proxy.default.id
  port_range = "443"
  ip_address = google_compute_global_address.default.address
}

# IAP Access Control
resource "google_iap_web_backend_service_iam_binding" "binding" {
  project             = var.project_id
  web_backend_service = google_compute_backend_service.default.name
  role                = "roles/iap.httpsResourceAccessor"
  members = [
    for email in jsondecode(file("${path.module}/allowed_users.json")) : "user:${email}"
  ]
}
