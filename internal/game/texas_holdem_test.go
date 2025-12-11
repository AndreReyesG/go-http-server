package game

import (
	"fmt"
	"io"
	"testing"
	"time"

	"poker/internal/testutils"
)

// scheduledAlert holds information about when an alert is scheduled.
type scheduledAlert struct {
	at     time.Duration
	amount int
}

// String() method print nicely if the test fails.
func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

// SpyBlindAlerter allows you to spy on ScheduleAlertAt calls.
type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

// ScheduleAlertAt records alerts that have been scheduled.
func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

var dummyBlindAlerter = &SpyBlindAlerter{}

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, testutils.DummyPlayerStore)

		game.Start(5, io.Discard)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, testutils.DummyPlayerStore)

		game.Start(7, io.Discard)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	store := &testutils.StubPlayerStore{}
	game := NewTexasHoldem(dummyBlindAlerter, store)
	winner := "Milky"

	game.Finish(winner)
	testutils.AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(t *testing.T, cases []scheduledAlert, blindAlerter *SpyBlindAlerter) {
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.alerts) <= 1 {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}

func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
