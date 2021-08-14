package server

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/trackmyfish/backend/internal/db"
	trackmyfishv1alpha1 "github.com/trackmyfish/proto/trackmyfish/v1alpha1"
)

type fishQuerier interface {
	ListFish(context.Context) ([]db.Fish, error)
}

type fishModifier interface {
	InsertFish(context.Context, db.Fish) (db.Fish, error)
	DeleteFish(context.Context, int32) (db.Fish, error)
}

type tankStatQuerier interface {
	ListTankStatistics(context.Context) ([]db.TankStatistic, error)
}

type tankStatModifier interface {
	InsertTankStatistic(context.Context, db.TankStatistic) (db.TankStatistic, error)
	DeleteTankStatistic(context.Context, int32) (db.TankStatistic, error)
}

// Server is the implementation of the trackmyfishv1alpha1.TrackMyFishServiceServer
type Server struct {
	fishQuerier      fishQuerier
	fishModifier     fishModifier
	tankStatQuerier  tankStatQuerier
	tankStatModifier tankStatModifier
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
		return nil, errors.Wrap(err, "unable to create db instance")
	}

	return &Server{
		fishQuerier:      dbManager,
		fishModifier:     dbManager,
		tankStatQuerier:  dbManager,
		tankStatModifier: dbManager,
	}, nil
}

func (s *Server) Heartbeat(ctx context.Context, req *trackmyfishv1alpha1.HeartbeatRequest) (*trackmyfishv1alpha1.HeartbeatResponse, error) {
	return &trackmyfishv1alpha1.HeartbeatResponse{}, nil
}

func (s *Server) AddFish(ctx context.Context, req *trackmyfishv1alpha1.AddFishRequest) (*trackmyfishv1alpha1.AddFishResponse, error) {
	rsp, err := s.fishModifier.InsertFish(ctx, db.Fish{
		Type:         req.GetFish().GetType(),
		Subtype:      req.GetFish().GetSubtype(),
		Color:        req.GetFish().GetColor(),
		Gender:       req.GetFish().GetGender().String(),
		PurchaseDate: req.GetFish().GetPurchaseDate(),
		Count:        req.GetFish().GetCount(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to add fish")
	}

	return &trackmyfishv1alpha1.AddFishResponse{
		Fish: &trackmyfishv1alpha1.Fish{
			Id:           rsp.ID,
			Type:         rsp.Type,
			Subtype:      rsp.Subtype,
			Color:        rsp.Color,
			Gender:       stringToGender(rsp.Gender),
			PurchaseDate: rsp.PurchaseDate,
			Count:        rsp.Count,
		},
	}, nil
}

func (s *Server) ListFish(ctx context.Context, req *trackmyfishv1alpha1.ListFishRequest) (*trackmyfishv1alpha1.ListFishResponse, error) {
	rsp, err := s.fishQuerier.ListFish(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list fish")
	}

	f := make([]*trackmyfishv1alpha1.Fish, len(rsp))
	for i, fish := range rsp {
		f[i] = &trackmyfishv1alpha1.Fish{
			Id:           fish.ID,
			Type:         fish.Type,
			Subtype:      fish.Subtype,
			Color:        fish.Color,
			Gender:       stringToGender(fish.Gender),
			PurchaseDate: fish.PurchaseDate,
			Count:        fish.Count,
		}
	}

	return &trackmyfishv1alpha1.ListFishResponse{
		Fish: f,
	}, nil
}

func (s *Server) DeleteFish(ctx context.Context, req *trackmyfishv1alpha1.DeleteFishRequest) (*trackmyfishv1alpha1.DeleteFishResponse, error) {
	rsp, err := s.fishModifier.DeleteFish(ctx, req.GetId())
	if err != nil {
		return nil, errors.Wrap(err, "unable to delete fish")
	}

	return &trackmyfishv1alpha1.DeleteFishResponse{
		Fish: &trackmyfishv1alpha1.Fish{
			Id:           rsp.ID,
			Type:         rsp.Type,
			Subtype:      rsp.Subtype,
			Color:        rsp.Color,
			Gender:       stringToGender(rsp.Gender),
			PurchaseDate: rsp.PurchaseDate,
			Count:        rsp.Count,
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

	rsp, err := s.tankStatModifier.InsertTankStatistic(ctx, ts)
	if err != nil {
		return nil, errors.Wrap(err, "unable to add tank statistic")
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
	rsp, err := s.tankStatQuerier.ListTankStatistics(ctx)
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
	rsp, err := s.tankStatModifier.DeleteTankStatistic(ctx, req.GetId())
	if err != nil {
		return nil, errors.Wrap(err, "unable to delete tank statistic")
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
	switch strings.ToUpper(gender) {
	case trackmyfishv1alpha1.Fish_MALE.String():
		return trackmyfishv1alpha1.Fish_MALE
	case trackmyfishv1alpha1.Fish_FEMALE.String():
		return trackmyfishv1alpha1.Fish_FEMALE
	}

	return trackmyfishv1alpha1.Fish_UNSPECIFIED
}
