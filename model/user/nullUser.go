package user

// NullUser is "placeholder" when client is disconnected
type NullUser struct {
}

// ReportStatus for null user do nothing
func (user *NullUser) ReportStatus(payload interface{}) {
	return
}
