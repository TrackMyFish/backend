//go:build integration
// +build integration

package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/trackmyfish/backend/internal/db"
)

var mgr *db.Manager

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	fmt.Println("Creating test container...")

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=username",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://username:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	var conn *sql.DB

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		conn, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err := createTables(conn); err != nil {
		log.Fatalf("Could not create tables: %s", err)
	}

	mgr, err = db.New(db.Config{
		Host:     "localhost",
		Port:     resource.GetPort("5432/tcp"),
		Username: "username",
		Password: "secret",
		Database: "dbname",
	})
	if err != nil {
		log.Fatalf("Could not create new db instance: %s", err)
	}

	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables(conn *sql.DB) error {
	// Fish Table
	query := `CREATE TABLE IF NOT EXISTS "fish" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "color" VARCHAR(255) DEFAULT '',
  "gender" VARCHAR(255) DEFAULT '',
  "purchase_date" VARCHAR(255) DEFAULT '',
	"count" INT DEFAULT 0,
	"type" VARCHAR(255) DEFAULT '',
	"subtype" VARCHAR(255) DEFAULT '',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	if _, err := conn.Exec(query); err != nil {
		return err
	}

	// Tank Statistics Table
	query = `CREATE TABLE IF NOT EXISTS "tank_statistics" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "test_date" VARCHAR(255) DEFAULT NULL,
  "ph" FLOAT DEFAULT NULL,
  "gh" FLOAT DEFAULT NULL,
  "kh" FLOAT DEFAULT NULL,
  "ammonia" FLOAT DEFAULT NULL,
  "nitrite" FLOAT DEFAULT NULL,
  "nitrate" FLOAT DEFAULT NULL,
  "phosphate" FLOAT DEFAULT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	if _, err := conn.Exec(query); err != nil {
		return err
	}

	// Tanks table
	query = `CREATE TABLE IF NOT EXISTS "tanks" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "make" VARCHAR(40) DEFAULT '',
  "model" VARCHAR(40) DEFAULT '',
  "name" VARCHAR(40) DEFAULT '',
  "location" VARCHAR(40) DEFAULT '',
  "capacity_measurement" VARCHAR(10) DEFAULT '',
  "capacity" FLOAT DEFAULT NULL,
  "description" VARCHAR(255) DEFAULT '',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	if _, err := conn.Exec(query); err != nil {
		return err
	}

	return nil
}

func TestPing(t *testing.T) {
	t.Run("Given a initialised db manager", func(t *testing.T) {
		t.Run("When ping is called", func(t *testing.T) {
			t.Run("Then no error is returned", func(t *testing.T) {
				assert.NoError(t, mgr.Ping(context.Background()))
			})
		})
	})
}

func TestFish(t *testing.T) {
	t.Run("Given a valid Fish object", func(t *testing.T) {
		var inserted db.Fish
		var err error

		fish := db.Fish{Type: "Gourami", Subtype: "Pearl", Color: "Red", Gender: "Male", PurchaseDate: "2021-08-01", Count: 10}

		t.Run("When it is passed to InsertFish", func(t *testing.T) {
			t.Run("Then it should create the record without error", func(t *testing.T) {
				inserted, err = mgr.InsertFish(context.Background(), fish)
				assert.NoError(t, err)
				assert.NotNil(t, inserted)

				assert.Equal(t, fish.Type, inserted.Type)
				assert.Equal(t, fish.Subtype, inserted.Subtype)
				assert.Equal(t, fish.Color, inserted.Color)
				assert.Equal(t, fish.Gender, inserted.Gender)
				assert.Equal(t, fish.PurchaseDate, inserted.PurchaseDate)
				assert.Equal(t, fish.Count, inserted.Count)
			})
		})

		t.Run("When ListFish is called", func(t *testing.T) {
			t.Run("Then the inserted Fish should exist", func(t *testing.T) {
				f, err := mgr.ListFish(context.Background())
				assert.NoError(t, err)

				assert.Len(t, f, 1)

				assert.Equal(t, inserted.ID, f[0].ID)
				assert.Equal(t, fish.Type, f[0].Type)
				assert.Equal(t, fish.Subtype, f[0].Subtype)
				assert.Equal(t, fish.Color, f[0].Color)
				assert.Equal(t, fish.Gender, f[0].Gender)
				assert.Equal(t, fish.PurchaseDate, f[0].PurchaseDate)
				assert.Equal(t, fish.Count, f[0].Count)
			})
		})

		t.Run("When DeleteFish is called", func(t *testing.T) {
			t.Run("Then the Fish is deleted", func(t *testing.T) {
				f, err := mgr.DeleteFish(context.Background(), inserted.ID)
				assert.NoError(t, err)
				assert.NotNil(t, f)

				assert.Equal(t, inserted.ID, f.ID)
				assert.Equal(t, fish.Type, f.Type)
				assert.Equal(t, fish.Subtype, f.Subtype)
				assert.Equal(t, fish.Color, f.Color)
				assert.Equal(t, fish.Gender, f.Gender)
				assert.Equal(t, fish.PurchaseDate, f.PurchaseDate)
				assert.Equal(t, fish.Count, f.Count)

				// Make sure the fish doesn't exist
				lf, err := mgr.ListFish(context.Background())
				assert.NoError(t, err)

				assert.Len(t, lf, 0)
			})
		})
	})
}

func TestTankStatistics(t *testing.T) {
	t.Run("Given a valid TankStatistics object", func(t *testing.T) {
		var inserted db.TankStatistic
		var err error

		tankStat := db.TankStatistic{
			TestDate:  "2021-04-03",
			PH:        pointy.Float32(7.2),
			GH:        pointy.Float32(1.3),
			KH:        pointy.Float32(23),
			Ammonia:   pointy.Float32(10),
			Nitrite:   pointy.Float32(0),
			Nitrate:   pointy.Float32(0.3),
			Phosphate: pointy.Float32(13),
		}

		t.Run("When it is passed to InsertTankStatistic", func(t *testing.T) {
			t.Run("Then it should create the record without error", func(t *testing.T) {
				inserted, err = mgr.InsertTankStatistic(context.Background(), tankStat)
				assert.NoError(t, err)
				assert.NotNil(t, inserted)

				assert.Equal(t, tankStat.TestDate, inserted.TestDate)
				assert.Equal(t, tankStat.PH, inserted.PH)
				assert.Equal(t, tankStat.GH, inserted.GH)
				assert.Equal(t, tankStat.KH, inserted.KH)
				assert.Equal(t, tankStat.Ammonia, inserted.Ammonia)
				assert.Equal(t, tankStat.Nitrite, inserted.Nitrite)
				assert.Equal(t, tankStat.Nitrate, inserted.Nitrate)
				assert.Equal(t, tankStat.Phosphate, inserted.Phosphate)
			})
		})

		t.Run("When ListTankStatistics is called", func(t *testing.T) {
			t.Run("Then the inserted TankStatistic should exist", func(t *testing.T) {
				ts, err := mgr.ListTankStatistics(context.Background())
				assert.NoError(t, err)

				assert.Len(t, ts, 1)

				assert.Equal(t, inserted.ID, ts[0].ID)
				assert.Equal(t, tankStat.TestDate, ts[0].TestDate)
				assert.Equal(t, tankStat.PH, ts[0].PH)
				assert.Equal(t, tankStat.GH, ts[0].GH)
				assert.Equal(t, tankStat.KH, ts[0].KH)
				assert.Equal(t, tankStat.Ammonia, ts[0].Ammonia)
				assert.Equal(t, tankStat.Nitrite, ts[0].Nitrite)
				assert.Equal(t, tankStat.Nitrate, ts[0].Nitrate)
				assert.Equal(t, tankStat.Phosphate, ts[0].Phosphate)
			})
		})

		t.Run("When DeleteTankStatistic is called", func(t *testing.T) {
			t.Run("Then the TankStatistic is deleted", func(t *testing.T) {
				ts, err := mgr.DeleteTankStatistic(context.Background(), inserted.ID)
				assert.NoError(t, err)
				assert.NotNil(t, ts)

				assert.Equal(t, inserted.ID, ts.ID)
				assert.Equal(t, tankStat.TestDate, ts.TestDate)
				assert.Equal(t, tankStat.PH, ts.PH)
				assert.Equal(t, tankStat.GH, ts.GH)
				assert.Equal(t, tankStat.KH, ts.KH)
				assert.Equal(t, tankStat.Ammonia, ts.Ammonia)
				assert.Equal(t, tankStat.Nitrite, ts.Nitrite)
				assert.Equal(t, tankStat.Nitrate, ts.Nitrate)
				assert.Equal(t, tankStat.Phosphate, ts.Phosphate)

				// Make sure the stat doesn't exist
				lts, err := mgr.ListTankStatistics(context.Background())
				assert.NoError(t, err)

				assert.Len(t, lts, 0)
			})
		})
	})
}

func TestTanks(t *testing.T) {
	t.Run("Given a valid Tank object", func(t *testing.T) {
		var inserted db.Tank
		var err error

		tank := db.Tank{
			Make:                "Jewel",
			Model:               "Rio 180 LED",
			Name:                "Main",
			Location:            "Office",
			CapacityMeasurement: "Litres",
			Capacity:            pointy.Float32(180),
			Description:         "Semi-Aggressive tank",
		}

		t.Run("When it is passed to InsertTank", func(t *testing.T) {
			t.Run("Then it should create the record without error", func(t *testing.T) {
				inserted, err = mgr.InsertTank(context.Background(), tank)
				assert.NoError(t, err)
				assert.NotNil(t, inserted)

				assert.Equal(t, tank.Make, inserted.Make)
				assert.Equal(t, tank.Model, inserted.Model)
				assert.Equal(t, tank.Name, inserted.Name)
				assert.Equal(t, tank.Location, inserted.Location)
				assert.Equal(t, tank.CapacityMeasurement, inserted.CapacityMeasurement)
				assert.Equal(t, tank.Capacity, inserted.Capacity)
				assert.Equal(t, tank.Description, inserted.Description)
			})
		})

		t.Run("When ListTanks is called", func(t *testing.T) {
			t.Run("Then the inserted Tank should exist", func(t *testing.T) {
				ts, err := mgr.ListTanks(context.Background())
				assert.NoError(t, err)

				assert.Len(t, ts, 1)

				assert.Equal(t, tank.Make, ts[0].Make)
				assert.Equal(t, tank.Model, ts[0].Model)
				assert.Equal(t, tank.Name, ts[0].Name)
				assert.Equal(t, tank.Location, ts[0].Location)
				assert.Equal(t, tank.CapacityMeasurement, ts[0].CapacityMeasurement)
				assert.Equal(t, tank.Capacity, ts[0].Capacity)
				assert.Equal(t, tank.Description, ts[0].Description)
			})
		})

		t.Run("When DeleteTank is called", func(t *testing.T) {
			t.Run("Then the Tank is deleted", func(t *testing.T) {
				tank, err := mgr.DeleteTank(context.Background(), inserted.ID)
				assert.NoError(t, err)
				assert.NotNil(t, tank)

				assert.Equal(t, tank.Make, inserted.Make)
				assert.Equal(t, tank.Model, inserted.Model)
				assert.Equal(t, tank.Name, inserted.Name)
				assert.Equal(t, tank.Location, inserted.Location)
				assert.Equal(t, tank.CapacityMeasurement, inserted.CapacityMeasurement)
				assert.Equal(t, tank.Capacity, inserted.Capacity)
				assert.Equal(t, tank.Description, inserted.Description)

				// Make sure the stat doesn't exist
				lts, err := mgr.ListTanks(context.Background())
				assert.NoError(t, err)

				assert.Len(t, lts, 0)
			})
		})
	})
}
