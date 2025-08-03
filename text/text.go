package text

import (
	"math/rand"
	"strings"

	"github.com/peymanier/donkeytype/utils"
)

var passages = []string{
	"The quick brown fox jumps over the lazy dog.",
	"Pack my box with five dozen liquor jugs.",
	"How razorback-jumping frogs can level six piqued gymnasts!",
	"Grumpy wizards make toxic brew for the evil queen and jack.",
	"Watch Jeopardy!, Alex Trebek's fun TV quiz game.",
	"Crazy Fredrick bought many very exquisite opal jewels.",

	"Life is really simple, but we insist on making it complicated.",
	"It always seems impossible until it's done.",
	"Success is not final, failure is not fatal: it is the courage to continue that counts.",
	"Do not dwell in the past, do not dream of the future, concentrate the mind on the present moment.",
	"Knowing yourself is the beginning of all wisdom.",
	"The only true wisdom is in knowing you know nothing.",
	"All we have to decide is what to do with the time that is given us.",
	"Fear is the mind-killer. Fear is the little-death that brings total obliteration.",

	"Typing is not about speed; it's about rhythm, focus, and precision.",
	"Don't look at the keyboard — let your fingers learn the way.",
	"Terminal applications can be beautiful, interactive, and powerful.",
	"Charm's Bubble Tea framework brings functional UI architecture to the terminal.",
	"In Go, composition is preferred over inheritance.",
	"A good developer types with clarity and intent, not haste.",

	"Rain gently tapped the window as the cat curled deeper into its nap.",
	"Silence stretched across the room like a taut wire, ready to snap.",
	"Her fingers danced across the keyboard, painting stories with each stroke.",
	"Somewhere in the world, someone is learning to type for the very first time.",
	"A single keystroke can start a revolution, or end an empire.",
	"The stars whispered secrets to the night, and the sky listened patiently.",

	"Keep your eyes on the screen and trust your muscle memory.",
	"Accuracy first. Speed comes with practice.",
	"Rest your wrists. Breathe deeply. Try again.",
	"Every typo is a step toward mastery.",
	"Practice daily, and soon your fingers will fly.",
}

func RandomPassage() []rune {
	return []rune(passages[rand.Intn(len(passages))])
}

func SamplePassages(count int) []rune {
	joined := strings.Join(utils.Sample(passages, count), " ")
	return []rune(joined)
}

func RemoveLastRune(r []rune) []rune {
	if len(r) == 0 {
		return r
	}

	return r[:len(r)-1]
}

func RandomTextFromChars(chars []rune) []rune {
	var charset []rune

	charPool := append(chars, ' ', ' ')
	for range 200 {
		randIdx := rand.Intn(len(charPool))
		charset = append(charset, charPool[randIdx])
	}

	trimCharset := strings.Trim(string(charset), " ")
	splitCharset := strings.FieldsSeq(trimCharset)

	var result []string
	for substr := range splitCharset {
		if len(substr) < 4 {
			continue
		}

		result = append(result, substr)
	}

	return []rune(strings.Join(result, " "))
}
