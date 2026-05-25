package service

import (
	"time"
)

type DeliveryTimeFactory interface {
	CalculateDeadline(transportType string) time.Time
}

type DeliveryFactory struct{}

func NewDeliveryFactory() *DeliveryFactory {
	return &DeliveryFactory{}
}

func (f *DeliveryFactory) CalculateDeadline(transportType string) time.Time {
	now := time.Now()

	switch transportType {
	case "on_foot":
		return now.Add(30 * time.Minute)
	case "scooter":
		return now.Add(15 * time.Minute)
	case "car":
		return now.Add(5 * time.Minute)
	}
	return now.Add(30 * time.Minute)
}
