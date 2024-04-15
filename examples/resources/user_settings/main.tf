resource "prosimo_user_settings" "user_settings" {
    allow_list {
    email   = "def@def.gz"  
    reason = "def" 
  }    
    allow_list {
    email   = "abd@def.gz"  
    reason = "def abc" 
  }         
}
