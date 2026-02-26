package logic

import (
	"fmt"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func GeneratePDF(userID string, content string) ([]byte, string, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	darkBlue := color.Color{Red: 44, Green: 62, Blue: 80}
	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("VAELIA FORGE", props.Text{Size: 16, Style: consts.Bold, Align: consts.Center, Color: darkBlue})
		})
	})
	m.Row(10, func() { m.Col(12, func() { m.Text(content, props.Text{Size: 10}) }) })
	pdfBytes, err := m.Output()
	if err != nil {
		return nil, "", err
	}
	return pdfBytes.Bytes(), fmt.Sprintf("lesson_%s_%d.pdf", userID, time.Now().Unix()), nil
}