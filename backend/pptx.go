package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/presentation"
)

func GeneratePPTX(userID string, content string) (string, error) {
	ppt := presentation.New()

	// Split content by '##' per your existing architecture
	slidesContent := strings.Split(content, "##")

	for _, section := range slidesContent {
		trimmed := strings.TrimSpace(section)
		if trimmed == "" {
			continue
		}

		slide := ppt.AddSlide()
		lines := strings.Split(trimmed, "\n")

		// 1. Add Title Box
		titleBox := slide.AddTextBox()
		titleBox.Properties().SetPosition(0.5*measurement.Inch, 0.5*measurement.Inch)
		titleBox.Properties().SetSize(9*measurement.Inch, 1*measurement.Inch)
		
		titlePara := titleBox.AddParagraph()
		titleRun := titlePara.AddRun()
		
		// Set Title Text
		titleRun.SetText(lines[0])
		titleRun.Properties().SetSize(32)

		// 2. Add Body text if it exists
		if len(lines) > 1 {
			bodyBox := slide.AddTextBox()
			bodyBox.Properties().SetPosition(0.5*measurement.Inch, 1.7*measurement.Inch)
			bodyBox.Properties().SetSize(9*measurement.Inch, 5*measurement.Inch)
			
			bodyPara := bodyBox.AddParagraph()
			bodyRun := bodyPara.AddRun()
			
			// Join remaining lines as body text
			bodyText := strings.Join(lines[1:], "\n")
			bodyRun.SetText(bodyText)
			bodyRun.Properties().SetSize(18)
		}
	}

	dir := "./output"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("presentation_%s_%d.pptx", userID, os.Getpid())
	path := filepath.Join(dir, filename)

	if err := ppt.SaveToFile(path); err != nil {
		return "", err
	}

	return path, nil
}