package dto

import "time"

type BaseRequest[T any] struct {
	RequestID   string    `json:"request_id,omitempty"`
	Lang        string    `json:"lang,omitempty"`
	RequestTime time.Time `json:"request_time,omitempty"`
	Data        T         `json:"data"`
}
