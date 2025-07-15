package main

import "time"

func (m model) calculateWPM() int {
	var endTime time.Time
	if m.endTime != nil {
		endTime = *m.endTime
	} else {
		endTime = time.Now()
	}

	duration := endTime.Sub(m.startTime).Minutes()
	correctCount := m.calculateCorrectCount()
	wpm := float64(correctCount) / 5.0 / duration
	return int(wpm)
}

func (m model) calculateCorrectCount() int {
	runesWant := []rune(m.wantedText)
	runesGot := []rune(m.gottenText)

	minLen := min(len(runesWant), len(runesGot))

	var correct int
	for i := range minLen {
		if runesGot[i] == runesWant[i] {
			correct++
		}
	}

	return correct
}
