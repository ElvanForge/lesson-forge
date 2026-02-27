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
	
	// Split slides by "---"
	slides := strings.Split(content, "---")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" || strings.Contains(cleanContent, "stripe.com") {
			continue
		}
		
		slide := ppt.AddSlide()

		// 1. CLEAN THE TEXT (Remove all Markdown garbage)
		// This ensures # and ** don't show up on your slides
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

		// 2. DESIGN: THE ACCENT BAR
		// We add a dark sidebar to give it a "Designed" look
		accentBar := slide.AddTextBox()
		accentBar.Properties().SetPosition(0, 0)
		accentBar.Properties().SetSize(2.5*measurement.Inch, 7.5*measurement.Inch)
		accentBar.Properties().SetSolidFill(color.SlateGray)

		// 3. THE TITLE (Large, Bold, Professional)
		titleTb := slide.AddTextBox()
		titleTb.Properties().SetPosition(2.8*measurement.Inch, 0.5*measurement.Inch)
		titleTb.Properties().SetSize(6.5*measurement.Inch, 1.5*measurement.Inch)
		
		titleP := titleTb.AddParagraph()
		run := titleP.AddRun()
		run.SetText(strings.ToUpper(titleText))
		run.Properties().SetSize(36)
		run.Properties().SetBold(true)
		run.Properties().SetSolidFill(color.DarkBlue)

		// 4. THE BODY (Clean Bullets)
		if len(bodyLines) > 0 {
			bodyTb := slide.AddTextBox()
			bodyTb.Properties().SetPosition(2.8*measurement.Inch, 2.2*measurement.Inch)
			bodyTb.Properties().SetSize(6.5*measurement.Inch, 4.5*measurement.Inch)
			
			for _, line := range bodyLines {
				p := bodyTb.AddParagraph()
				p.Properties().SetLevel(0) // This forces a bullet point indent
				
				bodyRun := p.AddRun()
				bodyRun.SetText(line)
				bodyRun.Properties().SetSize(20)
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