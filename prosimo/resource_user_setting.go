package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserSettings() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify user settings.",
		CreateContext: resourceUserSettingsUpdate,
		ReadContext:   resourceUserSettingsRead,
		DeleteContext: resourceUserSettingsDelete,
		UpdateContext: resourceUserSettingsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allow_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "To prevent known users from getting locked out erroneously, their email addresses can be added to this list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User Email",
						},
						"reason": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Reason to allow",
						},
						"createdtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceUserSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)

	userinputList := []client.Users{}
	payload := &client.PostAllowlist{}

	addUsers := []client.Users{}
	deleteUsers := []client.Users{}

	if v, ok := d.GetOk("allow_list"); ok {
		userinputListBYCityName := v.([]interface{})
		for _, userInput := range userinputListBYCityName {
			userinput := client.Users{}
			val := userInput.(map[string]interface{})
			if email, ok := val["email"].(string); ok {
				userinput.Email = email
			}
			if reason, ok := val["reason"].(string); ok {
				userinput.Reason = reason
			}
			userinputList = append(userinputList, userinput)
		}
	}
	existingUsers, err := prosimoClient.SearchAllowList(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(existingUsers.GetAllowlist.Records) > 0 {
		for _, inputUser := range userinputList {
			isAddList := true
			for _, existVal := range existingUsers.GetAllowlist.Records {
				if inputUser.Email == existVal.Email {
					isAddList = false
				} else {
					continue
				}
			}
			if isAddList == true {
				addUsers = append(addUsers, inputUser)
			}
		}
	} else {
		addUsers = userinputList
	}
	if len(userinputList) > 0 {
		for _, existVal := range existingUsers.GetAllowlist.Records {
			isAddDeleteList := true
			for _, inputVal := range userinputList {
				if existVal.Email == inputVal.Email {
					isAddDeleteList = false
				} else {
					continue
				}
			}
			if isAddDeleteList == true {
				deleteUsers = append(deleteUsers, existVal)
			}
		}
	} else {
		for _, existingVal := range existingUsers.GetAllowlist.Records {
			deleteUsers = append(deleteUsers, existingVal)
		}
	}
	payload.AddUsers = addUsers
	payload.DeleteUsers = deleteUsers

	createUser, err := prosimoClient.UpdateAllowList(ctx, payload)
	_ = createUser
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceUserSettingsRead(ctx, d, meta)
}

func resourceUserSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	UserAllowList, err := prosimoClient.SearchAllowList(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var geoLoc *client.Users
	var userInput []interface{}
	//d.Set("totalCount", UserAllowList.GetAllowlist.TotalCount)
	for _, retresponse := range UserAllowList.GetAllowlist.Records {
		geoLoc = &retresponse
		d.Set("email", geoLoc.Email)
		d.Set("reason", geoLoc.Reason)
		d.Set("createdTime", geoLoc.CreatedTime)
		userInput = append(userInput, geoLoc)
	}
	d.Set("records", userInput)
	d.SetId("User Settings")
	return diags
}

func resourceUserSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	err := prosimoClient.DeleteAllUsers(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] Deleted User AllowList")
	d.SetId("")
	return diags
}
