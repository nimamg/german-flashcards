// Package srs schedules flashcard reviews using the Leitner box system.
//
// Every card sits in a numbered box. A correct answer moves it up one box,
// a wrong answer sends it back to box 1. The higher the box, the longer the
// card waits before it is shown again.
package srs

import "time"

// Clock tells the current time. Production code uses SystemClock; tests pass
// a fake so results don't depend on when the test happens to run.
type Clock interface {
	Now() time.Time
}

// SystemClock is a Clock backed by the real wall clock.
type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }

const day = 24 * time.Hour

// intervals[i] is how long a card waits after landing in box i+1.
var intervals = [...]time.Duration{
	1 * day,  // box 1
	2 * day,  // box 2
	4 * day,  // box 3
	8 * day,  // box 4
	16 * day, // box 5
}

// NumBoxes is the highest box a card can reach.
const NumBoxes = len(intervals)

// State is the scheduling state of a single card. The zero value is a new
// card: never reviewed (box 0) and due immediately.
type State struct {
	Box int       // 1..NumBoxes, or 0 if never reviewed
	Due time.Time // show the card at or after this time
}

// Scheduler applies the Leitner rules, reading "now" from its Clock.
type Scheduler struct {
	clock Clock
}

// New returns a Scheduler that uses clock for the current time.
func New(clock Clock) Scheduler {
	return Scheduler{clock: clock}
}

// Review returns a card's next state after it was answered. A new card that
// is answered correctly goes straight to box 2.
func (s Scheduler) Review(st State, correct bool) State {
	box := 1
	if correct {
		box = min(max(st.Box, 1)+1, NumBoxes)
	}
	return State{
		Box: box,
		Due: s.clock.Now().Add(intervals[box-1]),
	}
}

// IsDue reports whether the card should be reviewed now.
func (s Scheduler) IsDue(st State) bool {
	return !s.clock.Now().Before(st.Due)
}
