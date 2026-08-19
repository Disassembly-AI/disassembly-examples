terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region
}

# API key in Secret Manager (set the value out-of-band).
resource "google_secret_manager_secret" "api_key" {
  secret_id = "disassembly-api-key"
  replication {
    auto {}
  }
}

resource "google_cloud_run_v2_job" "scan" {
  name     = "${var.prefix}-scan"
  location = var.region
  template {
    template {
      containers {
        image = "ghcr.io/disassembly-ai/toolkit:latest"
        args  = ["scan", var.target, "--ci"]
        env {
          name = "DISASSEMBLY_API_KEY"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.api_key.secret_id
              version = "latest"
            }
          }
        }
      }
    }
  }
}

# Weekly trigger via Cloud Scheduler -> Cloud Run Admin API run endpoint.
resource "google_cloud_scheduler_job" "weekly" {
  name     = "${var.prefix}-scan-weekly"
  schedule = "0 6 * * 1"
  http_target {
    http_method = "POST"
    uri         = "https://run.googleapis.com/v2/${google_cloud_run_v2_job.scan.id}:run"
    oauth_token {
      service_account_email = var.scheduler_sa_email
    }
  }
}
