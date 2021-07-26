package server

import (
	"context"

	"github.com/trackmyfish/backend/internal/db"
	"github.com/trackmyfish/backend/internal/fishbase"
	trackmyfishv1alpha1 "github.com/trackmyfish/proto/trackmyfish/v1alpha1"
)

const fishbaseURL = "https://fishbase.ropensci.org/"

var errNotFound *db.ErrNotFound

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

func (s *Server) AddFish(ctx context.Context, req *trackmyfishv1alpha1.AddFishRequest) (*trackmyfishv1alpha1.AddFishResponse, error) {
	_, err := fishbase.GetByGenusSpecies(req.GetFish().GetGenus(), req.GetFish().GetSpecies())
	if err != nil {
		return nil, err
	}

	return nil, nil
	// fish, err := s.dbManager.InsertFish(ctx, db.Fish{
	// 	Genus:              req.GetFish().GetGenus(),
	// 	Species:            req.GetFish().GetSpecies(),
	// 	CommonName:         req.GetFish().GetCommonName(),
	// 	Name:               req.GetFish().GetName(),
	// 	Color:              req.GetFish().GetColor(),
	// 	Gender:             req.GetFish().GetGender().String(),
	// 	PurchaseDate:       req.GetFish().GetPurchaseDate(),
	// 	EcosystemName:      req.GetFish().GetEcosystemName(),
	// 	EcosystemType:      req.GetFish().GetEcosystemType(),
	// 	EchosystemLocation: req.GetFish().GetEcosystemLocation(),
	// 	Salinity:           req.GetFish().GetSalinity(),
	// 	Climate:            req.GetFish().GetClimate(),
	// })
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return &trackmyfishv1alpha1.AddFishResponse{
	// 	Fish: &trackmyfishv1alpha1.Fish{
	// 		Id:                fish.ID,
	// 		Genus:             fish.Genus,
	// 		Species:           fish.Species,
	// 		CommonName:        fish.CommonName,
	// 		Name:              fish.Name,
	// 		Color:             fish.Color,
	// 		Gender:            stringToGender(fish.Gender),
	// 		PurchaseDate:      fish.PurchaseDate,
	// 		EcosystemName:     fish.EcosystemName,
	// 		EcosystemType:     fish.EcosystemType,
	// 		EcosystemLocation: fish.EchosystemLocation,
	// 		Salinity:          fish.Salinity,
	// 		Climate:           fish.Climate,
	// 	},
	// }, nil
}

func (s *Server) ListFish(context.Context, *trackmyfishv1alpha1.ListFishRequest) (*trackmyfishv1alpha1.ListFishResponse, error) {
	return nil, nil
}

func (s *Server) DeleteFish(context.Context, *trackmyfishv1alpha1.DeleteFishRequest) (*trackmyfishv1alpha1.DeleteFishResponse, error) {
	return nil, nil
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
