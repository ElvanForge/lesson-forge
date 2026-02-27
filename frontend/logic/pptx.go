package logic

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/presentation"
)

func GeneratePPTX(userID string, content string) ([]byte, string, error) {
	ppt := presentation.New()
	
	// Split slides by "---"
	slides := strings.Split(content, "---")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" || strings.Contains(cleanContent, "stripe.com") {
			continue
		}
		
		slide := ppt.AddSlide()

		// 1. STRIP MARKDOWN (Cleans up #, *, **, etc.)
		replacer := strings.NewReplacer("#", "", "*", "", "_", "", "**", "", "__", "")
		lines := strings.Split(cleanContent, "\n")
		
		var titleText string
		var bodyLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(replacer.Replace(line))
			if trimmed == "" { continue }
			if titleText == "" {
				titleText = trimmed
			} else {
				bodyLines = append(bodyLines, trimmed)
			}
		}

		// 2. DESIGN: THE SIDEBAR (Professional Blue-Gray)
		// This makes the slide look like a template rather than a blank page
		sidebar := slide.AddTextBox()
		sidebar.Properties().SetPosition(0, 0)
		sidebar.Properties().SetSize(2.2*measurement.Inch, 7.5*measurement.Inch)
		sidebar.Properties().SetSolidFill(color.SlateGray)

		// 3. THE TITLE (Right-aligned next to sidebar)
		titleTb := slide.AddTextBox()
		titleTb.Properties().SetPosition(2.5*measurement.Inch, 0.6*measurement.Inch)
		titleTb.Properties().SetSize(7.0*measurement.Inch, 1.2*measurement.Inch)
		
		titleP := titleTb.AddParagraph()
		run := titleP.AddRun()
		run.SetText(strings.ToUpper(titleText))
		run.Properties().SetSize(32)
		run.Properties().SetBold(true)
		run.Properties().SetSolidFill(color.SteelBlue)

		// 4. THE BODY (Clean list)
		if len(bodyLines) > 0 {
			bodyTb := slide.AddTextBox()
			bodyTb.Properties().SetPosition(2.5*measurement.Inch, 2.0*measurement.Inch)
			bodyTb.Properties().SetSize(7.0*measurement.Inch, 4.8*measurement.Inch)
			
			for _, line := range bodyLines {
				p := bodyTb.AddParagraph()
				p.Properties().SetLevel(0) // Indents text as a bullet point
				
				bodyRun := p.AddRun()
				bodyRun.SetText(line)
				bodyRun.Properties().SetSize(18)
				bodyRun.Properties().SetSolidFill(color.DimGray)
			}
		}
	}

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}