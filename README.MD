# GTranslate

A Go package for translation powered by Google Translate. This package provides various functionalities to extract comprehensive translation details.

## Features
* Translation of text from one language to another.
* Automatic detection of source language if not provided.
* Provides overall translations and detailed translations for each sentence.
* Provides synonyms, definitions, and examples for each word in the translated text.
* Provides alternate translations.

## Installation

To install GTranslate, run the following command in your terminal:

```bash
go get github.com/mshafiee/gtranslate
```

## Usage

Here is a simple example of using GTranslate:

```go
package main

import (
	"context"
	"fmt"
	"github.com/mshafiee/gtranslate"
	"golang.org/x/text/language"
	"log"
)

func main() {
	text := "Hello, world!"
	translate, err := gtranslate.Translate(context.Background(), text, language.Tag{}, language.Persian)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v", translate)
}
```

## Structure

The package includes several struct types:

* `Sentence`: Represents a translated sentence and related details.
* `Equivalent`: Represents equivalent translations.
* `WordTranslation`: Represents detailed translations of a word.
* `AlternateTranslation`: Represents alternative translations.
* `Translation`: Represents a translation.
* `Synonym`: Represents synonyms of a word.
* `WordSynonym`: Represents a word and its synonyms.
* `Definition`: Represents a definition of a word.
* `WordDefinition`: Represents a word and its definitions.
* `TranslationResult`: Represents the result of a translation request.

Each of these types contains various fields that represent different aspects of the translation.

## Contributing

Contributions to GTranslate are welcome! Please read the [contributing guidelines](CONTRIBUTING.md) before getting started.

## License

GTranslate is licensed under the [GNU license](LICENSE).
