data "prosimo_edge" "test-app" {
filter = "cloudtype=@GC&cloudregion!=us-east1"
}

output "app_onboard_output" {
  description = "edgelist"
  value       = data.prosimo_edge.test-app
}
