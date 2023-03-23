package client

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateGeoVelocity(t *testing.T) {

	ctx := context.Background()
	var cityCodeList []GetGeoVelocity
	//var addLocation GetGeoVelocity
	deletelocations := []GetGeoVelocity{}
	addlocations := []GetGeoVelocity{}

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getCityCode2 := GetGeoVelocity{
		//CityName:    "Sihala",
		StateName:   "Keryneia",
		CountryName: "Cyprus",
	}
	getCityCode3 := GetGeoVelocity{
		//CityName:    "Sihala",x`
		//StateName:   "Odisha",
		CountryName: "Philippines",
	}
	cityCodeList = append(cityCodeList, getCityCode2, getCityCode3)

	respString, err := prosimo_client.GetCityCode(ctx, cityCodeList)
	for i, val := range respString {
		_ = val
		deleteLocation := GetGeoVelocity{
			CityCode: respString[i],
		}
		deletelocations = append(deletelocations, deleteLocation)
	}

	fmt.Printf("respString2 - %d", respString)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

	postgeovelocity := &PostGeoVelocity{
		AddLocations:    addlocations,
		DeleteLocations: deletelocations,
	}
	fmt.Println(postgeovelocity)
	respString1, err := prosimo_client.UpdateGeoVelocity(ctx, postgeovelocity)
	_ = respString1
}

func TestGetGeoVelocity(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getgeovelocty, err := prosimo_client.GetGeoVelocity(ctx)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

	fmt.Printf("respString2 %v\n", getgeovelocty.GetGeoVelocities)

}
