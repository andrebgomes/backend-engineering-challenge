// This package implements a CLI tool that calculates and outputs the
// moving average of translations delivery time.
package main

import (
	"backend-engineering-challenge/internal/translationdeliverytime"
	"flag"
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	// Parse the flags
	var inputFile *string
	var windowSize *int
	inputFile = flag.String("input_file", "", "file containing the stream events")
	windowSize = flag.Int("window_size", 10, "number of past minutes that will be handled")
	flag.Parse()

	// Create logger
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.DisableStacktrace = true
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalf("failed to create logger: %s", err)
	}
	sugaredLogger := logger.With(
		zap.String("input_file", *inputFile),
		zap.Int("window_size", *windowSize)).Sugar()

	// Read the input file
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		sugaredLogger.Error("reading input_file: ", err)
		return
	}

	// Create instance of TranslationDeliveryTimeApp and run it
	app := translationdeliverytime.NewTranslationDeliveryTimeApp(string(input), *windowSize)
	output, err := app.Run()
	if err != nil {
		sugaredLogger.Error(err)
		return
	}

	// Create and write the output file
	outputFile, err := os.Create("output")
	if err != nil {
		sugaredLogger.Error("creating output file: ", err)
		return
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(output)
	if err != nil {
		sugaredLogger.Error("writing output file: ", err)
	}
}
