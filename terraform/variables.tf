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
  default     = "v5"
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

variable "allowed_user_emails" {
  description = "A list of email addresses of users to allow access via IAP"
  type        = list(string)
}

variable "google_api_key" {
  description = "The Google API Key for Gemini"
  type        = string
  sensitive   = true
}
