
resource "prosimo_dp_profile" "test_device_posture_cards" {

    # High risk policy for mac
    inprofile_list {
      name = "test-high-risk-1643666943"
      enabled = true
      risk_level = "high"
      edr_profiles = ["high"]
      criteria {
        os = "mac"
        firewall_status = "disabled"
        disk_encryption_status = "disabled"
        domain_of_interest = []
        running_process = []
      }
    }

    # High risk policy for windows
    inprofile_list {
      name = "test-high-risk-1643666944"
      enabled = true
      risk_level = "high"
      edr_profiles = ["high"]
      criteria {
        os = "windows"
        firewall_status = "disabled"
        disk_encryption_status = "disabled"
        domain_of_interest = []
        running_process = []
      }
    }

    # Medium risk policy for mac
    inprofile_list {
      name = "test-medium-risk-1643666946"
      enabled = true
      risk_level = "medium"
      edr_profiles = ["medium"]
      criteria {
        os = "mac"
        firewall_status = "disabled"
        disk_encryption_status = "enabled"
        domain_of_interest = []
        running_process = []
      }
    }

    inprofile_list {
      name = "test-medium-risk-1643666947"
      enabled = true
      risk_level = "medium"
      edr_profiles = ["medium"]
      criteria {
        os = "windows"
        firewall_status = "disabled"
        disk_encryption_status = "enabled"
        domain_of_interest = []
        running_process = []
      }
    }

    # Low risk policy for mac
    inprofile_list {
      name = "test-low-risk-1643666949"
      enabled = true
      risk_level = "low"
      edr_profiles = ["basic"]
      criteria {
        os = "mac"
        firewall_status = "enabled"
        disk_encryption_status = "enabled"
        domain_of_interest = []
        running_process = []
      }
    }

    inprofile_list {
      name = "test-low-risk-1643666950"
      enabled = true
      risk_level = "low"
      edr_profiles = ["basic"]
      criteria {
        os = "windows"
        firewall_status = "enabled"
        disk_encryption_status = "enabled"
        domain_of_interest = []
        running_process = []
      }
    }

}




