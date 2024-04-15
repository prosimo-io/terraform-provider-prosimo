# resource "prosimo_cloud_creds" "gcp" {
#   cloud_type = "GCP"
#   nickname   = "demo-1"
#   gcp {
#     file_path = "/Path/To/File"
#   }
# }

resource "prosimo_cloud_creds" "bulk" {
  cloud_type = "GCP"
  acctid  = "demo-1"
exterid = vale
}


# resource "prosimo_cloud_creds" "azure" {
#   cloud_type = "AZURE"
#   nickname   = "replace_me_with_nickname"

#   azure {
#     subscription_id = var.subscription_id
#     tenant_id       = var.tenant_id
#     client_id       = var.client_id
#     secret_id       = var.secret_id
#   }
# }

# resource "prosimo_cloud_creds" "aws" {
#   cloud_type = "AWS"
#   nickname   = "replace_me_with_nickname"

#   aws {
#     preffered_auth = "AWSKEY"

#     access_keys {
#       access_key_id = var.access_key_id
#       secret_key_id = var.secret_key_id
#     }
#   }
# }

data "prosimo_cloud_creds" "cloud_creds"{
    filter="nickname==prosimo-aws-iam"
}

output "cloud_creds_output" {
  description = "cloud_creds"
  value       = data.prosimo_cloud_creds.cloud_creds
}







