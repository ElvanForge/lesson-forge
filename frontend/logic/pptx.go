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
	
	// Split slides by the separator
	slides := strings.Split(content, "---")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" || strings.Contains(cleanContent, "stripe.com") { 
			continue 
		}
		
		slide := ppt.AddSlide()

		// --- 1. SET REAL BACKGROUND COLOR ---
		// We use the slide's background property correctly here
		bg := slide.Background()
		bg.Fill().SetSolidFill(color.SlateGray) // This forces the whole slide to NOT be white

		// --- 2. PARSE CONTENT (Remove Markdown) ---
		lines := strings.Split(cleanContent, "\n")
		var titleText string
		var bodyLines []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" { continue }
			// Strip all markdown characters
			line = strings.NewReplacer("#", "", "*", "", "_", "", "`", "").Replace(line)
			
			if titleText == "" {
				titleText = line
			} else {
				bodyLines = append(bodyLines, line)
			}
		}

		// --- 3. TITLE BOX ---
		titleTb := slide.AddTextBox()
		titleTb.Properties().SetPosition(0.5*measurement.Inch, 0.5*measurement.Inch)
		titleTb.Properties().SetSize(9*measurement.Inch, 1.5*measurement.Inch)
		
		titleP := titleTb.AddParagraph()
		titleP.Properties().SetAlign(dml.ST_TextAlignTypeCtr)
		run := titleP.AddRun()
		run.SetText(strings.ToUpper(titleText))
		run.Properties().SetSize(44)
		run.Properties().SetBold(true)
		run.Properties().SetSolidFill(color.White)

		// --- 4. BODY BOX ---
		if len(bodyLines) > 0 {
			bodyTb := slide.AddTextBox()
			bodyTb.Properties().SetPosition(1*measurement.Inch, 2.2*measurement.Inch)
			bodyTb.Properties().SetSize(8*measurement.Inch, 4.5*measurement.Inch)
			
			for _, line := range bodyLines {
				p := bodyTb.AddParagraph()
				p.Properties().SetLevel(0) // Adds bullet points
				
				bodyRun := p.AddRun()
				bodyRun.SetText(line)
				bodyRun.Properties().SetSize(24)
				bodyRun.Properties().SetSolidFill(color.LightGray)
			}
		}
	}

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}