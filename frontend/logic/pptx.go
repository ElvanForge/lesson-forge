package logic

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/presentation"
	"baliance.com/gooxml/schema/soo/dml"
)

func GeneratePPTX(userID string, content string) ([]byte, string, error) {
	ppt := presentation.New()
	
	// Split by slide separator
	slides := strings.Split(content, "---")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" {
			continue
		}
		
		slide := ppt.AddSlide()

		// --- 1. THE CANVAS (Background) ---
		// We create a box that covers the WHOLE slide to simulate a background color
		bgCanvas := slide.AddTextBox()
		bgCanvas.Properties().SetPosition(0, 0)
		bgCanvas.Properties().SetSize(10*measurement.Inch, 7.5*measurement.Inch)
		bgCanvas.Properties().SetSolidFill(color.Gray) // A nice Slate/Gray base

		// --- 2. CLEANING & PARSING ---
		lines := strings.Split(cleanContent, "\n")
		var titleText string
		var bodyLines []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" { continue }

			// Remove Markdown syntax
			line = strings.ReplaceAll(line, "#", "")
			line = strings.ReplaceAll(line, "**", "")
			line = strings.ReplaceAll(line, "__", "")
			line = strings.TrimSpace(line)

			if titleText == "" {
				titleText = line
			} else {
				bodyLines = append(bodyLines, line)
			}
		}

		// --- 3. STYLIZED TITLE AREA ---
		// A dark header bar at the top
		headerBar := slide.AddTextBox()
		headerBar.Properties().SetPosition(0, 0)
		headerBar.Properties().SetSize(10*measurement.Inch, 1.5*measurement.Inch)
		headerBar.Properties().SetSolidFill(color.DarkBlue)

		titleP := headerBar.AddParagraph()
		titleP.Properties().SetAlign(dml.ST_TextAlignTypeCtr)
		run := titleP.AddRun()
		run.SetText(strings.ToUpper(titleText))
		run.Properties().SetSize(40)
		run.Properties().SetBold(true)
		run.Properties().SetSolidFill(color.White)

		// --- 4. STYLIZED BODY AREA ---
		if len(bodyLines) > 0 {
			bodyTb := slide.AddTextBox()
			bodyTb.Properties().SetPosition(0.75*measurement.Inch, 2.0*measurement.Inch)
			bodyTb.Properties().SetSize(8.5*measurement.Inch, 5.0*measurement.Inch)
			
			for _, line := range bodyLines {
				// Clean bullet characters
				text := strings.TrimPrefix(line, "*")
				text = strings.TrimPrefix(text, "-")
				text = strings.TrimSpace(text)

				p := bodyTb.AddParagraph()
				p.Properties().SetLevel(0) // Indent as a bullet
				
				bodyRun := p.AddRun()
				bodyRun.SetText(text)
				bodyRun.Properties().SetSize(22)
				bodyRun.Properties().SetSolidFill(color.White)
			}
		}
	}

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}