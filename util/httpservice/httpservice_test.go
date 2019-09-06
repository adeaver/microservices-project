package httpservice

import (
	"testing"
)

func TestMakeRouteMethodString(t *testing.T) {
	getStr := RouteMethodGET.Str()
	if getStr != "GET" {
		t.Errorf("Expected GET, got %s", getStr)
	}
	postStr := RouteMethodPOST.Str()
	if postStr != "POST" {
		t.Errorf("Expected POST, got %s", postStr)
	}
}
