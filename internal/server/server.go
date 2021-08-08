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
	ts := db.TankStatistic{
		TestDate: req.GetTankStatistic().GetTestDate(),
	}

	if req.GetTankStatistic().GetOptionalPh() != nil {
		ph := req.GetTankStatistic().GetPh()
		ts.PH = &ph
	}

	if req.GetTankStatistic().GetOptionalGh() != nil {
		gh := req.GetTankStatistic().GetGh()
		ts.GH = &gh
	}

	if req.GetTankStatistic().GetOptionalKh() != nil {
		kh := req.GetTankStatistic().GetKh()
		ts.KH = &kh
	}

	if req.GetTankStatistic().GetOptionalAmmonia() != nil {
		a := req.GetTankStatistic().GetAmmonia()
		ts.Ammonia = &a
	}

	if req.GetTankStatistic().GetOptionalNitrite() != nil {
		n := req.GetTankStatistic().GetNitrite()
		ts.Nitrite = &n
	}

	if req.GetTankStatistic().GetOptionalNitrate() != nil {
		n := req.GetTankStatistic().GetNitrate()
		ts.Nitrate = &n
	}

	if req.GetTankStatistic().GetOptionalPhosphate() != nil {
		p := req.GetTankStatistic().GetPhosphate()
		ts.Phosphate = &p
	}

	rsp, err := s.dbManager.InsertTankStatistic(ctx, ts)
	if err != nil {
		return nil, err
	}

	tstat := &trackmyfishv1alpha1.TankStatistic{
		Id:       rsp.ID,
		TestDate: rsp.TestDate,
	}

	if rsp.PH != nil {
		tstat.OptionalPh = &trackmyfishv1alpha1.TankStatistic_Ph{Ph: *rsp.PH}
	}

	if rsp.GH != nil {
		tstat.OptionalGh = &trackmyfishv1alpha1.TankStatistic_Gh{Gh: *rsp.GH}
	}

	if rsp.KH != nil {
		tstat.OptionalKh = &trackmyfishv1alpha1.TankStatistic_Kh{Kh: *rsp.KH}
	}

	if rsp.Ammonia != nil {
		tstat.OptionalAmmonia = &trackmyfishv1alpha1.TankStatistic_Ammonia{Ammonia: *rsp.Ammonia}
	}

	if rsp.Nitrite != nil {
		tstat.OptionalNitrite = &trackmyfishv1alpha1.TankStatistic_Nitrite{Nitrite: *rsp.Nitrite}
	}

	if rsp.Nitrate != nil {
		tstat.OptionalNitrate = &trackmyfishv1alpha1.TankStatistic_Nitrate{Nitrate: *rsp.Nitrate}
	}

	if rsp.Phosphate != nil {
		tstat.OptionalPhosphate = &trackmyfishv1alpha1.TankStatistic_Phosphate{Phosphate: *rsp.Phosphate}
	}

	return &trackmyfishv1alpha1.AddTankStatisticResponse{TankStatistic: tstat}, nil
}

func (s *Server) ListTankStatistics(ctx context.Context, req *trackmyfishv1alpha1.ListTankStatisticsRequest) (*trackmyfishv1alpha1.ListTankStatisticsResponse, error) {
	rsp, err := s.dbManager.GetTankStatistics(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get tank statistics")
	}

	tankStats := make([]*trackmyfishv1alpha1.TankStatistic, len(rsp))
	for i, ts := range rsp {
		tstat := &trackmyfishv1alpha1.TankStatistic{
			Id:       ts.ID,
			TestDate: ts.TestDate,
		}

		if ts.PH != nil {
			tstat.OptionalPh = &trackmyfishv1alpha1.TankStatistic_Ph{Ph: *ts.PH}
		}

		if ts.GH != nil {
			tstat.OptionalGh = &trackmyfishv1alpha1.TankStatistic_Gh{Gh: *ts.GH}
		}

		if ts.KH != nil {
			tstat.OptionalKh = &trackmyfishv1alpha1.TankStatistic_Kh{Kh: *ts.KH}
		}

		if ts.Ammonia != nil {
			tstat.OptionalAmmonia = &trackmyfishv1alpha1.TankStatistic_Ammonia{Ammonia: *ts.Ammonia}
		}

		if ts.Nitrite != nil {
			tstat.OptionalNitrite = &trackmyfishv1alpha1.TankStatistic_Nitrite{Nitrite: *ts.Nitrite}
		}

		if ts.Nitrate != nil {
			tstat.OptionalNitrate = &trackmyfishv1alpha1.TankStatistic_Nitrate{Nitrate: *ts.Nitrate}
		}

		if ts.Phosphate != nil {
			tstat.OptionalPhosphate = &trackmyfishv1alpha1.TankStatistic_Phosphate{Phosphate: *ts.Phosphate}
		}

		tankStats[i] = tstat
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

	tstat := &trackmyfishv1alpha1.TankStatistic{
		Id:       rsp.ID,
		TestDate: rsp.TestDate,
	}

	if rsp.PH != nil {
		tstat.OptionalPh = &trackmyfishv1alpha1.TankStatistic_Ph{Ph: *rsp.PH}
	}

	if rsp.GH != nil {
		tstat.OptionalGh = &trackmyfishv1alpha1.TankStatistic_Gh{Gh: *rsp.GH}
	}

	if rsp.KH != nil {
		tstat.OptionalKh = &trackmyfishv1alpha1.TankStatistic_Kh{Kh: *rsp.KH}
	}

	if rsp.Ammonia != nil {
		tstat.OptionalAmmonia = &trackmyfishv1alpha1.TankStatistic_Ammonia{Ammonia: *rsp.Ammonia}
	}

	if rsp.Nitrite != nil {
		tstat.OptionalNitrite = &trackmyfishv1alpha1.TankStatistic_Nitrite{Nitrite: *rsp.Nitrite}
	}

	if rsp.Nitrate != nil {
		tstat.OptionalNitrate = &trackmyfishv1alpha1.TankStatistic_Nitrate{Nitrate: *rsp.Nitrate}
	}

	if rsp.Phosphate != nil {
		tstat.OptionalPhosphate = &trackmyfishv1alpha1.TankStatistic_Phosphate{Phosphate: *rsp.Phosphate}
	}

	return &trackmyfishv1alpha1.DeleteTankStatisticResponse{
		TankStatistic: tstat,
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
