package model

type NullUser struct {
}

func (user *NullUser) Command(relay string, workmode string, payload interface{}) {

}
func (user *NullUser) ReportStatus(payload interface{}) {

}

func (user *NullUser) AddDevice(sensorID string, sensorInfo interface{}) {

}
