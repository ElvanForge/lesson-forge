package logic

import (
	"bytes"
	"fmt"
	"strings" // Required for splitting the text
	"time"

	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/presentation"
)

func GeneratePPTX(userID string, content string) ([]byte, string, error) {
	ppt := presentation.New()

	// Split content by a delimiter (e.g., "---" for slides)
	sections := strings.Split(content, "---")

	for _, section := range sections {
		trimmed := strings.TrimSpace(section)
		if trimmed == "" {
			continue
		}

		slide := ppt.AddSlide()
		tb := slide.AddTextBox()
		// Set position and size to prevent text from running off
		tb.Properties().SetPosition(0.5*measurement.Inch, 0.5*measurement.Inch)
		tb.Properties().SetSize(9*measurement.Inch, 6*measurement.Inch) 
		
		// Add the section text to the new slide
		tb.AddParagraph().AddRun().SetText(trimmed)
	}

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}