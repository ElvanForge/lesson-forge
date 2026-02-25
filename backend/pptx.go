package main

import (
	"bytes"
	"fmt"
	"strings"

	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/presentation"
)

func GeneratePPTX(userID string, content string) ([]byte, string, error) {
	ppt := presentation.New()

	slidesContent := strings.Split(content, "##")

	for _, section := range slidesContent {
		trimmed := strings.TrimSpace(section)
		if trimmed == "" { continue }

		slide := ppt.AddSlide()
		lines := strings.Split(trimmed, "\n")

		titleBox := slide.AddTextBox()
		titleBox.Properties().SetPosition(0.5*measurement.Inch, 0.5*measurement.Inch)
		titleBox.Properties().SetSize(9*measurement.Inch, 1*measurement.Inch)
		
		titlePara := titleBox.AddParagraph()
		titleRun := titlePara.AddRun()
		titleRun.SetText(lines[0])
		titleRun.Properties().SetSize(32)

		if len(lines) > 1 {
			bodyBox := slide.AddTextBox()
			bodyBox.Properties().SetPosition(0.5*measurement.Inch, 1.7*measurement.Inch)
			bodyBox.Properties().SetSize(9*measurement.Inch, 5*measurement.Inch)
			
			bodyPara := bodyBox.AddParagraph()
			bodyRun := bodyPara.AddRun()
			bodyText := strings.Join(lines[1:], "\n")
			bodyRun.SetText(bodyText)
			bodyRun.Properties().SetSize(18)
		}
	}

	var buf bytes.Buffer
	if err := ppt.Write(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("presentation_%s_%d.pptx", userID, SystemTimeNow())
	return buf.Bytes(), filename, nil
}