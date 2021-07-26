package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type Manager struct {
	pool *pgxpool.Pool
}

var _ error = (*ErrNotFound)(nil) // ensure CustomError implements error

type ErrNotFound struct {
	message string
}

func (c *ErrNotFound) Error() string {
	return c.message
}

type Fish struct {
	ID                 int32
	Genus              string
	Species            string
	CommonName         string
	Name               string
	Color              string
	Gender             string
	PurchaseDate       string
	EcosystemName      string
	EcosystemType      string
	EchosystemLocation string
	Salinity           string
	Climate            string
}

func New(c Config) (*Manager, error) {
	if c.Host == "" {
		return nil, errors.New("host not defined")
	}

	if c.Port == 0 {
		return nil, errors.New("port not defined")
	}

	if c.Username == "" {
		return nil, errors.New("user not defined")
	}

	if c.Password == "" {
		return nil, errors.New("password not defined")
	}

	if c.Database == "" {
		return nil, errors.New("database not defined")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.Username, c.Password, c.Host, c.Port, c.Database)

	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	return &Manager{pool: pool}, nil
}

func (d *Manager) InsertFish(ctx context.Context, fish Fish) (Fish, error) {
	f := Fish{}

	err := d.pool.QueryRow(
		ctx,
		"INSERT INTO fish(genus, species, common_name, name, color, gender, purchase_date, ecosystem_name, ecosystem_type, ecosystem_location, salinity, climate) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, genus, species, common_name, name, color, gender, purchase_date, ecosystem_name, ecosystem_type, ecosystem_location, salinity, climate",
		fish.Genus, fish.Species, fish.CommonName, fish.Name, fish.Color, fish.Gender, fish.PurchaseDate, fish.EcosystemName, fish.EcosystemType, fish.EchosystemLocation, fish.Salinity, fish.Climate,
	).Scan(&f.ID, &f.Genus, &f.Species, &f.CommonName, &f.Name, &f.Color, &f.Gender, &f.PurchaseDate, &f.EcosystemName, &f.EcosystemType, &f.EchosystemLocation, &f.Salinity, &f.Climate)
	if err != nil {
		return f, errors.Wrap(err, "unable to add fish")
	}

	logrus.WithFields(logrus.Fields{
		"id": f.ID,
	}).Info("Fish inserted successfully")

	return f, nil
}
