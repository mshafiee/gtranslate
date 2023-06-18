package gtranslate

import (
	"encoding/json"
	"errors"
	"golang.org/x/text/language"
)

// safeGetBool is a utility function to safely retrieve a boolean value from a slice of interface{} at the given index
func safeGetBool(slice []interface{}, index int) bool {
	if len(slice) > index {
		if str, isValid := slice[index].(bool); isValid {
			return str
		}
	}
	return false
}

// safeGetString is a utility function to safely retrieve a string value from a slice of interface{} at the given index
func safeGetString(slice []interface{}, index int) string {
	if len(slice) > index {
		if str, isValid := slice[index].(string); isValid {
			return str
		}
	}
	return ""
}

// safeGetFloat64 is a utility function to safely retrieve a float64 value from a slice of interface{} at the given index
func safeGetFloat64(slice []interface{}, index int) float64 {
	if len(slice) > index {
		if f, isValid := slice[index].(float64); isValid {
			return f
		}
	}
	return 0.0
}

// safeGetInterfaceSlice is a utility function to safely retrieve a slice of interface{} from a slice of interface{} at the given index
func safeGetInterfaceSlice(slice []interface{}, index int) []interface{} {
	if len(slice) > index {
		if intfSlice, isValid := slice[index].([]interface{}); isValid {
			return intfSlice
		}
	}
	return nil
}

// parseOverallTranslation parses the overall translation result from Google Translate
func parseOverallTranslation(row []interface{}, result *TranslationResult) {
	result.TranslatedSentences = make([]Sentence, len(row))
	for i, rowItem := range row {
		if rowData, isValid := rowItem.([]interface{}); isValid {
			translation := safeGetString(rowData, 0)
			content := safeGetString(rowData, 1)
			pronunciation := safeGetString(rowData, 3)
			frequency := safeGetFloat64(rowData, 4)

			result.Translation += translation
			result.TranslatedSentences[i].Translation = translation

			result.Content += content
			result.TranslatedSentences[i].Content = content

			result.Pronunciation = pronunciation
			result.TranslatedSentences[i].Pronunciation = pronunciation

			result.TranslatedSentences[i].Frequency = frequency
		}
	}
}

// parseDetailedTranslation parses the detailed translation result for each word from Google Translate
func parseDetailedTranslation(row []interface{}, result *TranslationResult) {
	result.WordTranslations = make([]WordTranslation, len(row))
	for i, rowObj := range row {
		if rowData, isValid := rowObj.([]interface{}); isValid {
			partOfSentence := safeGetString(rowData, 0)
			translations := safeGetInterfaceSlice(rowData, 1)
			equivalents := safeGetInterfaceSlice(rowData, 2)
			frequency := safeGetFloat64(rowData, 4)

			result.WordTranslations[i].PartsOfSentence = partOfSentence

			for _, translation := range translations {
				if transStr, isString := translation.(string); isString {
					result.WordTranslations[i].Translations = append(result.WordTranslations[i].Translations, transStr)
				}
			}

			for _, equivalent := range equivalents {
				if equivArr, isArray := equivalent.([]interface{}); isArray {
					var equivalentDetail Equivalent
					equivalentDetail.Content = safeGetString(equivArr, 0)

					if equivalentWords, isWordsArray := equivArr[1].([]interface{}); isWordsArray {
						for _, equivWord := range equivalentWords {
							equivalentDetail.Equivalents = append(equivalentDetail.Equivalents, equivWord.(string))
						}
					}

					equivalentDetail.Frequency = safeGetFloat64(equivArr, 3)
					result.WordTranslations[i].Equivalents = append(result.WordTranslations[i].Equivalents, equivalentDetail)
				}
			}

			result.WordTranslations[i].Frequency = frequency
		}
	}
}

// parseAlternateTranslations parses the alternate translations for the text from Google Translate
func parseAlternateTranslations(row []interface{}, result *TranslationResult) {
	for _, item := range row {
		if rowData, isValid := item.([]interface{}); isValid {
			altTrans := AlternateTranslation{
				Content: safeGetString(rowData, 0),
			}
			altTransList := safeGetInterfaceSlice(rowData, 2)
			for _, alt := range altTransList {
				if transDetails, isValid := alt.([]interface{}); isValid {
					translation := Translation{
						Translation: safeGetString(transDetails, 0),
						IsCommon:    safeGetBool(transDetails, 2),
						IsInformal:  safeGetBool(transDetails, 3),
					}
					altTrans.Translations = append(altTrans.Translations, translation)
				}
			}
			result.AlternateTranslations = append(result.AlternateTranslations, altTrans)
		}
	}
}

// parseWordSynonyms parses the synonyms for the words in the text from Google Translate
func parseWordSynonyms(row []interface{}, translation *TranslationResult) {
	for _, row := range row {
		if rowData, isValid := row.([]interface{}); isValid {

			wordSynonym := WordSynonym{
				PartsOfSentence: safeGetString(rowData, 0),
				Contents:        safeGetString(rowData, 2),
				Frequency:       safeGetFloat64(rowData, 3),
			}

			nestedArray := safeGetInterfaceSlice(rowData, 1)
			for _, entry := range nestedArray {
				if arrayEntry, isValid := entry.([]interface{}); isValid {
					var synonymInfo Synonym

					firstSubEntry := safeGetInterfaceSlice(arrayEntry, 0)
					var synonymStrings []string
					for _, synonym := range firstSubEntry {
						if synonymStr, isValid := synonym.(string); isValid {
							synonymStrings = append(synonymStrings, synonymStr)
						}
					}
					synonymInfo.Synonyms = synonymStrings

					thirdSubEntry := safeGetInterfaceSlice(arrayEntry, 2)
					categorySubData := safeGetInterfaceSlice(thirdSubEntry, 0)
					synonymInfo.Category = safeGetString(categorySubData, 0)

					wordSynonym.Synonyms = append(wordSynonym.Synonyms, synonymInfo)
				}
			}
			translation.WordSynonyms = append(translation.WordSynonyms, wordSynonym)
		}
	}
}

// parseWordDefinitions parses the definitions for the words in the text from Google Translate
func parseWordDefinitions(row []interface{}, result *TranslationResult) {
	for _, row := range row {
		if rowData, isValid := row.([]interface{}); isValid {

			wordDefinition := WordDefinition{
				PartsOfSentence: safeGetString(rowData, 0),
			}

			nestedArray := safeGetInterfaceSlice(rowData, 1)
			for _, entry := range nestedArray {
				if arrayEntry, isValid := entry.([]interface{}); isValid {
					wordDefinition.Definitions = append(wordDefinition.Definitions, Definition{
						Definition: safeGetString(arrayEntry, 0),
						Example:    safeGetString(arrayEntry, 2),
					})
				}
			}
			result.WordDefinitions = append(result.WordDefinitions, wordDefinition)
		}
	}
}

// parseWordExamples parses the examples for the words in the text from Google Translate
func parseWordExamples(row []interface{}, result *TranslationResult) {
	for _, item := range row {
		if rowData, isValid := item.([]interface{}); isValid {
			for _, data := range rowData {
				if example, isValid := data.([]interface{}); isValid {
					result.WordExamples = append(result.WordExamples, safeGetString(example, 0))
				}
			}
		}
	}
}

// parseRow calls the appropriate parsing function based on the row index
func parseRow(row []interface{}, idx int, translationResult *TranslationResult) {
	if row, isValid := row[idx].([]interface{}); isValid {
		switch idx {
		case 0:
			parseOverallTranslation(row, translationResult)
		case 1:
			parseDetailedTranslation(row, translationResult)
		case 5:
			parseAlternateTranslations(row, translationResult)
		case 11:
			parseWordSynonyms(row, translationResult)
		case 12:
			parseWordDefinitions(row, translationResult)
		case 13:
			parseWordExamples(row, translationResult)
		default:
			// Handle other row indexes as needed
		}
	}
}

// extractTranslationData constructs a TranslationResult from the parsed JSON data
func extractTranslationData(jsonData []interface{}) *TranslationResult {
	var translationResult TranslationResult
	for rowIndex := range jsonData {
		parseRow(jsonData, rowIndex, &translationResult)
	}

	translationResult.SourceLanguage, _ = language.Parse(safeGetString(jsonData, 2))

	return &translationResult
}

// parseTranslationJSON takes the raw JSON data from Google Translate API and returns a TranslationResult structure
func parseTranslationJSON(jsonData []byte) (*TranslationResult, error) {
	var rawTranslationData []interface{}
	err := json.Unmarshal(jsonData, &rawTranslationData)
	if err != nil {
		return nil, errors.Join(errors.New("unable to parse the response from google translate api"), err)
	}

	return extractTranslationData(rawTranslationData), nil
}
