resource "prosimo_dynamic_risk" "dynamic_risk" {
  threshold  {
    name  =  "alert" 
    enabled  =  true
    value  =  80
  }
    threshold  {
    name  =  "mfa" 
    enabled  =  true
    value  =  80
  }
    threshold  {
    name  =  "lockUser" 
    enabled  =  false
    value  =  90
  }
}

