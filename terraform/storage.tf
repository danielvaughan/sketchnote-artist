resource "google_storage_bucket" "visual_briefs" {
  name          = "sketchnote-visual-briefs-${var.project_id}-${local.env}"
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true
}

resource "google_storage_bucket" "images" {
  name          = "sketchnote-images-${var.project_id}-${local.env}"
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true
}

# Make images public
resource "google_storage_bucket_iam_member" "public_image_access" {
  bucket = google_storage_bucket.images.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
