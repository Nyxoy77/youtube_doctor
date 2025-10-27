package service

import (
	"context"
	"fmt"

	"github.com/nyxoy77/B2C_YouTube_Doctor/client"
	"github.com/nyxoy77/B2C_YouTube_Doctor/models"
	"github.com/sirupsen/logrus"
)

type DoctorServiceInterface interface {
	TriggerService(ctx context.Context, channelName string) (*models.Response, error)
}

type DoctorService struct {
	doctorClient client.DoctorClientInterface
	logger       *logrus.Logger
}

func NewDoctorService(doctorClient client.DoctorClientInterface, logger *logrus.Logger) DoctorServiceInterface {
	return &DoctorService{
		doctorClient: doctorClient,
		logger:       logger,
	}
}

func (s *DoctorService) TriggerService(ctx context.Context, channelName string) (*models.Response, error) {
	fetchTitles, err := s.doctorClient.FetchTitles(ctx, channelName)
	if err != nil {
		return nil, err
	}
	if len(fetchTitles) == 0 {
		return nil, fmt.Errorf("no titles found ")
	}
	finalResponse, err := s.doctorClient.LlmCall(ctx, fetchTitles)
	if err != nil {
		return nil, err
	}
	return finalResponse, nil
}
