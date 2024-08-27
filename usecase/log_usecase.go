package usecase

import (
	"context"
	"loan-tracker/domain"
)

type logUsecase struct {
	logRepo domain.LogRepository
}

func NewLogUsecase(logRepo domain.LogRepository) domain.LogUsecase {
	return &logUsecase{
		logRepo: logRepo,
	}
}

func (uc *logUsecase) LogEvent(ctx context.Context, log domain.Log) error {
	return uc.logRepo.CreateLog(ctx, log)
}

func (uc *logUsecase) GetSystemLogs(ctx context.Context) ([]domain.Log, error) {
	return uc.logRepo.GetAllLogs(ctx)
}
