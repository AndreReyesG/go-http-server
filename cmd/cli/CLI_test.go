package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"poker/internal/testutils"
)

type SpyBlindAlerter struct {
	alerts []struct {
		scheduledAt time.Duration
		amount      int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct {
		scheduledAt time.Duration
		amount      int
	}{duration, amount})
}

var dummySpyAlerter = &SpyBlindAlerter{}

func TestCLI(t *testing.T) {
	t.Run("record Moka win from user input", func(t *testing.T) {
		in := strings.NewReader("Moka wins\n")
		playerStore := &testutils.StubPlayerStore{}

		cli := NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()

		testutils.AssertPlayerWin(t, playerStore, "Moka")
	})

	t.Run("record Milky win from user input", func(t *testing.T) {
		in := strings.NewReader("Milky wins\n")
		playerStore := &testutils.StubPlayerStore{}

		cli := NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()

		testutils.AssertPlayerWin(t, playerStore, "Milky")
	})

	t.Run("it schedules printing of bind values", func(t *testing.T) {
		in := strings.NewReader("Milky wins\n")
		playerStore := &testutils.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		cases := []struct {
			expectedScheduleTime time.Duration
			expectedAmount       int
		}{
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

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduleTime), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				amountGot := alert.amount
				if amountGot != c.expectedAmount {
					t.Errorf("got amount %d, want %d", amountGot, c.expectedAmount)
				}

				gotScheduledTime := alert.scheduledAt
				if gotScheduledTime != alert.scheduledAt {
					t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, alert.scheduledAt)
				}
			})
		}
	})
}
