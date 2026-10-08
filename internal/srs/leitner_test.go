package srs

import (
	"testing"
	"time"
)

// fakeClock is a Clock whose time only moves when a test moves it.
type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time          { return c.now }
func (c *fakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

var t0 = time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC)

func TestReview(t *testing.T) {
	tests := []struct {
		name    string
		box     int
		correct bool
		wantBox int
		wantDue time.Time
	}{
		{"new card, correct", 0, true, 2, t0.Add(2 * day)},
		{"new card, wrong", 0, false, 1, t0.Add(1 * day)},
		{"box 1, correct", 1, true, 2, t0.Add(2 * day)},
		{"box 3, correct", 3, true, 4, t0.Add(8 * day)},
		{"top box stays on top", NumBoxes, true, NumBoxes, t0.Add(16 * day)},
		{"box 4, wrong", 4, false, 1, t0.Add(1 * day)},
		{"top box, wrong", NumBoxes, false, 1, t0.Add(1 * day)},
	}

	s := New(&fakeClock{now: t0})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.Review(State{Box: tt.box}, tt.correct)
			if got.Box != tt.wantBox {
				t.Errorf("Box = %d, want %d", got.Box, tt.wantBox)
			}
			if !got.Due.Equal(tt.wantDue) {
				t.Errorf("Due = %v, want %v", got.Due, tt.wantDue)
			}
		})
	}
}

func TestIsDue(t *testing.T) {
	tests := []struct {
		name string
		due  time.Time
		want bool
	}{
		{"new card (zero time)", time.Time{}, true},
		{"due in the past", t0.Add(-time.Hour), true},
		{"due exactly now", t0, true},
		{"due in the future", t0.Add(time.Second), false},
	}

	s := New(&fakeClock{now: t0})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.IsDue(State{Due: tt.due}); got != tt.want {
				t.Errorf("IsDue = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCardLifecycle follows one card across several days of reviews.
func TestCardLifecycle(t *testing.T) {
	clock := &fakeClock{now: t0}
	s := New(clock)

	st := s.Review(State{}, true) // box 2, due in 2 days
	clock.Advance(1 * day)
	if s.IsDue(st) {
		t.Fatalf("due after 1 day, want not due until day 2")
	}

	clock.Advance(1 * day)
	if !s.IsDue(st) {
		t.Fatalf("not due after 2 days")
	}

	st = s.Review(st, false) // forgot it: back to box 1
	if st.Box != 1 || !st.Due.Equal(clock.now.Add(day)) {
		t.Fatalf("after wrong answer got %+v, want box 1 due tomorrow", st)
	}
}
