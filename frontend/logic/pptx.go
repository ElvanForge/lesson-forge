package logic

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/presentation"
)

func GeneratePPTX(userID string, content string) ([]byte, string, error) {
	ppt := presentation.New()
	
	// Split content into slides by the "#" character
	slides := strings.Split(content, "#")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" {
			continue
		}
		
		slide := ppt.AddSlide()
		tb := slide.AddTextBox()
		// Set position and size of the text box
		tb.Properties().SetPosition(0.5*measurement.Inch, 0.5*measurement.Inch)
		
		lines := strings.Split(cleanContent, "\n")
		for i, line := range lines {
			p := tb.AddParagraph()
			run := p.AddRun()
			run.SetText(strings.TrimSpace(line))
			
			// FIX: Access SetSize through Properties()
			if i == 0 {
				run.Properties().SetSize(32) // Title size
			} else {
				run.Properties().SetSize(18) // Body size
			}
		}
	}

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}