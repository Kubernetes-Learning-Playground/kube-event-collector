package model

import "time"

// Event 事件对象
type Event struct {
	Type      string
	Kind      string
	Name      string
	Namespace string
	Timestamp time.Time
	Message   string
	Reason    string
	Source    string
	Host      string
	Count     int32
}
