package models

import "time"

type Transaction struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	KIND        string
	AmountCents int64
	Category    string
	OccurredOn  time.Time `gorm:"type:date"`
	Note        string
	UpdatedAt   time.Time
}
