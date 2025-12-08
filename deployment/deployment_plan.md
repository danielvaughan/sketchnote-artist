# Deployment Plan: Sketchnote Artist on Cloud Run with IAP

This document outlines the steps to deploy the Sketchnote Artist application to Google Cloud Run and protect it with Identity-Aware Proxy (IAP).

## Prerequisites

1.  **Google Cloud Project**: An active project with billing enabled.
2.  **Domain Name**: Required for IAP (Load Balancer requires an SSL certificate).
3.  **gcloud CLI**: Installed and authorized (`gcloud auth login`).

## 1. Containerization

We need to package the application and its web assets into a Docker image.

### Dockerfile

A `Dockerfile` has been created in this directory. It uses a multi-stage build:
1.  **Builder**: Compiles the Go binary.
2.  **Runtime**: Minimal image (Debian Slim) containing the binary and web assets.

**Action**: Copy `deployment/Dockerfile` to the project root.
```bash
cp deployment/Dockerfile .
```

## 2. Infrastructure Setup (Google Cloud)

### Enable APIs
```bash
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  compute.googleapis.com \
  iap.googleapis.com
```

### Create Artifact Registry
```bash
gcloud artifacts repositories create sketchnote-repo \
  --repository-format=docker \
  --location=us-central1 \
  --description="Docker repository for Sketchnote Artist"
```

### Build and Push Image
```bash
gcloud builds submit --tag us-central1-docker.pkg.dev/YOUR_PROJECT_ID/sketchnote-repo/sketchnote-artist:v1 .
```

## 3. Deploy to Cloud Run

Deploy the service initially. Note that for IAP, we will put it behind a Load Balancer later.

```bash
gcloud run deploy sketchnote-artist \
  --image us-central1-docker.pkg.dev/YOUR_PROJECT_ID/sketchnote-repo/sketchnote-artist:v1 \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars GOOGLE_API_KEY=your_api_key_here
```
*Note: We allow unauthenticated access at the **Cloud Run level** because the Load Balancer will handle authentication. Alternatively, we can restrict Cloud Run to only accept traffic from the Load Balancer (ingress=internal-and-cloud-load-balancing).*

## 4. Identity-Aware Proxy (IAP) Setup

IAP requires a Global External HTTP(S) Load Balancer.

### Step 4.1: Reserve a Static IP
```bash
gcloud compute addresses create sketchnote-ip --global
```
*Note the IP address: `gcloud compute addresses describe sketchnote-ip --global --format="get(address)"`*

### Step 4.2: Domain Configuration
Update your domain's DNS A record to point to the IP address reserved above.

### Step 4.3: OAuth Consent Screen
1.  Go to **APIs & Services > OAuth consent screen**.
2.  Select **Internal** (if you have a Workspace organization) or **External** (for personal projects, requires adding test users).
3.  Fill in the application steps.

### Step 4.4: Create OAuth Credentials
1.  Go to **APIs & Services > Credentials**.
2.  Create **OAuth client ID**.
3.  Application type: **Web application**.
4.  Name: `IAP-Sketchnote`.
5.  Authorized redirect URIs: `https://YOUR_DOMAIN/_gcp_gatekeeper/authenticate` (replace `YOUR_DOMAIN` with your actual domain).
6.  Note the **Client ID** and **Client Secret**.

### Step 4.5: Setup Load Balancer with IAP

This is complex to do via CLI alone due to certificate provisioning. Using the Console is recommended for the first time, but here is the process:

1.  **Create Serverless NEG**: Network Endpoint Group pointing to your Cloud Run service.
    ```bash
    gcloud compute network-endpoint-groups create sketchnote-neg \
        --region=us-central1 \
        --network-endpoint-type=serverless \
        --cloud-run-service=sketchnote-artist
    ```
2.  **Create Backend Service**:
    ```bash
    gcloud compute backend-services create sketchnote-backend \
        --global \
        --iap=enabled,oauth2-client-id=CLIENT_ID,oauth2-client-secret=CLIENT_SECRET
    ```
3.  **Add NEG to Backend Service**:
    ```bash
    gcloud compute backend-services add-backend sketchnote-backend \
        --global \
        --network-endpoint-group=sketchnote-neg \
        --network-endpoint-group-region=us-central1
    ```
4.  **Create URL Map & Proxy**: Standard LB setup.
5.  **Create Forwarding Rule**: Links the IP to the Proxy.

### Step 4.6: Access Control
Grant access to specific users:
```bash
gcloud iap web add-iam-policy-binding \
    --resource-type=backend-services \
    --service=sketchnote-backend \
    --member='user:allowed-user@gmail.com' \
    --role='roles/iap.httpsResourceAccessor'
```

## Alternative: Simpler Auth Option
If setting up a Domain and Load Balancer is too heavy, you can use **Cloud Run's native authentication**:

1.  Deploy with `--no-allow-unauthenticated`.
2.  Add users:
    ```bash
    gcloud run services add-iam-policy-binding sketchnote-artist \
      --region us-central1 \
      --member='user:allowed-user@gmail.com' \
      --role='roles/run.invoker'
    ```
3.  Users must access the URL using a tool that sends an identity token (like `gcloud run services proxy` or a browser extension), OR you can use **IAP** (as requested) which is the robust way for browser-based access.

For a browser-based "Login with Google" flow, **IAP (Steps 1-4.6) is the correct solution**.

## Notes on State
Cloud Run is stateless. Images generated in the `sketchnotes/` directory will **disappear** when the container restarts. To persist images, consider:
1.  Using **Google Cloud Storage** to save images.
2.  Updating the code to write to a GCS bucket instead of local disk.
