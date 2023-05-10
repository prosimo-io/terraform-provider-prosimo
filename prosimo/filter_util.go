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
	fmt.Println("VLAUE", value)

	operand := filtered["oprand"]
	values := filtered["values"]
	log.Println("Passed value to match", value)
	result := false

	if operand == client.PatternEqual {
		if strings.ToLower(values) == value {
			result = true
		}

	}
	if operand == client.PatternEqualNoCase {
		if strings.EqualFold(strings.ToLower(values), value) {
			result = true
		}

	}
	if operand == client.PatternNotEqual {
		if !(strings.EqualFold(strings.ToLower(values), value)) {
			result = true
		}
	}

	if operand == client.PatternContains {
		if strings.Contains(value, strings.ToLower(values)) {
			result = true
		}
	}
	return result
}

func traverseStruct(path string, filtered map[string]string, val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			sliceVal := val.Index(i)
			fieldPath := fmt.Sprintf("%s[%d]", path, i)
			if traverseStruct(fieldPath, filtered, sliceVal) {
				return true
			}
		}

	case reflect.Ptr:
		if val.IsNil() {
			fmt.Printf("%s: <nil>\n", path)
		} else {
			elem := val.Elem()
			fieldPath := path + "*"
			if traverseStruct(fieldPath, filtered, elem) {
				return true
			}
		}

	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			dataValue := val.Field(i)
			key := val.Type().Field(i).Name

			result := matchKV(filtered, key, dataValue)
			if result {
				return true
			}
			fieldPath := path + "." + key

			if traverseStruct(fieldPath, filtered, dataValue) {
				return true
			}
		}

	case reflect.String:
		fmt.Printf("%s: %s\n", path, val.String())

	case reflect.Int:
		fmt.Printf("%s: %d\n", path, val.Int())

	case reflect.Bool:
		fmt.Printf("%s: %t\n", path, val.Bool())

	default:
		fmt.Printf("unhandled data type at this location %s: %v\n", path, val)

	}
	return false
}

func matchKV(filtered map[string]string, key string, value reflect.Value) bool {

	if strings.ToLower(key) != strings.ToLower(filtered["key"]) {

		return false
	}
	log.Println("Key Matched")

	return matchOperand(filtered, strings.ToLower(value.String()))

}

func checkMainOperand(filter string, value reflect.Value) (diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	log.Println("Value Map:", value)

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
			filtered := checkFilter(v)
			fmt.Println("filtered", filtered)
			var res bool
			res = traverseStruct("root", filtered, value) //pass filtered
			log.Println("final result:", filtered["key"], res)
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

func checkFilter(filter string) map[string]string {
	filtered := make(map[string]string)

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

	} else {
		noteqSlice := strings.Split(filter, client.PatternNotEqual)
		filtered["key"] = noteqSlice[0]
		filtered["oprand"] = client.PatternNotEqual
		filtered["values"] = noteqSlice[1]
	}
	return filtered
}
