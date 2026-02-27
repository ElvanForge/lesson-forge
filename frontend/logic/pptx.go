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
	
	// Split content into slides by the "---" separator often used by AI
	slides := strings.Split(content, "---")
	
	for _, slideContent := range slides {
		cleanContent := strings.TrimSpace(slideContent)
		if cleanContent == "" {
			continue
		}
		
		slide := ppt.AddSlide()
		
		// 1. CREATE A STYLIZED HEADER BOX
		// We use a TextBox with a background fill instead of AddShape
		headerBox := slide.AddTextBox()
		headerBox.Properties().SetPosition(0, 0)
		headerBox.Properties().SetSize(10*measurement.Inch, 1.25*measurement.Inch)
		headerBox.Properties().SetSolidFill(color.SlateGray) // Dark header background

		// Split title from the rest of the slide content
		parts := strings.SplitN(cleanContent, "\n", 2)
		titleText := strings.TrimPrefix(strings.TrimSpace(parts[0]), "# ")
		
		// Add the title text to the header box
		titleP := headerBox.AddParagraph()
		titleP.Properties().SetAlign(dml.ST_TextAlignTypeL)
		
		titleRun := titleP.AddRun()
		titleRun.SetText("  " + strings.ToUpper(titleText)) // Added padding
		titleRun.Properties().SetSize(32)
		titleRun.Properties().SetBold(true)
		titleRun.Properties().SetSolidFill(color.White)

		// 2. STYLIZED BODY BOX
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
				
				// Handle bullets manually if AI provides them
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