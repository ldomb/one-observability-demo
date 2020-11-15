package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

// links endpoints to transport
type Service interface {
	HealthCheck(ctx context.Context) (string, error)
	CompleteAdoption(ctx context.Context, petId, petType string) (Adoption, error)
	CleanupAdoptions(ctx context.Context) (string, error)
}

// object that handles the logic and complies with interface
type service struct {
	logger     log.Logger
	repository Repository
}

//inject dependencies into core logic
func NewService(logger log.Logger, rep Repository) Service {
	return &service{
		logger:     logger,
		repository: rep,
	}
}

// health check logic
func (s service) HealthCheck(ctx context.Context) (string, error) {
	return "alive", nil
}

// /api/completeadoption logic
func (s service) CompleteAdoption(ctx context.Context, petId, petType string) (Adoption, error) {
	logger := log.With(s.logger, "method", "CompleteAdoption")

	uuid, _ := uuid.NewV4()
	a := Adoption{
		TransactionID: uuid.String(),
		PetID:         petId,
		PetType:       petType,
		AdoptionDate:  time.Now(),
	}

	logger.Log(
		"traceId", "xray Trace to retrieve",
		"transaction", fmt.Sprintf("%#v", a),
	)

	if err := s.repository.CreateTransaction(ctx, a); err != nil {
		level.Error(logger).Log("err", err)
		return Adoption{}, err
	}

	return a, s.repository.UpdateAvailability(ctx, a)
}

// /api/completeadoption logic
func (s service) CleanupAdoptions(ctx context.Context) (string, error) {
	logger := log.With(s.logger, "method", "CleanupAdoptions")

	logger.Log(
		"traceId", "xray Trace to retrieve",
	)

	if err := s.repository.DropTransactions(ctx); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return "", nil
}
