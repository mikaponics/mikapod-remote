package internal // github.com/mikaponics/mikapod-remote/internal

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

// The time-series data structure used to store all the data that will be
// returned by the `Mikapod Soil` Arduino device.
type TimeSeriesData struct {
    HumidityValue float32 `json:"humidity_value,omitempty"`
    HumidityUnit string `json:"humidity_unit,omitempty"`
    TemperatureValue float32 `json:"temperature_primary_value,omitempty"`
    TemperatureUnit string `json:"temperature_primary_unit,omitempty"`
    PressureValue float32 `json:"pressure_value,omitempty"`
    PressureUnit string `json:"pressure_unit,omitempty"`
    TemperatureBackupValue float32 `json:"temperature_secondary_value,omitempty"`
    TemperatureBackupUnit string `json:"temperature_secondary_unit,omitempty"`
    AltitudeValue float32 `json:"altitude_value,omitempty"`
    AltitudeUnit string `json:"altitude_unit,omitempty"`
    IlluminanceValue float32 `json:"illuminance_value,omitempty"`
    IlluminanceUnit string `json:"illuminance_unit,omitempty"`
    SoilMoistureValue float32 `json:"soil_moisture_value,omitempty"`
    SoilMoistureUnit string `json:"soil_moisture_unit,omitempty"`
    Timestamp *timestamp.Timestamp `json:"timestamp,omitempty"`
}

// The time-series datum structure used to store the instrument datum of a single
// reading.
type TimeSeriesDatum struct {
	Id int64
    Instrument int32
    Value float64
    Timestamp *timestamp.Timestamp
}
