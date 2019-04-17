package provider

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestQuantity(t *testing.T) {
	testCases := []map[string]interface{}{
		// <decimalSI> ::= m | "" | k | M | G | T | P | E
		{
			"in": "",		// error
		},
		{
			"in": "0.0",
		},
		{
			"in": "0.1",	// 100m
		},
		{
			"in": "1",		// 1000m
		},
		{
			"in": "0.0m",
		},
		{
			"in": "0.1m",	// 100u
		},
		{
			"in": "1m",
		},
		{
			"in": "0.1M",	// 100k
		},
		{
			"in": "1M",		// 1000 * 1000 = 10^6
		},
		// <binarySI> ::= Ki | Mi | Gi | Ti | Pi | Ei
		{
			"in": "1023Ki",	// 1204 * 1024 = 1048576B
		},
		{
			"in": "1Mi",	// 1204 * 1024 = 1048576B
		},
		{
			"in": "1.1Mi",	// 1.1 * 1204 * 1024
		},
	}
	for _, testCase := range testCases {
		q, err := resource.ParseQuantity(testCase["in"].(string))
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Printf("quantity: %+v, quantity string: %s\n", q, q.String())
	}
}
