# Either of these filter options can be used.Max one filter condition per request. 
data "prosimo_idp" "idp_account" {
  # filter = "idpname==azure_ad&selecttype!=primary,authtype==oidc"
  filter = "idpname!=azure_ad"
  # filter = "idpname=@azure"
}

  output "idp_details" {
  description = "idp"
  value       = data.prosimo_idp.idp_account
}