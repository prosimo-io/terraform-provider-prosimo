# Either of these filter options can be used.Max one filter condition per request. 
data "prosimo_s3bucket" "s3_bucket" {
  input_nickname = "prosimo-aws-app-iam"
  # input_region = "us-east-2"
}

output "s3_aws_creds" {
    description = "aws_s3_bucket_credentials"
    value       = data.prosimo_s3bucket.s3_bucket.data
}