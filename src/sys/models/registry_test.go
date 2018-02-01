package models

import (
	"testing"
	"fmt"
)

const (
	REGISTRY_DOMAIN = "service-registry.fuyoukache.com"
	REGISTRY_USERNAME = "linyang"
	REGISTRY_PASSWORD = "123456"
)

func TestRegistry_GetImages(t *testing.T) {
	r := NewRegistry(REGISTRY_USERNAME, REGISTRY_PASSWORD, REGISTRY_DOMAIN)
	res, err := r.GetImages()
	if err != nil {
		t.Errorf("TestRegistry_GetImages Error %v", err.Error())
	} else {
		fmt.Println(res)
	}
}

func TestRegistry_GetImageVersions(t *testing.T) {
	r := NewRegistry(REGISTRY_USERNAME, REGISTRY_PASSWORD, REGISTRY_DOMAIN)
	res, err := r.GetImageVersions("yaf")
	if err != nil {
		t.Errorf("TestRegistry_GetImageVersions Error %v", err.Error())
	} else {
		fmt.Println(res)
	}
}

func TestRegistry_HasImage(t *testing.T) {
	r := NewRegistry(REGISTRY_USERNAME, REGISTRY_PASSWORD, REGISTRY_DOMAIN)
	res, err := r.HasImage("a", "")
	if err != nil {
		t.Errorf("TestRegistry_HasImage Error %v", err.Error())
	} else {
		fmt.Println(res)
	}
}

func TestRegistry_GetAllImages(t *testing.T) {
	r := NewRegistry(REGISTRY_USERNAME, REGISTRY_PASSWORD, REGISTRY_DOMAIN)
	res, err := r.GetAllImages()
	if err != nil {
		t.Errorf("TestRegistry_GetAllImages Error %v", err.Error())
	} else {
		fmt.Println(res)
	}
}

func TestRegistry_DeleteImage(t *testing.T) {
	r := NewRegistry(REGISTRY_USERNAME, REGISTRY_PASSWORD, REGISTRY_DOMAIN)
	res, err := r.DeleteImage("yaf", "v2.0.0")
	if err != nil {
		t.Errorf("TestRegistry_DeleteImage Error %v", err.Error())
	} else {
		fmt.Println(res)
	}
}