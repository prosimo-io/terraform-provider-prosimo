

resource "prosimo_edr_integration" "crowdstrike" {
    crowdstrike {
        name = "demo3"
        vendor = "CrowdStrike"
        criteria {
            sensor_active = "enabled"
            status = "na"
            zta_score {
                from = 1
                to = 3
            }
        }
    }
    
    crowdstrike {
        name = "demo2"
        vendor = "CrowdStrike"
        criteria {
            sensor_active = "na"
            status = "na"
            zta_score {
                from = 2
                to = 3
            }
        }
    }
    
      crowdstrike {
          name = "demo4"
          vendor = "CrowdStrike"
          criteria {
              sensor_active = "disabled"
              status = "disabled"
              zta_score {
                  from = 2
                  to = 3
              }
          }
      }

    crowdstrike {
        name = "demo9"
        vendor = "CrowdStrike"
        criteria {
            sensor_active = "disabled"
            status = "disabled"
            zta_score {
                from = 2
                to = 5
            }
        }
    }

    crowdstrike {
        name = "demo5"
        vendor = "CrowdStrike"
        criteria {
            sensor_active = "disabled"
            status = "disabled"
            zta_score {
                from = 2
                to = 3
            }
        }
    }
}