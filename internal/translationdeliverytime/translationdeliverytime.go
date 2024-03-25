// Package translationdeliverytime implements the parsing and calculation
// of the delivery time of a translation as a SLA.
//
// This SLA represents a moving average of the translation delivery
// time for the last X minutes.
package translationdeliverytime

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TranslationDeliveryTimeApp is to be used for calculating the Translation Delivery Time SLA.
type TranslationDeliveryTimeApp struct {
	inputEvents string
	windowSize  int
}

// NewTranslationDeliveryTimeApp returns an instance of TranslationDeliveryTimeApp.
// inputEvents 	- string containing the stream of events
// windowSize 	- minutes taken into account for the average
func NewTranslationDeliveryTimeApp(inputEvents string, windowSize int) TranslationDeliveryTimeApp {
	return TranslationDeliveryTimeApp{
		inputEvents: inputEvents,
		windowSize:  windowSize,
	}
}

// Run will calculate the moving average according to the app configs.
// Returns a string with the output.
func (t TranslationDeliveryTimeApp) Run() (string, error) {
	events, first, last, err := t.parseTranslationDeliveredEvents()
	if err != nil {
		return "", err
	}
	output := t.createOutput(events, first, last)
	return output, nil
}

func (t TranslationDeliveryTimeApp) parseTranslationDeliveredEvents() (eventsMap, time.Time, time.Time, error) {
	m := make(eventsMap)
	events := []TranslationDelivered{}

	// Go through each line of the stream events
	scanner := bufio.NewScanner(strings.NewReader(t.inputEvents))
	for scanner.Scan() {
		// Convert it into a struct
		var td TranslationDelivered
		if err := json.Unmarshal([]byte(scanner.Text()), &td); err != nil {
			return nil, time.Time{}, time.Time{}, err
		}

		// Append the event into a slice in order to later get the first and last
		events = append(events, td)

		// Fill the eventsMap with the sum and ammount of events per minute
		key := ceilTime(td.Timestamp).Format("2006-01-02 15:04")
		d := durations{
			sum: td.Duration,
			q:   1,
		}
		if val, ok := m[key]; ok {
			d.sum += val.sum
			d.q += val.q
		}
		m[key] = d
	}

	if len(m) == 0 {
		return nil, time.Time{}, time.Time{}, ErrNoEvents
	}

	// Assert the first and last minute of the output
	first := events[0].Timestamp.Truncate(time.Minute)
	last := ceilTime(events[len(events)-1].Timestamp)

	return m, first, last, nil
}

func (t TranslationDeliveryTimeApp) createOutput(events eventsMap, first, last time.Time) string {
	var sb strings.Builder
	for current := first; current.Before(last) || current.Equal(last); current = current.Add(time.Minute) {
		line := AverageDeliveryTime{
			Date:                current,
			AverageDeliveryTime: t.calculateWindowAverage(events, current, t.windowSize),
		}
		sb.WriteString(fmt.Sprintf("%s\n", line))
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

func (t TranslationDeliveryTimeApp) calculateWindowAverage(events eventsMap, current time.Time, windowSize int) float64 {
	// Get the sums and amount of events for each minute for the last windowSize minutes
	windowDurations := []durations{}
	for i := 0; i < windowSize; i++ {
		s := current.Format("2006-01-02 15:04")
		windowDurations = append(windowDurations, events[s])
		current = current.Add(-time.Minute)
	}

	// Calculate the average
	sum := 0
	q := 0
	for _, d := range windowDurations {
		sum += d.sum
		q += d.q
	}
	return float64(sum) / float64(q)
}

// ceilTime will return the time truncated to the minute above if its not already truncated (0 seconds and 0 nanoseconds).
func ceilTime(t time.Time) time.Time {
	if t.Second() > 0 || t.Nanosecond() > 0 {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Add(time.Minute).Minute(), 0, 0, time.UTC)
	}
	return t
}
