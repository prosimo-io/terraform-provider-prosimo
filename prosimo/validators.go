package prosimo

import (
	"fmt"
)

func validateCloudType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "AWS" && value != "AZURE" && value != "GCP" {
		errors = append(errors, fmt.Errorf(
			"%q must be one of \"AWS\", \"AZURE\", or \"GCP\"", k))
	}
	return
}
