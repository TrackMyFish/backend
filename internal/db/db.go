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
	Count              int32
}

type TankStatistic struct {
	ID        int32
	TestDate  string
	PH        float32
	GH        float32
	KH        float32
	Ammonia   float32
	Nitrite   float32
	Nitrate   float32
	Phosphate float32
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
		"INSERT INTO fish(genus, species, common_name, name, color, gender, purchase_date, ecosystem_name, ecosystem_type, ecosystem_location, salinity, climate, count) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, genus, species, common_name, name, color, gender, purchase_date, ecosystem_name, ecosystem_type, ecosystem_location, salinity, climate, count",
		fish.Genus, fish.Species, fish.CommonName, fish.Name, fish.Color, fish.Gender, fish.PurchaseDate, fish.EcosystemName, fish.EcosystemType, fish.EchosystemLocation, fish.Salinity, fish.Climate, fish.Count,
	).Scan(&f.ID, &f.Genus, &f.Species, &f.CommonName, &f.Name, &f.Color, &f.Gender, &f.PurchaseDate, &f.EcosystemName, &f.EcosystemType, &f.EchosystemLocation, &f.Salinity, &f.Climate, &f.Count)
	if err != nil {
		return f, errors.Wrap(err, "unable to add fish")
	}

	logrus.WithFields(logrus.Fields{
		"id": f.ID,
	}).Info("Fish inserted successfully")

	return f, nil
}

func (d *Manager) GetFish(ctx context.Context) ([]Fish, error) {
	fish := make([]Fish, 0)

	rows, err := d.pool.Query(ctx, "SELECT id, genus, species, common_name, name, color, gender, purchase_date, ecosystem_name, ecosystem_type, ecosystem_location, salinity, climate, count FROM fish")
	if err != nil {
		return fish, errors.Wrap(err, "unable to get fish")
	}

	rowCount := 0
	for rows.Next() {
		f := Fish{}

		if err := rows.Scan(&f.ID, &f.Genus, &f.Species, &f.CommonName, &f.Name, &f.Color, &f.Gender, &f.PurchaseDate, &f.EcosystemName, &f.EcosystemType, &f.EchosystemLocation, &f.Salinity, &f.Climate, &f.Count); err != nil {
			return nil, errors.Wrap(err, "unable to scan row")
		}

		fish = append(fish, f)

		rowCount++
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "erroring reading rows")
	}

	logrus.WithFields(logrus.Fields{"rowCount": rowCount}).Info("Fish queried successfully")

	return fish, nil
}

func (d *Manager) DeleteFish(ctx context.Context, id int32) (Fish, error) {
	f := Fish{}

	err := d.pool.QueryRow(
		ctx,
		"DELETE FROM fish WHERE id=$1 RETURNING id, genus, species, common_name, name, color, gender, purchase_date, ecosystem_name, ecosystem_type, ecosystem_location, salinity, climate",
		id,
	).Scan(&f.ID, &f.Genus, &f.Species, &f.CommonName, &f.Name, &f.Color, &f.Gender, &f.PurchaseDate, &f.EcosystemName, &f.EcosystemType, &f.EchosystemLocation, &f.Salinity, &f.Climate)
	if err != nil {
		return f, errors.Wrap(err, "unable to delete fish")
	}

	logrus.WithFields(logrus.Fields{
		"id": f.ID,
	}).Info("Fish deleted successfully")

	return f, nil
}

func (d *Manager) InsertTankStatistic(ctx context.Context, tankStatistic TankStatistic) (TankStatistic, error) {
	ts := TankStatistic{}

	err := d.pool.QueryRow(
		ctx,
		"INSERT INTO tank_statistics(ph, gh, kh, ammonia, nitrite, nitrate, phosphate) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, ph, gh, kh, ammonia, nitrite, nitrate, phosphate",
		tankStatistic.PH, tankStatistic.GH, tankStatistic.KH, tankStatistic.Ammonia, tankStatistic.Nitrite, tankStatistic.Nitrate, tankStatistic.Phosphate,
	).Scan(&ts.ID, &ts.PH, &ts.GH, &ts.KH, &ts.Ammonia, &ts.Nitrite, &ts.Nitrate, &ts.Phosphate)
	if err != nil {
		return ts, errors.Wrap(err, "unable to add fish")
	}

	logrus.WithFields(logrus.Fields{
		"id": ts.ID,
	}).Info("Tank Statistic inserted successfully")

	return ts, nil
}

func (d *Manager) GetTankStatistics(ctx context.Context) ([]TankStatistic, error) {
	tankStats := make([]TankStatistic, 0)

	rows, err := d.pool.Query(ctx, "SELECT id, ph, gh, kh, ammonia, nitrite, nitrate, phosphate FROM tank_statistics")
	if err != nil {
		return tankStats, errors.Wrap(err, "unable to get tank statistics")
	}

	rowCount := 0
	for rows.Next() {
		ts := TankStatistic{}

		if err := rows.Scan(&ts.ID, &ts.PH, &ts.GH, &ts.KH, &ts.Ammonia, &ts.Nitrite, &ts.Nitrate, &ts.Phosphate); err != nil {
			return nil, errors.Wrap(err, "unable to scan row")
		}

		tankStats = append(tankStats, ts)

		rowCount++
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "erroring reading rows")
	}

	logrus.WithFields(logrus.Fields{"rowCount": rowCount}).Info("Tank Statistics queried successfully")

	return tankStats, nil
}

func (d *Manager) DeleteTankStatistic(ctx context.Context, id int32) (TankStatistic, error) {
	ts := TankStatistic{}

	err := d.pool.QueryRow(
		ctx,
		"DELETE FROM tank_statistics WHERE id=$1 RETURNING id, ph, gh, kh, ammonia, nitrite, nitrate, phosphate",
		id,
	).Scan(&ts.ID, &ts.PH, &ts.GH, &ts.KH, &ts.Ammonia, &ts.Nitrite, &ts.Nitrate, &ts.Phosphate)
	if err != nil {
		return ts, errors.Wrap(err, "unable to delete tank statistic")
	}

	logrus.WithFields(logrus.Fields{
		"id": ts.ID,
	}).Info("Tank Statistic deleted successfully")

	return ts, nil
}
