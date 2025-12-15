variable "project_id" {
  description = "The Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "The Google Cloud region to deploy to"
  type        = string
  default     = "us-central1"
}

variable "image_tag" {
  description = "The tag of the container image to deploy (e.g., v1)"
  type        = string
  default     = "v4"
}

variable "domain" {
  description = "The domain name for the Load Balancer SSL certificate"
  type        = string
}

variable "iap_client_id" {
  description = "The OAuth Client ID for IAP"
  type        = string
  sensitive   = true
}

variable "iap_client_secret" {
  description = "The OAuth Client Secret for IAP"
  type        = string
  sensitive   = true
}


variable "google_api_key" {
  description = "The Google API Key for Gemini"
  type        = string
  sensitive   = true
}

variable "github_repository_id" {
  description = "The Cloud Build Repository ID (usually owner-repo)"
  type        = string
  default     = "danielvaughan-sketchnote-artist"
}

variable "github_connection_name" {
  description = "The name of the Cloud Build 2nd Gen Repository Connection"
  type        = string
  default     = "sketchnote-artist"
}

