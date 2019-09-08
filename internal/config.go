package internal

import (
    "os"
    "strconv"
)

func GetTenantId() int64 {
    tenantId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_TENANT_ID"), 10, 64)
    return tenantId
}

func GetSensorIdFromInstrumentId(instrumentId int32) int64 {
    switch instrumentId {
    case 1:
        sensorId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_HUMIDITY_SENSOR_ID"), 10, 64)
        return sensorId
    case 2:
        sensorId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_TEMPERATURE_SENSOR_ID"), 10, 64)
        return sensorId
    case 3:
        sensorId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_PRESSURE_SENSOR_ID"), 10, 64)
        return sensorId
    case 4:
        sensorId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_TEMPERATURE_BACKEND_SENSOR_ID"), 10, 64)
        return sensorId
    case 5:
        sensorId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_ALTITUDE_SENSOR_ID"), 10, 64)
        return sensorId
    case 6:
        sensorId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_ILLUMINANCE_SENSOR_ID"), 10, 64)
        return sensorId
    case 7:
        sensorId, _ := strconv.ParseInt(os.Getenv("MIKAPOD_SOIL_MOISTURE_SENSOR_ID"), 10, 64)
        return sensorId
    default:
        return 0
    }
}
