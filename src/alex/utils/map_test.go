package utils

import (
	"testing"
	"net/url"
)

func TestMapSort(t *testing.T) {
	m := map[string]interface{}{"d":1, "e":2, "b":4}
	t.Logf("TestMapSort: %v", MapSort(m))
}

func TestMapToQuery(t *testing.T) {
	m := map[string]interface{}{"d":1, "e":2, "b":4}
	t.Logf("TestMapToQuery Sorted: %v", MapToQuery(m, true))
	t.Logf("TestMapToQuery UnSorted: %v", MapToQuery(m, false))
}

func TestMapFromQuery(t *testing.T) {
	query := "w=1&b=2&a=3"
	m := MapFromQuery(query)
	t.Logf("TestMapFromQuery %v ,%v", query, m)
}

func TestValuesSort(t *testing.T) {
	instance, _ := url.Parse("http://www.a.com?d=1&e=2&b=4")
	t.Logf("TestValuesSort %v", ValuesSort(instance.Query()))
}

func TestValuesToQuery(t *testing.T) {
	instance, _ := url.Parse("http://www.a.com?d=1&e=2&b=4")
	t.Logf("TestValuesToQuery Sorted: %v", ValuesToQuery(instance.Query(), true))
	t.Logf("TestValuesToQuery UnSorted: %v", ValuesToQuery(instance.Query(), false))
}

func TestValuesFromQuery(t *testing.T) {
	t.Logf("TestValuesFromQuery %v", ValuesFromQuery("w=1&b=2&a=3"))
}
