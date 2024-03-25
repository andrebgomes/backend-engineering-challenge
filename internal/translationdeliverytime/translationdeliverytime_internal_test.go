package translationdeliverytime

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParseTranslationDeliveredEvents(t *testing.T) {
	type exp struct {
		events eventsMap
		first  time.Time
		last   time.Time
	}
	type testCase struct {
		name        string
		inputEvents string
		windowSize  int
		exp         exp
		err         error
	}
	tcs := []testCase{
		{
			name:        "success with events",
			inputEvents: readFile(t, "testdata/input"),
			windowSize:  10,
			exp: exp{
				events: eventsMap{
					"2018-12-26 18:12": {sum: 20, q: 1},
					"2018-12-26 18:16": {sum: 31, q: 1},
					"2018-12-26 18:24": {sum: 54, q: 1},
				},
				first: time.Date(2018, 12, 26, 18, 11, 0, 0, time.UTC),
				last:  time.Date(2018, 12, 26, 18, 24, 0, 0, time.UTC),
			},
			err: nil,
		},
		{
			name:        "success with events in the same minute",
			inputEvents: readFile(t, "testdata/input2"),
			windowSize:  10,
			exp: exp{
				events: eventsMap{
					"2018-12-26 18:12": {sum: 20, q: 1},
					"2018-12-26 18:16": {sum: 71, q: 2},
					"2018-12-26 18:24": {sum: 54, q: 1},
				},
				first: time.Date(2018, 12, 26, 18, 11, 0, 0, time.UTC),
				last:  time.Date(2018, 12, 26, 18, 24, 0, 0, time.UTC),
			},
			err: nil,
		},
		{
			name:        "success without events",
			inputEvents: "",
			windowSize:  10,
			exp: exp{
				events: nil,
				first:  time.Time{},
				last:   time.Time{},
			},
			err: ErrNoEvents,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := TranslationDeliveryTimeApp{
				inputEvents: tc.inputEvents,
				windowSize:  tc.windowSize,
			}

			events, first, last, err := app.parseTranslationDeliveredEvents()
			if !errors.Is(err, tc.err) {
				t.Errorf("unexpected error. got: %s, exp: %s\n", err, tc.err)
			}

			if !reflect.DeepEqual(events, tc.exp.events) || !first.Equal(tc.exp.first) || !last.Equal(tc.exp.last) {
				t.Errorf("unexpected restult. got: %v %v %v, exp: %v %v %v\n", events, first, last, tc.exp.events, tc.exp.first, tc.exp.last)
			}
		})
	}
}

func readFile(t *testing.T, file string) string {
	t.Helper()
	events, err := os.ReadFile(file)
	if err != nil {
		t.Errorf("unable to read file[%v]: %s", file, err)
	}
	return string(events)
}
