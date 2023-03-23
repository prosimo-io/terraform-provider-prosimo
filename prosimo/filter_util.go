package prosimo

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func matchOperand(filtered map[string]string, value string) bool {

	operand := filtered["oprand"]
	values := filtered["values"]
	var result bool

	if operand == client.PatternEqual {
		if values == value {
			result = true
		}

	}
	if operand == client.PatternEqualNoCase {
		if strings.EqualFold(values, value) {
			result = true
		}

	}
	if operand == client.PatternNotEqual {
		if !(strings.EqualFold(values, value)) {
			result = true
		}
	}

	if operand == client.PatternContains {
		if strings.Contains(value, values) {
			result = true
		}
	}
	return result
}

func lowerCaseConvertMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range input {
		switch v := v.(type) {
		case map[string]interface{}:
			result[strings.ToLower(k)] = lowerCaseConvertMap(v)
		case []interface{}:
			for i, iv := range v {
				if reflect.TypeOf(iv).Kind() == reflect.Map {
					v[i] = lowerCaseConvertMap(iv.(map[string]interface{}))
				}
			}
			result[strings.ToLower(k)] = v
		default:
			result[strings.ToLower(k)] = v
		}
	}
	return result
}

func checkMainOperand(filter string, value map[string]interface{}) (diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	// Conveting the value map (keys) to lowercase form camelCase
	lowercaseValueMap := lowerCaseConvertMap(value)
	log.Println("lowercaseValueMap", lowercaseValueMap)

	var strArr []string
	var ifAnd bool

	if (strings.Contains(filter, client.PatternAND)) && (strings.Contains(filter, client.PatternOR)) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Input, either use '&' or ',' in a filter",
			Detail:   fmt.Sprintln("Invalid Input, please use AND(&) / OR(,) in the filter field. Do not use multiple or same operand twice"),
		})
		return diags, false

	} else if strings.Count(filter, client.PatternAND) > 1 || strings.Count(filter, client.PatternOR) > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Input, either use '&' or ',' once in a filter",
			Detail:   fmt.Sprintln("Invalid Input, please use AND(&) / OR(,) in the filter field only once"),
		})

		return diags, false

	} else {
		if strings.Contains(filter, client.PatternOR) {
			strArr = strings.Split(filter, client.PatternOR)
		}
		if strings.Contains(filter, client.PatternAND) {
			strArr = strings.Split(filter, client.PatternAND)
			ifAnd = true
		}
		if len(strArr) == 0 {
			strArr = append(strArr, filter)
		}
		log.Println("StrArr :", strArr)
		for _, v := range strArr {
			diags, filtered := checkFilter(v)
			if diags != nil {
				return diags, false
			}
			tempKey := filtered["key"]
			var res bool
			if v, ok := lowercaseValueMap[tempKey]; ok {
				if data, ok := v.(string); ok {
					res = matchOperand(filtered, data)
				} else {
					return nil, false
				}
			}
			// res := matchOperand(a, value)
			if ifAnd && !res {
				return nil, false
			}
			if !ifAnd && res {
				return nil, true
			}
		}
		if ifAnd {
			return nil, true
		}
	}
	return nil, false
}

func checkFilter(filter string) (diag.Diagnostics, map[string]string) {
	filtered := make(map[string]string)
	var diags diag.Diagnostics
	if strings.Contains(filter, client.PatternEqual) {
		equalSlice := strings.Split(filter, client.PatternEqual)
		filtered["key"] = equalSlice[0]
		filtered["oprand"] = client.PatternEqual
		filtered["values"] = equalSlice[1]

	} else if strings.Contains(filter, client.PatternEqualNoCase) {
		equalnocaseSlice := strings.Split(filter, client.PatternEqualNoCase)
		filtered["key"] = equalnocaseSlice[0]
		filtered["oprand"] = client.PatternEqualNoCase
		filtered["values"] = equalnocaseSlice[1]

	} else if strings.Contains(filter, client.PatternContains) {
		containsSlice := strings.Split(filter, client.PatternContains)
		filtered["key"] = containsSlice[0]
		filtered["oprand"] = client.PatternContains
		filtered["values"] = containsSlice[1]

	} else if strings.Contains(filter, client.PatternNotEqual) {
		noteqSlice := strings.Split(filter, client.PatternNotEqual)
		filtered["key"] = noteqSlice[0]
		filtered["oprand"] = client.PatternNotEqual
		filtered["values"] = noteqSlice[1]
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Input, Conditional Operator should be one of ==, !=, =*, =@",
			Detail:   fmt.Sprintln("Invalid Input, Conditional Operator should be one of ==, !=, =*, =@"),
		})
		return diags, nil
	}
	return nil, filtered
}
