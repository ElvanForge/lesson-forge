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
	
	// Split slides by "---" separator
	slides := strings.Split(content, "---")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" { continue }
		
		slide := ppt.AddSlide()

		// 1. CLEAN TEXT & REMOVE LINKS
		// Strip markdown symbols and filter out raw URLs/Stripe links
		replacer := strings.NewReplacer("#", "", "**", "", "__", "", "* ", "• ")
		lines := strings.Split(cleanContent, "\n")
		
		var titleText string
		var bodyLines []string
		for _, line := range lines {
			if strings.Contains(line, "stripe.com") || strings.Contains(line, "http") {
				continue
			}

			trimmed := strings.TrimSpace(replacer.Replace(line))
			if trimmed == "" { continue }
			
			if titleText == "" {
				titleText = trimmed
			} else {
				bodyLines = append(bodyLines, trimmed)
			}
		}

		// 2. DESIGN: THE SIDEBAR ACCENT
		// A vertical bar on the left makes the slide look professionally designed
		sidebar := slide.AddTextBox()
		sidebar.Properties().SetPosition(0, 0)
		sidebar.Properties().SetSize(0.3*measurement.Inch, 7.5*measurement.Inch)
		sidebar.Properties().SetSolidFill(color.SlateGray)

		// 3. THE TITLE
		titleTb := slide.AddTextBox()
		titleTb.Properties().SetPosition(0.6*measurement.Inch, 0.6*measurement.Inch)
		titleTb.Properties().SetSize(8.8*measurement.Inch, 1.2*measurement.Inch)
		
		titleP := titleTb.AddParagraph()
		run := titleP.AddRun()
		run.SetText(strings.ToUpper(titleText))
		
		// Correct font size syntax for gooxml build
		run.Properties().SetSize(34 * measurement.Point)
		run.Properties().SetBold(true)
		run.Properties().SetSolidFill(color.SteelBlue)

		// 4. THE BODY (With Dynamic Font Scaling)
		if len(bodyLines) > 0 {
			bodyTb := slide.AddTextBox()
			bodyTb.Properties().SetPosition(0.8*measurement.Inch, 2.0*measurement.Inch)
			bodyTb.Properties().SetSize(8.5*measurement.Inch, 4.8*measurement.Inch)
			
			// Adjust font size based on content density to prevent "Wall of Text"
			fontSize := 22.0
			if len(bodyLines) > 6 { fontSize = 18.0 }
			if len(bodyLines) > 9 { bodyLines = bodyLines[:9] } // Hard limit for legibility

			for _, line := range bodyLines {
				p := bodyTb.AddParagraph()
				p.Properties().SetLevel(0) 
				
				bodyRun := p.AddRun()
				bodyRun.SetText(line)
				bodyRun.Properties().SetSize(measurement.Distance(fontSize) * measurement.Point)
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