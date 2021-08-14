//+build integration

package db_test

import (
	"context"
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/trackmyfish/backend/internal/db"
)

func TestFish(t *testing.T) {
	m, err := db.New(db.Config{
		Host:     "localhost",
		Port:     15432,
		Username: "user",
		Password: "password",
		Database: "trackmyfishtests",
	})

	assert.NoError(t, err)
	assert.NotNil(t, m)

	t.Run("Given a valid Fish object", func(t *testing.T) {
		var inserted db.Fish
		var err error

		fish := db.Fish{Type: "Gourami", Subtype: "Pearl", Color: "Red", Gender: "Male", PurchaseDate: "2021-08-01", Count: 10}

		t.Run("When it is passed to InsertFish", func(t *testing.T) {
			t.Run("Then it should create the record without error", func(t *testing.T) {
				inserted, err = m.InsertFish(context.Background(), fish)
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
				f, err := m.ListFish(context.Background())
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
				f, err := m.DeleteFish(context.Background(), inserted.ID)
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
				lf, err := m.ListFish(context.Background())
				assert.NoError(t, err)

				assert.Len(t, lf, 0)
			})
		})
	})
}

func TestTankStatistics(t *testing.T) {
	m, err := db.New(db.Config{
		Host:     "localhost",
		Port:     15432,
		Username: "user",
		Password: "password",
		Database: "trackmyfishtests",
	})

	assert.NoError(t, err)
	assert.NotNil(t, m)

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
				inserted, err = m.InsertTankStatistic(context.Background(), tankStat)
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
				ts, err := m.ListTankStatistics(context.Background())
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
				ts, err := m.DeleteTankStatistic(context.Background(), inserted.ID)
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
				lts, err := m.ListTankStatistics(context.Background())
				assert.NoError(t, err)

				assert.Len(t, lts, 0)
			})
		})
	})
}
