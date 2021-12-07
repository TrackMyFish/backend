package server

import (
	"context"
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/trackmyfish/backend/internal/db"
	trackmyfishv1alpha1 "github.com/trackmyfish/proto/trackmyfish/v1alpha1"
)

func TestHeartbeat(t *testing.T) {
	s := Server{}

	t.Run("Given an initialised Server", func(t *testing.T) {
		t.Run("When a request is made to Heartbeat", func(t *testing.T) {
			t.Run("Then it should return an empty HeartbeatResponse", func(t *testing.T) {
				r, err := s.Heartbeat(context.Background(), &trackmyfishv1alpha1.HeartbeatRequest{})
				assert.NoError(t, err)
				assert.Equal(t, &trackmyfishv1alpha1.HeartbeatResponse{}, r)
			})
		})
	})
}

func TestAddFish(t *testing.T) {
	fm := &fishMock{}
	s := Server{fishModifier: fm}

	t.Run("Given a request to AddFish", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				fm.err = errors.New("an error")

				r, err := s.AddFish(context.Background(), &trackmyfishv1alpha1.AddFishRequest{})
				assert.EqualError(t, err, "unable to add fish: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When no error is returned", func(t *testing.T) {
			t.Run("Then the Fish is returned to the caller", func(t *testing.T) {
				fm.err = nil
				fm.insertFishResponse = db.Fish{ID: 45, Type: "Snail", Subtype: "Assassin", Gender: "mAlE"}

				r, err := s.AddFish(context.Background(), &trackmyfishv1alpha1.AddFishRequest{})
				assert.NoError(t, err)

				assert.Equal(t, fm.insertFishResponse.ID, r.Fish.Id)
				assert.Equal(t, fm.insertFishResponse.Type, r.Fish.Type)
				assert.Equal(t, fm.insertFishResponse.Subtype, r.Fish.Subtype)
				assert.Equal(t, trackmyfishv1alpha1.Fish_MALE, r.Fish.Gender)
			})
		})
	})
}

func TestListFish(t *testing.T) {
	fm := &fishMock{}
	s := Server{fishQuerier: fm}

	t.Run("Given a request to ListFish", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				fm.err = errors.New("an error")

				r, err := s.ListFish(context.Background(), &trackmyfishv1alpha1.ListFishRequest{})
				assert.EqualError(t, err, "unable to list fish: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When no error is returned", func(t *testing.T) {
			t.Run("Then the Fish are returned to the caller", func(t *testing.T) {
				fm.err = nil
				fm.listFishResponse = []db.Fish{{ID: 45, Type: "Snail", Subtype: "Assassin", Gender: "mAlE"}, {ID: 1, Type: "Gourami", Subtype: "Pearl", Gender: "female"}}

				r, err := s.ListFish(context.Background(), &trackmyfishv1alpha1.ListFishRequest{})
				assert.NoError(t, err)
				assert.Len(t, r.Fish, 2)

				expected := make(map[int32]db.Fish)
				for _, e := range fm.listFishResponse {
					expected[e.ID] = db.Fish{
						ID:           e.ID,
						Type:         e.Type,
						Subtype:      e.Subtype,
						Color:        e.Color,
						Gender:       e.Gender,
						PurchaseDate: e.PurchaseDate,
						Count:        e.Count,
					}
				}

				for _, a := range r.Fish {
					assert.Equal(t, expected[a.Id].Type, a.Type)
					assert.Equal(t, expected[a.Id].Subtype, a.Subtype)
					assert.Equal(t, expected[a.Id].Color, a.Color)
					assert.Equal(t, stringToGender(expected[a.Id].Gender), a.Gender)
					assert.Equal(t, expected[a.Id].PurchaseDate, a.PurchaseDate)
					assert.Equal(t, expected[a.Id].Count, a.Count)
				}
			})
		})
	})
}

func TestDeleteFish(t *testing.T) {
	fm := &fishMock{}
	s := Server{fishModifier: fm}

	t.Run("Given a request to DeleteFish", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				fm.err = errors.New("an error")

				r, err := s.DeleteFish(context.Background(), &trackmyfishv1alpha1.DeleteFishRequest{})
				assert.EqualError(t, err, "unable to delete fish: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When no error is returned", func(t *testing.T) {
			t.Run("Then the Deleted Fish is returned to the caller", func(t *testing.T) {
				fm.err = nil
				fm.deleteFishResponse = db.Fish{ID: 45, Type: "Snail", Subtype: "Assassin", Gender: "mAlE"}

				r, err := s.DeleteFish(context.Background(), &trackmyfishv1alpha1.DeleteFishRequest{})
				assert.NoError(t, err)

				assert.Equal(t, fm.deleteFishResponse.ID, r.Fish.Id)
				assert.Equal(t, fm.deleteFishResponse.Type, r.Fish.Type)
				assert.Equal(t, fm.deleteFishResponse.Subtype, r.Fish.Subtype)
				assert.Equal(t, trackmyfishv1alpha1.Fish_MALE, r.Fish.Gender)
			})
		})
	})
}

func TestAddTankStatistic(t *testing.T) {
	tsm := &tankStatsMock{}
	s := Server{tankStatModifier: tsm}

	t.Run("Given a request to AddTankStatistic", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				tsm.err = errors.New("an error")

				r, err := s.AddTankStatistic(context.Background(), &trackmyfishv1alpha1.AddTankStatisticRequest{})
				assert.EqualError(t, err, "unable to add tank statistic: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When Optional fields aren't specified", func(t *testing.T) {
			t.Run("Then they shouldn't be used", func(t *testing.T) {
				tsm.err = nil

				r, err := s.AddTankStatistic(context.Background(), &trackmyfishv1alpha1.AddTankStatisticRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Nil(t, tsm.insertTankStatisticsRequest.PH)
				assert.Nil(t, tsm.insertTankStatisticsRequest.GH)
				assert.Nil(t, tsm.insertTankStatisticsRequest.KH)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Ammonia)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Nitrite)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Nitrate)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Phosphate)
			})
		})
		t.Run("When Optional fields are specified", func(t *testing.T) {
			t.Run("Then they should be used", func(t *testing.T) {
				tsm.err = nil

				req := &trackmyfishv1alpha1.AddTankStatisticRequest{
					TankStatistic: &trackmyfishv1alpha1.TankStatistic{
						OptionalPh:        &trackmyfishv1alpha1.TankStatistic_Ph{Ph: 1.7},
						OptionalGh:        &trackmyfishv1alpha1.TankStatistic_Gh{Gh: 2.4},
						OptionalKh:        &trackmyfishv1alpha1.TankStatistic_Kh{Kh: 12},
						OptionalAmmonia:   &trackmyfishv1alpha1.TankStatistic_Ammonia{Ammonia: 12.4},
						OptionalNitrite:   &trackmyfishv1alpha1.TankStatistic_Nitrite{Nitrite: 124},
						OptionalNitrate:   &trackmyfishv1alpha1.TankStatistic_Nitrate{Nitrate: 0.25},
						OptionalPhosphate: &trackmyfishv1alpha1.TankStatistic_Phosphate{Phosphate: 122.31},
					},
				}

				r, err := s.AddTankStatistic(context.Background(), req)
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.NotNil(t, tsm.insertTankStatisticsRequest.PH)
				assert.NotNil(t, tsm.insertTankStatisticsRequest.GH)
				assert.NotNil(t, tsm.insertTankStatisticsRequest.KH)
				assert.NotNil(t, tsm.insertTankStatisticsRequest.Ammonia)
				assert.NotNil(t, tsm.insertTankStatisticsRequest.Nitrite)
				assert.NotNil(t, tsm.insertTankStatisticsRequest.Nitrate)
				assert.NotNil(t, tsm.insertTankStatisticsRequest.Phosphate)

				assert.Equal(t, req.GetTankStatistic().GetPh(), *tsm.insertTankStatisticsRequest.PH)
				assert.Equal(t, req.GetTankStatistic().GetGh(), *tsm.insertTankStatisticsRequest.GH)
				assert.Equal(t, req.GetTankStatistic().GetKh(), *tsm.insertTankStatisticsRequest.KH)
				assert.Equal(t, req.GetTankStatistic().GetAmmonia(), *tsm.insertTankStatisticsRequest.Ammonia)
				assert.Equal(t, req.GetTankStatistic().GetNitrite(), *tsm.insertTankStatisticsRequest.Nitrite)
				assert.Equal(t, req.GetTankStatistic().GetNitrate(), *tsm.insertTankStatisticsRequest.Nitrate)
				assert.Equal(t, req.GetTankStatistic().GetPhosphate(), *tsm.insertTankStatisticsRequest.Phosphate)
			})
		})
		t.Run("When pointer values are nil", func(t *testing.T) {
			t.Run("Then they aren't returned", func(t *testing.T) {
				tsm.err = nil

				r, err := s.AddTankStatistic(context.Background(), &trackmyfishv1alpha1.AddTankStatisticRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Nil(t, r.TankStatistic.OptionalPh)
				assert.Nil(t, r.TankStatistic.OptionalGh)
				assert.Nil(t, r.TankStatistic.OptionalKh)
				assert.Nil(t, r.TankStatistic.OptionalAmmonia)
				assert.Nil(t, r.TankStatistic.OptionalNitrite)
				assert.Nil(t, r.TankStatistic.OptionalNitrate)
				assert.Nil(t, r.TankStatistic.OptionalPhosphate)
			})
		})
		t.Run("When pointer values are not nil", func(t *testing.T) {
			t.Run("Then they are returned", func(t *testing.T) {
				tsm.err = nil
				tsm.insertTankStatisticsResponse = db.TankStatistic{
					ID:        4,
					TestDate:  "2020-03-09",
					PH:        pointy.Float32(23),
					GH:        pointy.Float32(2),
					KH:        pointy.Float32(92.3),
					Ammonia:   pointy.Float32(12.9),
					Nitrite:   pointy.Float32(9),
					Nitrate:   pointy.Float32(0.2),
					Phosphate: pointy.Float32(4),
				}

				r, err := s.AddTankStatistic(context.Background(), &trackmyfishv1alpha1.AddTankStatisticRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Equal(t, *tsm.insertTankStatisticsResponse.PH, r.GetTankStatistic().GetPh())
				assert.Equal(t, *tsm.insertTankStatisticsResponse.GH, r.GetTankStatistic().GetGh())
				assert.Equal(t, *tsm.insertTankStatisticsResponse.KH, r.GetTankStatistic().GetKh())
				assert.Equal(t, *tsm.insertTankStatisticsResponse.Ammonia, r.GetTankStatistic().GetAmmonia())
				assert.Equal(t, *tsm.insertTankStatisticsResponse.Nitrite, r.GetTankStatistic().GetNitrite())
				assert.Equal(t, *tsm.insertTankStatisticsResponse.Nitrate, r.GetTankStatistic().GetNitrate())
				assert.Equal(t, *tsm.insertTankStatisticsResponse.Phosphate, r.GetTankStatistic().GetPhosphate())
			})
		})
	})
}

func TestListTankStatistics(t *testing.T) {
	tsm := &tankStatsMock{}
	s := Server{tankStatQuerier: tsm}

	t.Run("Given a request to ListTankStatistics", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				tsm.err = errors.New("an error")

				r, err := s.ListTankStatistics(context.Background(), &trackmyfishv1alpha1.ListTankStatisticsRequest{})
				assert.EqualError(t, err, "unable to get tank statistics: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When no results are returned", func(t *testing.T) {
			t.Run("Then a length of 0 is returned to the caller", func(t *testing.T) {
				tsm.err = nil

				r, err := s.ListTankStatistics(context.Background(), &trackmyfishv1alpha1.ListTankStatisticsRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Len(t, r.TankStatistics, 0)
			})
		})
		t.Run("When pointer values are not nil", func(t *testing.T) {
			t.Run("Then they are returned", func(t *testing.T) {
				tsm.err = nil
				tsm.listTankStatisticsResponse = []db.TankStatistic{{
					ID:        4,
					TestDate:  "2020-03-09",
					PH:        pointy.Float32(23),
					GH:        pointy.Float32(2),
					KH:        pointy.Float32(92.3),
					Ammonia:   pointy.Float32(12.9),
					Nitrite:   pointy.Float32(9),
					Nitrate:   pointy.Float32(0.2),
					Phosphate: pointy.Float32(4),
				}}

				r, err := s.ListTankStatistics(context.Background(), &trackmyfishv1alpha1.ListTankStatisticsRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Len(t, r.TankStatistics, 1)

				assert.Equal(t, *tsm.listTankStatisticsResponse[0].PH, r.GetTankStatistics()[0].GetPh())
				assert.Equal(t, *tsm.listTankStatisticsResponse[0].GH, r.GetTankStatistics()[0].GetGh())
				assert.Equal(t, *tsm.listTankStatisticsResponse[0].KH, r.GetTankStatistics()[0].GetKh())
				assert.Equal(t, *tsm.listTankStatisticsResponse[0].Ammonia, r.GetTankStatistics()[0].GetAmmonia())
				assert.Equal(t, *tsm.listTankStatisticsResponse[0].Nitrite, r.GetTankStatistics()[0].GetNitrite())
				assert.Equal(t, *tsm.listTankStatisticsResponse[0].Nitrate, r.GetTankStatistics()[0].GetNitrate())
				assert.Equal(t, *tsm.listTankStatisticsResponse[0].Phosphate, r.GetTankStatistics()[0].GetPhosphate())
			})
		})
	})
}

func TestDeleteTankStatistic(t *testing.T) {
	tsm := &tankStatsMock{}
	s := Server{tankStatModifier: tsm}

	t.Run("Given a request to DeleteTankStatistic", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				tsm.err = errors.New("an error")

				r, err := s.DeleteTankStatistic(context.Background(), &trackmyfishv1alpha1.DeleteTankStatisticRequest{})
				assert.EqualError(t, err, "unable to delete tank statistic: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When Optional fields aren't specified", func(t *testing.T) {
			t.Run("Then they shouldn't be used", func(t *testing.T) {
				tsm.err = nil

				r, err := s.DeleteTankStatistic(context.Background(), &trackmyfishv1alpha1.DeleteTankStatisticRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Nil(t, tsm.insertTankStatisticsRequest.PH)
				assert.Nil(t, tsm.insertTankStatisticsRequest.GH)
				assert.Nil(t, tsm.insertTankStatisticsRequest.KH)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Ammonia)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Nitrite)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Nitrate)
				assert.Nil(t, tsm.insertTankStatisticsRequest.Phosphate)
			})
		})
		t.Run("When pointer values are not nil", func(t *testing.T) {
			t.Run("Then they are returned", func(t *testing.T) {
				tsm.err = nil
				tsm.deleteTankStatisticsResponse = db.TankStatistic{
					ID:        4,
					TestDate:  "2020-03-09",
					PH:        pointy.Float32(23),
					GH:        pointy.Float32(2),
					KH:        pointy.Float32(92.3),
					Ammonia:   pointy.Float32(12.9),
					Nitrite:   pointy.Float32(9),
					Nitrate:   pointy.Float32(0.2),
					Phosphate: pointy.Float32(4),
				}

				r, err := s.DeleteTankStatistic(context.Background(), &trackmyfishv1alpha1.DeleteTankStatisticRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Equal(t, *tsm.deleteTankStatisticsResponse.PH, r.GetTankStatistic().GetPh())
				assert.Equal(t, *tsm.deleteTankStatisticsResponse.GH, r.GetTankStatistic().GetGh())
				assert.Equal(t, *tsm.deleteTankStatisticsResponse.KH, r.GetTankStatistic().GetKh())
				assert.Equal(t, *tsm.deleteTankStatisticsResponse.Ammonia, r.GetTankStatistic().GetAmmonia())
				assert.Equal(t, *tsm.deleteTankStatisticsResponse.Nitrite, r.GetTankStatistic().GetNitrite())
				assert.Equal(t, *tsm.deleteTankStatisticsResponse.Nitrate, r.GetTankStatistic().GetNitrate())
				assert.Equal(t, *tsm.deleteTankStatisticsResponse.Phosphate, r.GetTankStatistic().GetPhosphate())
			})
		})
	})
}

func TestAddTank(t *testing.T) {
	tm := &tankMock{}
	s := Server{tankModifier: tm}

	t.Run("Given a request to AddTank", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				tm.err = errors.New("an error")

				r, err := s.AddTank(context.Background(), &trackmyfishv1alpha1.AddTankRequest{})
				assert.EqualError(t, err, "unable to add tank: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When Optional fields aren't specified", func(t *testing.T) {
			t.Run("Then they shouldn't be used", func(t *testing.T) {
				tm.err = nil

				r, err := s.AddTank(context.Background(), &trackmyfishv1alpha1.AddTankRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Nil(t, tm.insertTankRequest.Capacity)
			})
		})
		t.Run("When Optional fields are specified", func(t *testing.T) {
			t.Run("Then they should be used", func(t *testing.T) {
				tm.err = nil

				req := &trackmyfishv1alpha1.AddTankRequest{
					Tank: &trackmyfishv1alpha1.Tank{
						OptionalCapacity: &trackmyfishv1alpha1.Tank_Capacity{Capacity: 180},
					},
				}

				r, err := s.AddTank(context.Background(), req)
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.NotNil(t, tm.insertTankRequest.Capacity)
				assert.Equal(t, req.GetTank().GetCapacity(), *tm.insertTankRequest.Capacity)
			})
		})
		t.Run("When pointer values are nil", func(t *testing.T) {
			t.Run("Then they aren't returned", func(t *testing.T) {
				tm.err = nil

				r, err := s.AddTank(context.Background(), &trackmyfishv1alpha1.AddTankRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Nil(t, r.Tank.OptionalCapacity)
			})
		})
		t.Run("When pointer values are not nil", func(t *testing.T) {
			t.Run("Then they are returned", func(t *testing.T) {
				tm.err = nil
				tm.insertTankResponse = db.Tank{
					ID:                  4,
					Make:                "Juewl",
					Model:               "Rio",
					Name:                "Main Semi-Aggressive",
					Location:            "Office",
					CapacityMeasurement: "Litres",
					Capacity:            pointy.Float32(180),
					Description:         "Semi-Aggressive community tank",
				}

				r, err := s.AddTank(context.Background(), &trackmyfishv1alpha1.AddTankRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Equal(t, *tm.insertTankResponse.Capacity, r.GetTank().GetCapacity())
			})
		})
	})
}

func TestListTanks(t *testing.T) {
	tm := &tankMock{}
	s := Server{tankQuerier: tm}

	t.Run("Given a request to ListTanks", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				tm.err = errors.New("an error")

				r, err := s.ListTanks(context.Background(), &trackmyfishv1alpha1.ListTanksRequest{})
				assert.EqualError(t, err, "unable to get tanks: an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When no results are returned", func(t *testing.T) {
			t.Run("Then a length of 0 is returned to the caller", func(t *testing.T) {
				tm.err = nil

				r, err := s.ListTanks(context.Background(), &trackmyfishv1alpha1.ListTanksRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Len(t, r.Tanks, 0)
			})
		})
		t.Run("When pointer values are not nil", func(t *testing.T) {
			t.Run("Then they are returned", func(t *testing.T) {
				tm.err = nil
				tm.listTankResponse = []db.Tank{{
					ID:                  4,
					Make:                "Juwel",
					Model:               "Rio",
					Name:                "Main",
					Location:            "Office",
					CapacityMeasurement: "Litres",
					Capacity:            pointy.Float32(180),
					Description:         "Semi-Aggressive community",
				}}

				r, err := s.ListTanks(context.Background(), &trackmyfishv1alpha1.ListTanksRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Len(t, r.Tanks, 1)

				assert.Equal(t, *tm.listTankResponse[0].Capacity, r.GetTanks()[0].GetCapacity())
			})
		})
	})
}

func TestDeleteTank(t *testing.T) {
	tm := &tankMock{}
	s := Server{tankModifier: tm}

	t.Run("Given a request to DeleteTank", func(t *testing.T) {
		t.Run("When an error is returned", func(t *testing.T) {
			t.Run("Then the error is returned to the caller", func(t *testing.T) {
				tm.err = errors.New("an error")

				r, err := s.DeleteTank(context.Background(), &trackmyfishv1alpha1.DeleteTankRequest{})
				assert.EqualError(t, err, "unable to delete tank : an error")
				assert.Nil(t, r)
			})
		})
		t.Run("When Optional fields aren't specified", func(t *testing.T) {
			t.Run("Then they shouldn't be used", func(t *testing.T) {
				tm.err = nil

				r, err := s.DeleteTank(context.Background(), &trackmyfishv1alpha1.DeleteTankRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Nil(t, tm.insertTankRequest.Capacity)
			})
		})
		t.Run("When pointer values are not nil", func(t *testing.T) {
			t.Run("Then they are returned", func(t *testing.T) {
				tm.err = nil
				tm.deleteTankResponse = db.Tank{
					ID:                  4,
					Make:                "Juwel",
					Model:               "Rio",
					Name:                "Main",
					Location:            "Office",
					CapacityMeasurement: "Litres",
					Capacity:            pointy.Float32(180),
					Description:         "Semi-Aggressive community",
				}

				r, err := s.DeleteTank(context.Background(), &trackmyfishv1alpha1.DeleteTankRequest{})
				assert.NoError(t, err)
				assert.NotNil(t, r)

				assert.Equal(t, *tm.deleteTankResponse.Capacity, r.GetTank().GetCapacity())
			})
		})
	})
}

func TestStringToGender(t *testing.T) {
	testCases := []struct {
		desc     string
		gender   string
		expected trackmyfishv1alpha1.Fish_Gender
	}{
		{
			desc:     "MALE returns Fish_MALE",
			gender:   "MALE",
			expected: trackmyfishv1alpha1.Fish_MALE,
		},
		{
			desc:     "FEMALE returns Fish_FEMALE",
			gender:   "FEMALE",
			expected: trackmyfishv1alpha1.Fish_FEMALE,
		},
		{
			desc:     "Unsupported value returns Fish_UNSPECIFIED",
			gender:   "UNSUPPORTED",
			expected: trackmyfishv1alpha1.Fish_UNSPECIFIED,
		},
		{
			desc:     "Can handle lower-case MALE",
			gender:   "male",
			expected: trackmyfishv1alpha1.Fish_MALE,
		},
		{
			desc:     "Can handle mixed-case MALE",
			gender:   "mAlE",
			expected: trackmyfishv1alpha1.Fish_MALE,
		},
		{
			desc:     "Can handle lower-case FEMALE",
			gender:   "female",
			expected: trackmyfishv1alpha1.Fish_FEMALE,
		},
		{
			desc:     "Can handle mixed-case FEMALE",
			gender:   "FeMaLe",
			expected: trackmyfishv1alpha1.Fish_FEMALE,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.expected, stringToGender(tC.gender))
		})
	}
}

type fishMock struct {
	insertFishResponse db.Fish
	deleteFishResponse db.Fish
	listFishResponse   []db.Fish
	err                error
}

func (f fishMock) InsertFish(context.Context, db.Fish) (db.Fish, error) {
	return f.insertFishResponse, f.err
}

func (f fishMock) DeleteFish(context.Context, int32) (db.Fish, error) {
	return f.deleteFishResponse, f.err
}

func (f fishMock) ListFish(context.Context) ([]db.Fish, error) {
	return f.listFishResponse, f.err
}

type tankStatsMock struct {
	insertTankStatisticsResponse db.TankStatistic
	insertTankStatisticsRequest  db.TankStatistic
	deleteTankStatisticsResponse db.TankStatistic
	listTankStatisticsResponse   []db.TankStatistic
	err                          error
}

func (f *tankStatsMock) InsertTankStatistic(ctx context.Context, req db.TankStatistic) (db.TankStatistic, error) {
	f.insertTankStatisticsRequest = req

	return f.insertTankStatisticsResponse, f.err
}

func (f *tankStatsMock) DeleteTankStatistic(context.Context, int32) (db.TankStatistic, error) {
	return f.deleteTankStatisticsResponse, f.err
}

func (f *tankStatsMock) ListTankStatistics(context.Context) ([]db.TankStatistic, error) {
	return f.listTankStatisticsResponse, f.err
}

type tankMock struct {
	insertTankResponse db.Tank
	insertTankRequest  db.Tank
	deleteTankResponse db.Tank
	listTankResponse   []db.Tank
	err                error
}

func (f *tankMock) InsertTank(ctx context.Context, req db.Tank) (db.Tank, error) {
	f.insertTankRequest = req

	return f.insertTankResponse, f.err
}

func (f *tankMock) DeleteTank(context.Context, int32) (db.Tank, error) {
	return f.deleteTankResponse, f.err
}

func (f *tankMock) ListTanks(context.Context) ([]db.Tank, error) {
	return f.listTankResponse, f.err
}
