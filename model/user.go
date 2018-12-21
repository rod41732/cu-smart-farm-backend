package model

// User : interface type of 'RealUser' and 'NullUser'
type User interface {
	ReportStatus(payload interface{})
}
