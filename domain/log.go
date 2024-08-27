package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Timestamp time.Time          `json:"timestamp"`
	Type      string             `json:"type"` // e.g., "login_attempt", "loan_submission", etc.
	Details   string             `json:"details"`
}

type LogRepository interface {
	CreateLog(ctx context.Context, log Log) error
	GetAllLogs(ctx context.Context) ([]Log, error)
}

type LogUsecase interface {
	LogEvent(ctx context.Context, log Log) error
	GetSystemLogs(ctx context.Context) ([]Log, error)
}
