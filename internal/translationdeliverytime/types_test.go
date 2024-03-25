package translationdeliverytime_test

import (
	"backend-engineering-challenge/internal/translationdeliverytime"
	"math"
	"testing"
	"time"
)

func TestTranslationDeliveredUnmarshalJSON(t *testing.T) {
	type testCase struct {
		name string
		json []byte
		exp  func() translationdeliverytime.TranslationDelivered
	}
	tcs := []testCase{
		{
			name: "average is an int",
			json: []byte(`{"timestamp": "2018-12-26 18:11:08.509654","translation_id": "5aa5b2f39f7254a75aa5","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 20}`),
			exp: func() translationdeliverytime.TranslationDelivered {
				time, err := time.Parse("2006-01-02 15:04:05.000000", "2018-12-26 18:11:08.509654")
				if err != nil {
					t.Errorf("unexpected error: %s\n", err)
				}
				return translationdeliverytime.TranslationDelivered{
					Timestamp: time,
					Duration:  20,
				}
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var td translationdeliverytime.TranslationDelivered
			err := td.UnmarshalJSON(tc.json)
			if err != nil {
				t.Errorf("unexpected error: %s\n", err)
			}

			exp := tc.exp()
			if !td.Timestamp.Equal(exp.Timestamp) || td.Duration != exp.Duration {
				t.Errorf("unexpected TranslationDelivered struct. got: %v, exp: %v\n", td, exp)
			}
		})
	}
}

func TestAverageDeliveryTimeString(t *testing.T) {
	type testCase struct {
		name string
		adt  translationdeliverytime.AverageDeliveryTime
		exp  string
	}
	tcs := []testCase{
		{
			name: "average is an int",
			adt: translationdeliverytime.AverageDeliveryTime{
				Date:                time.Date(2024, 2, 25, 15, 38, 1, 0, time.UTC),
				AverageDeliveryTime: 1,
			},
			exp: `{"date": "2024-02-25 15:38:01", "average_delivery_time": 1}`,
		},
		{
			name: "average is a float",
			adt: translationdeliverytime.AverageDeliveryTime{
				Date:                time.Date(2024, 2, 25, 15, 38, 1, 0, time.UTC),
				AverageDeliveryTime: 2.5,
			},
			exp: `{"date": "2024-02-25 15:38:01", "average_delivery_time": 2.5}`,
		},
		{
			name: "average is NaN",
			adt: translationdeliverytime.AverageDeliveryTime{
				Date:                time.Date(2024, 2, 25, 15, 38, 1, 0, time.UTC),
				AverageDeliveryTime: math.NaN(),
			},
			exp: `{"date": "2024-02-25 15:38:01", "average_delivery_time": 0}`,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.adt.String()
			if got != tc.exp {
				t.Errorf("unexpected adt string. got: %s, exp: %s\n", got, tc.exp)
			}
		})
	}
}
