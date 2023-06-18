package main

import (
	"context"
	"fmt"
	"github.com/mshafiee/gtranslate"
	"golang.org/x/text/language"
	"log"
)

func main() {
	// Text to be translated
	text := "The top two Republican leaders in the Senate remain silent a day after former President Donald Trump, the current GOP 2024 frontrunner, was indicted by the federal government." +
		"\nWhile the charges have yet to be unsealed, the top two Republicans in the Senate, Minority Leader Mitch McConnell and Minority Whip John Thune, have not put out statements." +
		"\nThat's a stark contrast to the swift reaction among House GOP leaders, who quickly rushed to Trump’s defense." +
		"\nToday is indeed a dark day for the United States of America. It is unconscionable for a President to indict the leading candidate opposing him. " +
		"\nJoe Biden kept classified documents for decades, House Speaker Kevin McCarthy tweeted Thursday night."

	// Translate function is called with the text, source language (empty for auto-detection), and target language.
	// The function returns a struct that contains the translation result or an error.
	translate, err := gtranslate.Translate(context.Background(), text, language.Tag{}, language.Persian)
	if err != nil {
		log.Println(err) // Logs error if any occurred during the translation
	}
	fmt.Printf("%+v", translate) // Prints the struct containing the translation result
}
