package gtranslate

import (
	"golang.org/x/text/language"
)

// Sentence represents a sentence's content, its translation, pronunciation, and frequency.
type Sentence struct {
	Content       string  // The original content of the sentence.
	Translation   string  // The translated content of the sentence.
	Pronunciation string  // The pronunciation of the sentence in the target language.
	Frequency     float64 // The frequency of the sentence appearing in the corpus of data.
}

// Equivalent represents an equivalent word in the target language.
type Equivalent struct {
	Content     string   // The original content of the word.
	Equivalents []string // The equivalent words in the target language.
	Frequency   float64  // The frequency of the equivalent words appearing in the corpus of data.
}

// WordTranslation represents a word translation along with its equivalent in the target language.
type WordTranslation struct {
	PartsOfSentence string       // The part of speech for the word.
	Frequency       float64      // The frequency of the word appearing in the corpus of data.
	Translations    []string     // The translations of the word.
	Equivalents     []Equivalent // The equivalent translations of the word.
}

// AlternateTranslation represents an alternate translation for the original content.
type AlternateTranslation struct {
	Content      string        // The original content of the translation.
	Translations []Translation // The alternate translations.
}

// Translation represents a translation and its type.
type Translation struct {
	Translation string // The translated content.
	IsCommon    bool   // A boolean value indicating if the translation is commonly used.
	IsInformal  bool   // A boolean value indicating if the translation is informal.
}

// Synonym represents a synonym for a word in the target language.
type Synonym struct {
	Category string   // The category of the synonym.
	Synonyms []string // The synonyms.
}

// WordSynonym represents a word synonym in the target language.
type WordSynonym struct {
	PartsOfSentence string    // The part of speech for the word.
	Synonyms        []Synonym // The synonyms for the word.
	Contents        string    // The original content of the word.
	Frequency       float64   // The frequency of the word appearing in the corpus of data.
}

// Definition represents a definition of a word along with an example.
type Definition struct {
	Definition string // The definition of the word.
	Example    string // An example using the word.
}

// WordDefinition represents a word definition in the target language.
type WordDefinition struct {
	PartsOfSentence string       // The part of speech for the word.
	Definitions     []Definition // The definitions of the word.
}

// TranslationResult represents the full result of a translation request.
type TranslationResult struct {
	Content               string                 // The original content to translate.
	Translation           string                 // The translated content.
	Pronunciation         string                 // The pronunciation of the translated content.
	TranslatedSentences   []Sentence             // The sentences after translation.
	WordTranslations      []WordTranslation      // The translations of the words in the content.
	SourceLanguage        language.Tag           // The source language of the original content.
	AlternateTranslations []AlternateTranslation // The alternate translations for the original content.
	WordSynonyms          []WordSynonym          // The synonyms of the words in the content.
	WordDefinitions       []WordDefinition       // The definitions of the words in the content.
	WordExamples          []string               // The examples of the words in the content.
}
