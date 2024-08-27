package repositories

import (
	"context"
	"loan-tracker/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type logRepository struct {
	db *mongo.Collection
}

func NewLogRepository(db *mongo.Client) domain.LogRepository {
	return &logRepository{
		db: db.Database("loan-tracker").Collection("logs"),
	}
}

func (r *logRepository) CreateLog(ctx context.Context, log domain.Log) error {
	_, err := r.db.InsertOne(ctx, log)
	return err
}

func (r *logRepository) GetAllLogs(ctx context.Context) ([]domain.Log, error) {
	var logs []domain.Log
	cursor, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &logs)
	return logs, err
}
