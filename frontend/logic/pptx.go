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
	
	// Split content into slides by the "---" or "#" character
	// The AI is prompted to use "# Slide Title"
	slides := strings.Split(content, "---")
	if len(slides) <= 1 {
		slides = strings.Split(content, "#")
	}
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" {
			continue
		}
		
		slide := ppt.AddSlide()
		
		// 1. CREATE A STYLIZED HEADER BAR (Faux Background)
		// This uses a shape since slide.Background can be buggy in some versions
		rect := slide.AddShape()
		rect.Properties().SetPosition(0, 0)
		rect.Properties().SetSize(10*measurement.Inch, 1.25*measurement.Inch)
		rect.Properties().SetSolidFill(color.SlateGray) 

		// Split title from the body
		parts := strings.SplitN(cleanContent, "\n", 2)
		titleText := strings.TrimPrefix(strings.TrimSpace(parts[0]), "# ")
		
		// 2. STYLIZED TITLE
		titleTb := slide.AddTextBox()
		titleTb.Properties().SetPosition(0.5*measurement.Inch, 0.25*measurement.Inch)
		titleTb.Properties().SetSize(9*measurement.Inch, 0.75*measurement.Inch)
		
		titleP := titleTb.AddParagraph()
		titleP.Properties().SetAlign(dml.ST_TextAlignTypeL)
		
		titleRun := titleP.AddRun()
		titleRun.SetText(strings.ToUpper(titleText))
		titleRun.Properties().SetSize(32)
		titleRun.Properties().SetBold(true)
		titleRun.Properties().SetSolidFill(color.White) // White text on SlateGray bar

		// 3. STYLIZED BODY
		if len(parts) > 1 {
			bodyTb := slide.AddTextBox()
			bodyTb.Properties().SetPosition(0.5*measurement.Inch, 1.75*measurement.Inch)
			bodyTb.Properties().SetSize(9*measurement.Inch, 5*measurement.Inch)
			
			lines := strings.Split(parts[1], "\n")
			for _, line := range lines {
				text := strings.TrimSpace(line)
				if text == "" {
					continue
				}
				
				p := bodyTb.AddParagraph()
				// Use Level to create automatic bullet point indentation
				if strings.HasPrefix(text, "*") || strings.HasPrefix(text, "-") {
					p.Properties().SetLevel(0)
					text = strings.TrimSpace(text[1:])
				}
				
				run := p.AddRun()
				run.SetText(text)
				run.Properties().SetSize(18)
				run.Properties().SetSolidFill(color.DimGray)
			}
		}
	}

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}