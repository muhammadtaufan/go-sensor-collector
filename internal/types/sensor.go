package types

import "time"

type SensorData struct {
	ID          string
	SensorValue float64
	SensorType  string
	ID1         string
	ID2         int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
