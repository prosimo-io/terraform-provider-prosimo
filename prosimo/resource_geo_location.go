package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGeoLocation() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify geo Velocity settings.",
		CreateContext: resourceGeoLocationUpdate,
		ReadContext:   resourceGeoLocationRead,
		DeleteContext: resourceGeoLocationDelete,
		UpdateContext: resourceGeoLocationUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"allow_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Add Locations you don’t want to have blocked",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"state_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"country_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceGeoLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)

	locByCityCode := client.GetGeoVelocity{}
	userinputList := []client.GetGeoVelocity{}
	payload := &client.PostGeoVelocity{}

	addLocations := []client.GetGeoVelocity{}
	deleteLocations := []client.GetGeoVelocity{}

	addLocationsByCityCode := []client.GetGeoVelocity{}
	deleteLocationsByCityCode := []client.GetGeoVelocity{}

	if v, ok := d.GetOk("allow_list"); ok {
		userinputListBYCityName := v.([]interface{})
		for _, userInput := range userinputListBYCityName {
			userinput := client.GetGeoVelocity{}
			val := userInput.(map[string]interface{})
			if city, ok := val["city_name"].(string); ok {
				userinput.CityName = city
			}
			if state, ok := val["state_name"].(string); ok {
				userinput.StateName = state
			}
			if country, ok := val["country_name"].(string); ok {
				userinput.CountryName = country
			}
			userinputList = append(userinputList, userinput)
		}
	}

	geolocationList, err := prosimoClient.GetGeoVelocity(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(geolocationList.GetGeoVelocities) > 0 {
		for _, val := range userinputList {
			if val.CountryName != "" {
				isAddList := true
				for _, val1 := range geolocationList.GetGeoVelocities {
					if val.CountryName == val1.CountryName && val.StateName == val1.StateName && val.CityName == val1.CityName {
						isAddList = false
					} else if val.CityName == "" {
						if val.CountryName == val1.CountryName && val.StateName == val1.StateName {
							isAddList = false
						}
					} else if val.StateName == "" {
						if val.CountryName == val1.CountryName {
							isAddList = false
						}
					} else {
						continue
					}

				}

				if isAddList == true {
					addLocations = append(addLocations, val)
				}
			} else {
				return diag.Errorf("Invalid Input, Country Missing")
			}
		}
	} else {
		addLocations = userinputList
	}
	res1, err := prosimoClient.GetCityCode(ctx, addLocations)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, cityCode := range res1 {
		locByCityCode.CityCode = cityCode
		addLocationsByCityCode = append(addLocationsByCityCode, locByCityCode)
	}
	if len(userinputList) > 0 {
		for _, val := range geolocationList.GetGeoVelocities {
			isAddDeleteList := true
			for _, val1 := range userinputList {
				if val.CountryName == val1.CountryName && val.StateName == val1.StateName && val.CityName == val1.CityName {
					isAddDeleteList = false
				} else if val.CityName == "?" {
					if val.CountryName == val1.CountryName && val.StateName == val1.StateName {
						isAddDeleteList = false
					}
				} else if val.StateName == "?" {
					if val.CountryName == val1.CountryName {
						isAddDeleteList = false
					}
				} else {
					continue
				}
			}
			if isAddDeleteList == true {
				deleteLocations = append(deleteLocations, *val)
			}
		}
	} else {
		for _, val := range geolocationList.GetGeoVelocities {
			deleteLocations = append(deleteLocations, *val)
		}
	}

	for _, cityCode := range deleteLocations {
		locByCityCode.CityCode = cityCode.CityCode
		deleteLocationsByCityCode = append(deleteLocationsByCityCode, locByCityCode)
	}

	payload.AddLocations = addLocationsByCityCode
	payload.DeleteLocations = deleteLocationsByCityCode

	createdGeoLoc, err := prosimoClient.UpdateGeoVelocity(ctx, payload)
	_ = createdGeoLoc
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGeoLocationRead(ctx, d, meta)
}

func resourceGeoLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	geolocationList, err := prosimoClient.GetGeoVelocity(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var geoLoc *client.GetGeoVelocity
	var userInput []interface{}
	for _, retresponse := range geolocationList.GetGeoVelocities {
		geoLoc = retresponse
		d.Set("city_name", geoLoc.CityName)
		d.Set("state_name", geoLoc.StateName)
		d.Set("country_name", geoLoc.CountryName)
		userInput = append(userInput, geoLoc)
	}
	d.Set("allow_list", userInput)
	return diags
}

func resourceGeoLocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	deleteLocations := []client.GetGeoVelocity{}
	locByCityCode := client.GetGeoVelocity{}
	addLocationsByCityCode := []client.GetGeoVelocity{}
	deleteLocationsByCityCode := []client.GetGeoVelocity{}
	prosimoClient := meta.(*client.ProsimoClient)
	payload := &client.PostGeoVelocity{}
	geolocationList, err := prosimoClient.GetGeoVelocity(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, val := range geolocationList.GetGeoVelocities {
		deleteLocations = append(deleteLocations, *val)
	}
	for _, GeoLoc := range deleteLocations {
		locByCityCode.CityCode = GeoLoc.CityCode
		deleteLocationsByCityCode = append(deleteLocationsByCityCode, locByCityCode)
	}
	payload.AddLocations = addLocationsByCityCode
	payload.DeleteLocations = deleteLocationsByCityCode
	deletetedGeoLoc, err := prosimoClient.UpdateGeoVelocity(ctx, payload)
	_ = deletetedGeoLoc
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] Deleted geo locations")
	d.SetId("")
	return diags
}
