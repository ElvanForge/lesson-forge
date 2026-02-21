package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/unidoc/unioffice/presentation"
)

func GeneratePPTX(userID string, content string) (string, error) {
	ppt := presentation.New()
	defer ppt.Close()

	// Title Slide
	titleSlide := ppt.AddSlide()
	titleSlide.AddTextBox().AddParagraph().AddRun().SetText("ESL Lesson Plan")

	// Content Slide
	contentSlide := ppt.AddSlide()
	tf := contentSlide.AddTextBox()
	tf.AddParagraph().AddRun().SetText(content)

	// Save to a temporary or persistent location
	dir := "./output"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	filename := fmt.Sprintf("%s_%d.pptx", userID, os.Getpid())
	path := filepath.Join(dir, filename)

	if err := ppt.SaveToFile(path); err != nil {
		return "", err
	}

	return path, nil
}
