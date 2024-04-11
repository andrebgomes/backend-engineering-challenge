package translationdeliverytime_test

import (
	"backend-engineering-challenge/internal/translationdeliverytime"
	"errors"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	type testCase struct {
		name        string
		inputEvents string
		windowSize  int
		exp         string
		err         error
	}
	tcs := []testCase{
		{
			name:        "success with events",
			inputEvents: readFile(t, "testdata/input"),
			windowSize:  10,
			exp:         readFile(t, "testdata/output"),
			err:         nil,
		},
		{
			name:        "success without events",
			inputEvents: "",
			windowSize:  10,
			exp:         "",
			err:         translationdeliverytime.ErrNoEvents,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := translationdeliverytime.NewTranslationDeliveryTimeApp(tc.inputEvents, tc.windowSize)

			got, err := app.Run()
			if !errors.Is(err, tc.err) {
				t.Errorf("unexpected error. got: %s, exp: %s\n", err, tc.err)
			}

			if got != tc.exp {
				t.Errorf("unexpected restult. got: %s, exp: %s\n", got, tc.exp)
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	inputEvents := readFile(&testing.T{}, "testdata/input")
	app := translationdeliverytime.NewTranslationDeliveryTimeApp(inputEvents, 10)

	for i := 0; i < b.N; i++ {
		_, _ = app.Run()
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
