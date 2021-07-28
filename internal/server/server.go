package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/trackmyfish/backend/internal/db"
	"github.com/trackmyfish/backend/internal/fishbase"
	trackmyfishv1alpha1 "github.com/trackmyfish/proto/trackmyfish/v1alpha1"
)

// Server is the implementation of the trackmyfishv1alpha1.TrackMyFishServiceServer
type Server struct {
	dbManager *db.Manager
}

type Config struct {
	DBHost     string
	DBPort     int
	DBUsername string
	DBPassword string
	DBName     string
}

func New(c Config) (*Server, error) {
	dbManager, err := db.New(db.Config{
		Host:     c.DBHost,
		Port:     c.DBPort,
		Username: c.DBUsername,
		Password: c.DBPassword,
		Database: c.DBName,
	})
	if err != nil {
		return nil, err
	}

	return &Server{
		dbManager: dbManager,
	}, nil
}

func (s *Server) Heartbeat(ctx context.Context, req *trackmyfishv1alpha1.HeartbeatRequest) (*trackmyfishv1alpha1.HeartbeatResponse, error) {
	d, err := fishbase.GetHeartbeat()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get heartbeat information")
	}

	return &trackmyfishv1alpha1.HeartbeatResponse{
		Fishbase: &trackmyfishv1alpha1.HeartbeatStatus{
			Status: fishbaseStatusToStatus(d.Status),
		},
	}, nil
}

func (s *Server) AddFish(ctx context.Context, req *trackmyfishv1alpha1.AddFishRequest) (*trackmyfishv1alpha1.AddFishResponse, error) {
	d, err := fishbase.GetDetails(req.GetFish().GetGenus(), req.GetFish().GetSpecies())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"stack": err,
		}).Error("Unable to get fish details")
	}

	rsp, err := s.dbManager.InsertFish(ctx, db.Fish{
		Genus:              req.GetFish().GetGenus(),
		Species:            req.GetFish().GetSpecies(),
		CommonName:         req.GetFish().GetCommonName(),
		Name:               req.GetFish().GetName(),
		Color:              req.GetFish().GetColor(),
		Gender:             req.GetFish().GetGender().String(),
		PurchaseDate:       req.GetFish().GetPurchaseDate(),
		EcosystemName:      d.Ecosystem.Name,
		EcosystemType:      d.Ecosystem.Type,
		EchosystemLocation: d.Ecosystem.Location,
		Salinity:           d.Ecosystem.Salinity,
		Climate:            d.Ecosystem.Climate,
	})
	if err != nil {
		return nil, err
	}

	return &trackmyfishv1alpha1.AddFishResponse{
		Fish: &trackmyfishv1alpha1.Fish{
			Id:                rsp.ID,
			Genus:             rsp.Genus,
			Species:           rsp.Species,
			CommonName:        rsp.CommonName,
			Name:              rsp.Name,
			Color:             rsp.Color,
			Gender:            stringToGender(rsp.Gender),
			PurchaseDate:      rsp.PurchaseDate,
			EcosystemName:     rsp.EcosystemName,
			EcosystemType:     rsp.EcosystemType,
			EcosystemLocation: rsp.EchosystemLocation,
			Salinity:          rsp.Salinity,
			Climate:           rsp.Climate,
		},
	}, nil
}

func (s *Server) ListFish(ctx context.Context, req *trackmyfishv1alpha1.ListFishRequest) (*trackmyfishv1alpha1.ListFishResponse, error) {
	rsp, err := s.dbManager.GetFish(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{})
	}

	f := make([]*trackmyfishv1alpha1.Fish, len(rsp))
	for i, fish := range rsp {
		f[i] = &trackmyfishv1alpha1.Fish{
			Id:                fish.ID,
			Genus:             fish.Genus,
			Species:           fish.Species,
			CommonName:        fish.CommonName,
			Name:              fish.Name,
			Color:             fish.Color,
			Gender:            stringToGender(fish.Gender),
			PurchaseDate:      fish.PurchaseDate,
			EcosystemName:     fish.EcosystemName,
			EcosystemType:     fish.EcosystemType,
			EcosystemLocation: fish.EchosystemLocation,
			Salinity:          fish.Salinity,
			Climate:           fish.Climate,
		}
	}

	return &trackmyfishv1alpha1.ListFishResponse{
		Fish: f,
	}, nil
}

func (s *Server) DeleteFish(ctx context.Context, req *trackmyfishv1alpha1.DeleteFishRequest) (*trackmyfishv1alpha1.DeleteFishResponse, error) {
	rsp, err := s.dbManager.DeleteFish(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// TODO - change proto and return the deleted fish
	_ = rsp

	return &trackmyfishv1alpha1.DeleteFishResponse{
		Fish: &trackmyfishv1alpha1.Fish{
			Id:                rsp.ID,
			Genus:             rsp.Genus,
			Species:           rsp.Species,
			CommonName:        rsp.CommonName,
			Name:              rsp.Name,
			Color:             rsp.Color,
			Gender:            stringToGender(rsp.Gender),
			PurchaseDate:      rsp.PurchaseDate,
			EcosystemName:     rsp.EcosystemName,
			EcosystemType:     rsp.EcosystemType,
			EcosystemLocation: rsp.EchosystemLocation,
			Salinity:          rsp.Salinity,
			Climate:           rsp.Climate,
		},
	}, nil
}

func stringToGender(gender string) trackmyfishv1alpha1.Fish_Gender {
	switch gender {
	case trackmyfishv1alpha1.Fish_MALE.String():
		return trackmyfishv1alpha1.Fish_MALE
	case trackmyfishv1alpha1.Fish_FEMALE.String():
		return trackmyfishv1alpha1.Fish_FEMALE
	}

	return trackmyfishv1alpha1.Fish_UNSPECIFIED
}

func fishbaseStatusToStatus(fs fishbase.HeartbeatStatus) trackmyfishv1alpha1.HeartbeatStatus_Status {
	switch fs {
	case fishbase.HeartbeatStatusDown:
		return trackmyfishv1alpha1.HeartbeatStatus_DOWN
	case fishbase.HeartbeatStatusDegraded:
		return trackmyfishv1alpha1.HeartbeatStatus_DEGRADED
	case fishbase.HeartbeatStatusOperational:
		return trackmyfishv1alpha1.HeartbeatStatus_OPERATIONAL
	}

	return trackmyfishv1alpha1.HeartbeatStatus_UNSPECIFIED
}
