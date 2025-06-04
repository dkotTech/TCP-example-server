package internal

import "math/rand"

var quotes = []string{
	"The only true wisdom is in knowing you know nothing. – Socrates",
	"In the middle of difficulty lies opportunity. – Einstein",
	"Knowing yourself is the beginning of all wisdom. – Aristotle",
}

func GetRandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
