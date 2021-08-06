package server

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/trackmyfish/backend/internal/db"
	"github.com/trackmyfish/backend/internal/fishbase"
	trackmyfishv1alpha1 "github.com/trackmyfish/proto/trackmyfish/v1alpha1"
)

// Server is the implementation of the trackmyfishv1alpha1.TrackMyFishServiceServer
type Server struct {
	dbManager      *db.Manager
	fishbaseClient *fishbase.Client
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

	fc := fishbase.New(http.DefaultClient)

	return &Server{
		dbManager:      dbManager,
		fishbaseClient: fc,
	}, nil
}

func (s *Server) Heartbeat(ctx context.Context, req *trackmyfishv1alpha1.HeartbeatRequest) (*trackmyfishv1alpha1.HeartbeatResponse, error) {
	d, err := s.fishbaseClient.GetHeartbeat()
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
		Count:              req.GetFish().GetCount(),
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
			Count:             rsp.Count,
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
			Count:             fish.Count,
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
			Count:             rsp.Count,
		},
	}, nil
}

func (s *Server) AddTankStatistic(ctx context.Context, req *trackmyfishv1alpha1.AddTankStatisticRequest) (*trackmyfishv1alpha1.AddTankStatisticResponse, error) {
	rsp, err := s.dbManager.InsertTankStatistic(ctx, db.TankStatistic{
		TestDate:  req.GetTankStatistic().GetTestDate(),
		PH:        req.GetTankStatistic().GetPH(),
		GH:        req.GetTankStatistic().GetGH(),
		KH:        req.GetTankStatistic().GetKH(),
		Ammonia:   req.GetTankStatistic().GetAmmonia(),
		Nitrite:   req.GetTankStatistic().GetNitrite(),
		Nitrate:   req.GetTankStatistic().GetNitrate(),
		Phosphate: req.GetTankStatistic().GetPhosphate(),
	})
	if err != nil {
		return nil, err
	}

	return &trackmyfishv1alpha1.AddTankStatisticResponse{
		TankStatistic: &trackmyfishv1alpha1.TankStatistic{
			Id:        rsp.ID,
			TestDate:  rsp.TestDate,
			PH:        rsp.PH,
			GH:        rsp.GH,
			KH:        rsp.KH,
			Ammonia:   rsp.Ammonia,
			Nitrite:   rsp.Nitrite,
			Nitrate:   rsp.Nitrate,
			Phosphate: rsp.Phosphate,
		},
	}, nil
}

func (s *Server) ListTankStatistics(ctx context.Context, req *trackmyfishv1alpha1.ListTankStatisticsRequest) (*trackmyfishv1alpha1.ListTankStatisticsResponse, error) {
	rsp, err := s.dbManager.GetTankStatistics(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get tank statistics")
	}

	tankStats := make([]*trackmyfishv1alpha1.TankStatistic, len(rsp))
	for i, ts := range rsp {
		tankStats[i] = &trackmyfishv1alpha1.TankStatistic{
			Id:        ts.ID,
			TestDate:  ts.TestDate,
			PH:        ts.PH,
			GH:        ts.GH,
			KH:        ts.KH,
			Ammonia:   ts.Ammonia,
			Nitrite:   ts.Nitrite,
			Nitrate:   ts.Nitrate,
			Phosphate: ts.Phosphate,
		}
	}

	return &trackmyfishv1alpha1.ListTankStatisticsResponse{
		TankStatistics: tankStats,
	}, nil
}

func (s *Server) DeleteTankStatistic(ctx context.Context, req *trackmyfishv1alpha1.DeleteTankStatisticRequest) (*trackmyfishv1alpha1.DeleteTankStatisticResponse, error) {
	rsp, err := s.dbManager.DeleteTankStatistic(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &trackmyfishv1alpha1.DeleteTankStatisticResponse{
		TankStatistic: &trackmyfishv1alpha1.TankStatistic{
			Id:        rsp.ID,
			TestDate:  rsp.TestDate,
			PH:        rsp.PH,
			GH:        rsp.GH,
			KH:        rsp.KH,
			Ammonia:   rsp.Ammonia,
			Nitrite:   rsp.Nitrite,
			Nitrate:   rsp.Nitrate,
			Phosphate: rsp.Phosphate,
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
