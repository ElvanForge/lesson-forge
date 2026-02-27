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
	
	// Split content into slides by the "#" character
	slides := strings.Split(content, "#")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" {
			continue
		}
		
		slide := ppt.AddSlide()
		
		// Fix: Set Background using the direct Background() method
		bg := slide.Background()
		bg.Fill().SetSolidFill(color.LightSlateGray) 

		// Split title from the rest of the slide content
		parts := strings.SplitN(cleanContent, "\n", 2)
		titleText := parts[0]
		
		// 1. Stylized Title Box
		titleTb := slide.AddTextBox()
		titleTb.Properties().SetPosition(0.5*measurement.Inch, 0.4*measurement.Inch)
		titleTb.Properties().SetSize(9*measurement.Inch, 1*measurement.Inch)
		
		titleP := titleTb.AddParagraph()
		titleP.Properties().SetAlign(dml.ST_TextAlignTypeCtr)
		
		titleRun := titleP.AddRun()
		titleRun.SetText(strings.ToUpper(titleText))
		titleRun.Properties().SetSize(36)
		titleRun.Properties().SetBold(true)
		titleRun.Properties().SetSolidFill(color.White)

		// 2. Stylized Body Box
		if len(parts) > 1 {
			bodyTb := slide.AddTextBox()
			bodyTb.Properties().SetPosition(0.75*measurement.Inch, 1.5*measurement.Inch)
			bodyTb.Properties().SetSize(8.5*measurement.Inch, 5*measurement.Inch)
			
			lines := strings.Split(parts[1], "\n")
			for _, line := range lines {
				text := strings.TrimSpace(line)
				if text == "" {
					continue
				}
				
				p := bodyTb.AddParagraph()
				// Adding indenting for bullet style
				p.Properties().SetLevel(0) 
				
				run := p.AddRun()
				run.SetText(text)
				run.Properties().SetSize(20)
				run.Properties().SetSolidFill(color.White)
			}
		}
	}

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}