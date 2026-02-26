package logic

import (
	"bytes"
	"fmt"
	"time"

	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/presentation"
)

func GeneratePPTX(userID string, content string) ([]byte, string, error) {
	ppt := presentation.New()
	slide := ppt.AddSlide()
	tb := slide.AddTextBox()
	tb.Properties().SetPosition(0.5*measurement.Inch, 0.5*measurement.Inch)
	tb.AddParagraph().AddRun().SetText(content)

	var buf bytes.Buffer
	if err := ppt.Save(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("presentation_%s_%d.pptx", userID, time.Now().Unix()), nil
}