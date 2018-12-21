package model

type RealUser struct {
	currentToken string
	devices      []string
}

type User interface {
	Command(relay string, workmode string, payload interface{})
	ReportStatus(payload interface{})
	AddDevice(sensorID string, sensorInfo interface{})
}
