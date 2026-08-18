package gomeassistant_test

import (
	"context"
	"testing"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/services"
)

func TestServiceAPIUsesPubliclyNameableTypes(t *testing.T) {
	request := services.BaseServiceRequest{
		Domain:  "light",
		Service: "turn_on",
		Target:  services.Entities("light.kitchen"),
	}

	var target ga.Target = request.Target
	var light *services.Light = (&ga.Service{}).Light
	var call func(*ga.App, context.Context, services.BaseServiceRequest, any) error = (*ga.App).Call
	var callAndForget func(*ga.App, services.BaseServiceRequest) error = (*ga.App).CallAndForget

	if len(target.EntityIDs) != 1 || target.EntityIDs[0] != "light.kitchen" {
		t.Fatalf("unexpected target: %#v", target)
	}
	if light != nil {
		t.Fatal("zero Service unexpectedly contains a Light service")
	}
	if call == nil || callAndForget == nil {
		t.Fatal("service call methods are unavailable")
	}
}
