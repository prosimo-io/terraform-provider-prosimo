resource "prosimo_certificates" "test1" {
    upload_domain_cert {
        cert_path = "/Users/sibaprasadtripathy/workspace/repos/certs/app_onboarding/domain.crt"
        private_key_path = "/Users/sibaprasadtripathy/workspace/repos/certs/app_onboarding/domain.key"
    }
}


