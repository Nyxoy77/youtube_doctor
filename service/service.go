package service

import (
	"context"

	"github.com/nyxoy77/B2C_YouTube_Doctor/client"
	"github.com/nyxoy77/B2C_YouTube_Doctor/models"
)

type DoctorServiceInterface interface {
	TriggerService(ctx context.Context, channelName string) (*models.Response, error)
}

type DoctorService struct {
	doctorClient client.DoctorClientInterface
}

func NewDoctorService(doctorClient client.DoctorClientInterface) DoctorServiceInterface {
	return &DoctorService{
		doctorClient: doctorClient,
	}
}

func (s *DoctorService) TriggerService(ctx context.Context, channelName string) (*models.Response, error) {
	fetchTitles, err := s.doctorClient.FetchTitles(ctx, channelName)
	if err != nil {
		return nil, err
	}
	finalResponse, err := s.doctorClient.LlmCall(ctx, fetchTitles)
	if err != nil {
		return nil, err
	}
	return finalResponse, nil
}
