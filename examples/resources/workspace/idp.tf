
resource "prosimo_idp" "google_ws" {
    idp_account = "google"
    auth_type = "oidc"
    account_url = "https://accounts.google.com"
    api_cred_provided = "yes"
    oidc {
      client_id = "227796663766-6db0damqp8r104269fs2ih8i3e3uijhf.apps.googleusercontent.com"
      secret_id = "var.google-secret"
      admin_email = "test0319@myeventarena.com"
      customer_id = "C02tdt1kg"
      domain = "myevetarena7.com"
      # file_path = "/User/file/path"
    }
    select_type = "partner"
    partner {
    user_domain =  ["myevetarena.com"]
    apps = ["all"]
    }
}


resource "prosimo_idp" "Ping_one" {
    idp_account = "ping-one"
    auth_type = "oidc"
    account_url = "https://auth.pingone.com/fe92b63d-a687-4964-be6a-047054e53ccb"
    api_cred_provided = "yes"
    oidc {
      api_client_id = "76f8e416-df36-46b2-bfa8-cb60de66e30f"
      api_secret_id = "var.pingone-secret"
      region = "default"
      env_id = "fe92b63d-a687-4964-be6a-047054e53ccb"
    }
    select_type = "partner"
    partner {
      user_domain =  ["myevetarena.com", "test1.com"]
      apps = ["all"]
    }
}