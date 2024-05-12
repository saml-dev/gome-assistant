package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type Scene struct {
	service Service
}

func NewScene(service Service) *Scene {
	return &Scene{
		service: service,
	}
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "scene", "apply",
		serviceData, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Create a scene entity.
func (s Scene) Create(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "scene", "create",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Reload the scenes.
func (s Scene) Reload() (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "scene", "reload", nil, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// TurnOn a scene entity.
func (s Scene) TurnOn(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "scene", "turn_on",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
