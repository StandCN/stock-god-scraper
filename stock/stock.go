package stock

import "time"

type SourceData interface {
	GetText() string
	GetDate() time.Time
	GetUrl() string
}
