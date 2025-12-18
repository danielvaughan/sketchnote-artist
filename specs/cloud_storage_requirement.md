add # Requirement: Cloud Storage Integration for Sketchnotes on Cloud Run

## Overview

This requirement details the migration of asset storage from the local file system to Google Cloud Storage (GCS) when the Sketchnote Artist application is deployed on Google Cloud Run. This ensures scalability and persistence in a serverless environment.

## Current State

- Visual briefs and generated sketchnote images are currently saved to the local file system within the container.
- This is not persistent across Cloud Run instances and does not scale.

## Objective

Enable the application to store and serve assets from GCS buckets when running in a cloud environment, while maintaining local file system storage for local development.

## Detailed Requirements

### 1. Environment Detection & Configuration

- **Mechanism**: The application must detect its runtime environment using an environment variable.
- **Variable**: `DEPLOYMENT_MODE` (or similar).
  - If set to `cloud_run`, use GCS storage strategy.
  - If missing or set to `local`, default to existing local file system strategy.
- **Configuration**: The application configuration must support defining GCS bucket names.

### 2. Infrastructure (Terraform)

- **New Resources**: Update the existing Terraform configuration to provision two new GCS buckets:
    1. `sketchnote-visual-briefs-{env}`: For storing intermediate visual brief text/data.
    2. `sketchnote-images-{env}`: For storing the final generated sketchnote images.
- **IAM Permissions**:
  - The Cloud Run service account must be granted:
    - `roles/storage.objectCreator` (or equivalent) for writing to both buckets.
    - `roles/storage.objectViewer` for reading from the buckets (especially the images bucket).
  - Ensure the sketchnotes image bucket is configured to allow public access or appropriate signed URL generation if the UI accesses it directly (see UI Integration). *Self-correction: If the UI serves "directly" from the bucket, the bucket might need public read access or the app needs to proxy/sign URLs. Given "serve directly", public read or signed URLs are implied.*

### 3. Application Logic Changes

- **Storage Interface**: Refactor the current file saving logic into a storage interface with two implementations:
  - `LocalStorage`: Existing implementation.
  - `CloudStorage`: New implementation using the Google Cloud Storage Client Library for Go.
- **Visual Briefs**:
  - When `DEPLOYMENT_MODE=cloud_run`, write brief content to the configured briefs bucket.
- **Sketchnotes**:
  - When `DEPLOYMENT_MODE=cloud_run`, write generated images to the configured images bucket.

### 4. UI Integration

- **Serving Images**:
  - The web UI currently loads images from a local path.
  - Update the UI/Server logic to return the public GCS URL (or a signed URL) for the image when in Cloud Run mode.
  - If the bucket is public, the URL will be `https://storage.googleapis.com/<bucket-name>/<filename>`.

## Acceptance Criteria

1. **Local Development**: Application continues to work locally, saving files to disk, without requiring GCS credentials.
2. **Cloud Deployment**:
    - `terraform apply` successfully creates the two buckets and assigns permissions.
    - Cloud Run service deploys with `DEPLOYMENT_MODE=cloud_run`.
3. **Functionality**:
    - Generating a sketchnote in Cloud Run results in files appearing in the respective GCS buckets.
    - The result page in the UI correctly displays the image served from GCS.
