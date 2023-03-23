package client

import (
	"context"
	"fmt"
	"log"
)

type GetGeoVelocityList struct {
	GetGeoVelocities []*GetGeoVelocity `json:"data,omitempty"`
}

type GetGeoVelocityListRead struct {
	GetGeoVelocities []*GetGeoVelocityRead `json:"data,omitempty"`
}
type PostGeoVelocityData struct {
	PostGeoVelocity *PostGeoVelocity `json:"data,omitempty"`
}

type PostGeoVelocity struct {
	AddLocations    []GetGeoVelocity `json:"addLocations,omitempty"`
	DeleteLocations []GetGeoVelocity `json:"deleteLocations,omitempty"`
}

type Location struct {
	CityName    string `json:"cityName,omitempty"`
	StateName   string `json:"stateName,omitempty"`
	CountryName string `json:"countryName,omitempty"`
}

type GetGeoVelocity struct {
	CityCode        int    `json:"cityCode,omitempty"`
	CityName        string `json:"cityName,omitempty"`
	CountryCodeis02 string `json:"countryCodeIso2,omitempty"`
	Region          string `json:"regionCode,omitempty"`
	StateName       string `json:"stateName,omitempty"`
	CountryName     string `json:"countryName,omitempty"`
}

type GetGeoVelocityRead struct {
	CityCode        int    `json:"city_code,omitempty"`
	CityName        string `json:"city_name,omitempty"`
	CountryCodeis02 string `json:"country_code_iso2,omitempty"`
	Region          string `json:"region_code,omitempty"`
	StateName       string `json:"state_name,omitempty"`
	CountryName     string `json:"country_name,omitempty"`
}

func (gv GetGeoVelocity) String() string {
	return fmt.Sprintf("{CityCode:%d, CityName:%s, CountryCodeis02:%s, Region:%s, StateName:%s, CountryName:%s}", gv.CityCode, gv.CityName, gv.CountryCodeis02, gv.Region, gv.StateName, gv.CountryName)
}

func (prosimoClient *ProsimoClient) UpdateGeoVelocity(ctx context.Context, geoVelocityops *PostGeoVelocity) (*PostGeoVelocityData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", GeoVelocityEndpoint, geoVelocityops)
	if err != nil {
		return nil, err
	}

	gv := &PostGeoVelocityData{}
	_, err = prosimoClient.api_client.Do(ctx, req, gv)
	if err != nil {
		return nil, err
	}

	return gv, nil

}

func (prosimoClient *ProsimoClient) GetGeoVelocity(ctx context.Context) (*GetGeoVelocityList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", GeoVelocityEndpoint, nil)
	if err != nil {
		return nil, err
	}

	geoVelocityList := &GetGeoVelocityList{}
	_, err = prosimoClient.api_client.Do(ctx, req, geoVelocityList)
	if err != nil {
		return nil, err
	}

	log.Println("geoVelocityList", geoVelocityList)
	return geoVelocityList, nil

}

func (prosimoClient *ProsimoClient) GetCityCode(ctx context.Context, inputData []GetGeoVelocity) ([]int, error) {
	var CountryCode string
	var RegionCode string
	var CityCodeList []int
	for _, inputval := range inputData {
		var CityCode int
		if inputval.CountryName != "" {
			req, err := prosimoClient.api_client.NewRequest("GET", GetCityCode, nil)
			if err != nil {
				return nil, err
			}

			geoVelocityList := &GetGeoVelocityListRead{}
			_, err = prosimoClient.api_client.Do(ctx, req, geoVelocityList)
			if err != nil {
				return nil, err
			}
			for _, existingval := range geoVelocityList.GetGeoVelocities {
				if inputval.CountryName == existingval.CountryName {
					if inputval.StateName != "" {
						CountryCode = existingval.CountryCodeis02
					} else {
						CityCode = existingval.CityCode
						log.Println("CityCode_new", CityCode)
						CityCodeList = append(CityCodeList, CityCode)
					}

				}

			}
			if CountryCode != "" {
				UpdateGetCityCode := fmt.Sprintf("%s/%s", GetCityCode, CountryCode)
				UpdateGetCityCode1 := fmt.Sprintf("%s/%s", UpdateGetCityCode, "region")
				req1, err := prosimoClient.api_client.NewRequest("GET", UpdateGetCityCode1, nil)
				if err != nil {
					return nil, err
				}
				geoVelocityList1 := &GetGeoVelocityListRead{}
				_, err = prosimoClient.api_client.Do(ctx, req1, geoVelocityList1)
				if err != nil {
					return nil, err
				}
				for _, existingtval := range geoVelocityList1.GetGeoVelocities {
					if inputval.StateName == existingtval.StateName {
						if inputval.CityName != "" {
							RegionCode = existingtval.Region
						} else {
							CityCode = existingtval.CityCode
							CityCodeList = append(CityCodeList, CityCode)
						}
					}
				}
			}
			if RegionCode != "" {
				UpdateGetCityCode := fmt.Sprintf("%s/%s", GetCityCode, CountryCode)
				UpdateGetCityCode1 := fmt.Sprintf("%s/%s", UpdateGetCityCode, "region")
				UpdateGetCityCode2 := fmt.Sprintf("%s/%s", UpdateGetCityCode1, RegionCode)
				UpdateGetCityCode3 := fmt.Sprintf("%s/%s", UpdateGetCityCode2, "city")
				req2, err := prosimoClient.api_client.NewRequest("GET", UpdateGetCityCode3, nil)
				if err != nil {
					return nil, err
				}
				geoVelocityList2 := &GetGeoVelocityListRead{}
				_, err = prosimoClient.api_client.Do(ctx, req2, geoVelocityList2)
				if err != nil {
					return nil, err
				}
				for _, existingtval := range geoVelocityList2.GetGeoVelocities {
					if inputval.CityName == existingtval.CityName {
						CityCode = existingtval.CityCode
						CityCodeList = append(CityCodeList, CityCode)
					}
				}
			}

		} else {
			fmt.Println("Invalid Input, Country name required")
		}
	}
	return CityCodeList, nil
}
