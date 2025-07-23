package main

import "time"

func (m typingModel) calculateWPM() int {
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

func (m typingModel) calculateAccuracy() int {
	correctCount := m.calculateCorrectCount()
	accuracy := float64(correctCount) / float64(len(m.gottenText)) * 100
	return int(accuracy)
}

func (m typingModel) calculateCorrectCount() int {
	minLen := min(len(m.wantedText), len(m.gottenText))

	var correct int
	for i := range minLen {
		if m.gottenText[i] == m.wantedText[i] {
			correct++
		}
	}

	return correct
}
