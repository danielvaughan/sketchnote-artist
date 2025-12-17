# Terraform Deployment

The Sketchnote Artist application uses Terraform to manage its infrastructure on Google Cloud Platform.

## Prerequisites

1. [Install Terraform](https://developer.hashicorp.com/terraform/install).
2. [Install Google Cloud SDK](https://cloud.google.com/sdk/docs/install).
3. Authenticate with GCP:

    ```bash
    gcloud auth login
    gcloud auth application-default login
    ```

## Terraform Setup

1. **Navigate to the terraform directory**:

    ```bash
    cd terraform
    ```

2. **Initialize Terraform**:

    ```bash
    terraform init
    ```

3. **Configure Variables**:
    Copy the template:

    ```bash
    cp terraform.tfvars.template terraform.tfvars
    ```

    Edit `terraform.tfvars` and fill in your details:
    * `project_id`: Your GCP Project ID.
    * `domain`: The domain for your load balancer (e.g., `app.example.com`).
    * `allowed_user_emails`: List of emails allowed to access the app via IAP.
    * `iap_client_id` & `iap_client_secret`: From GCP Console -> APIs & Services -> Credentials -> OAuth 2.0 Client IDs.
    * `google_api_key`: Your Gemini API Key.

4. **Deploy**:

    ```bash
    terraform apply
    ```

    Confirm the plan by typing `yes`.

## Workspaces

The project uses Terraform workspaces to manage separate environments. There are two primary workspaces: `dev` and `prod`.

### Switching Environments

To list current workspaces:

```bash
terraform workspace list
```

To switch to a workspace (or create it):

```bash
terraform workspace select dev || terraform workspace new dev
```

### Environment Differences

The `locals.tf` file contains logic that uses the current workspace to configure environment-specific resources:

| Feature | `dev` Workspace | `prod` Workspace |
| :--- | :--- | :--- |
| **Service Name** | `sketchnote-artist-dev` | `sketchnote-artist-prod` |
| **CI/CD Branch** | `dev` | `main` |
| **Load Balancer** | Not created | Created |
| **Safety Mapping** | `default` maps to `dev` | N/A |

> [!IMPORTANT]
> Always ensure you are in the correct workspace before running `terraform apply` to avoid impacting the wrong environment.

## Post-Deployment

* **DNS Setup**: Update your DNS A record to point to the `load_balancer_ip` output by Terraform.
* **IAP Configuration**: Add the callback URL to your OAuth Client ID in GCP Console: `https://iap.googleapis.com/v1/oauth/clientIds/YOUR_CLIENT_ID:handleRedirect`.
* **Service URL**: Retrieve the deployed service URL:

    ```bash
    export SERVICE_URL=$(cd terraform && terraform output -raw service_url)
    ```
