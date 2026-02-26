package main

import (
	"fmt"
	"strings"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func GeneratePDF(userID string, content string) ([]byte, string, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	headerBg := color.Color{Red: 245, Green: 247, Blue: 250}
	darkBlue := color.Color{Red: 44, Green: 62, Blue: 80}

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("VAELIA FORGE - LESSON PLAN", props.Text{
				Top:   5,
				Size:  16,
				Style: consts.Bold,
				Align: consts.Center,
				Color: darkBlue,
			})
		})
	})

	lines := strings.Split(content, "\n")
	var metadata = make(map[string]string)
	var bodyLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, ":") && (strings.HasPrefix(trimmed, "TOPIC") || strings.HasPrefix(trimmed, "LEVEL") || strings.HasPrefix(trimmed, "SUBJECT")) {
			parts := strings.SplitN(trimmed, ":", 2)
			metadata[strings.ToUpper(parts[0])] = strings.TrimSpace(parts[1])
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	m.Row(12, func() {
		m.Col(4, func() { m.Text("Subject: "+metadata["SUBJECT"], props.Text{Size: 10, Style: consts.Bold}) })
		m.Col(4, func() { m.Text("Topic: "+metadata["TOPIC"], props.Text{Size: 10, Style: consts.Bold}) })
		m.Col(4, func() { m.Text("Level: "+metadata["LEVEL"], props.Text{Size: 10, Style: consts.Bold}) })
	})
	m.Line(1.0)

	var tableRows [][]string
	inStructure := false

	for _, line := range bodyLines {
		line = strings.TrimSpace(line)
		if line == "" { continue }

		if strings.HasPrefix(strings.ToUpper(line), "STRUCTURE") {
			inStructure = true
			continue
		}

		if inStructure {
			parts := strings.Split(line, "|")
			if len(parts) >= 2 {
				row := []string{}
				for _, p := range parts { row = append(row, strings.TrimSpace(p)) }
				tableRows = append(tableRows, row)
				continue
			} else { inStructure = false }
		}

		if strings.HasPrefix(line, "##") {
			m.Row(12, func() {
				m.Col(12, func() { m.Text(strings.TrimSpace(strings.TrimPrefix(line, "##")), props.Text{Size: 12, Style: consts.Bold, Top: 5}) })
			})
		} else {
			m.Row(8, func() { m.Col(12, func() { m.Text(line, props.Text{Size: 10}) }) })
		}
	}

	if len(tableRows) > 0 {
		m.Row(10, func() { m.Col(12, func() { m.Text("Lesson Structure:", props.Text{Size: 12, Style: consts.Bold, Top: 10}) }) })
		header := []string{"Time", "Activities / Tasks", "Teaching Approach"}
		m.TableList(header, tableRows, props.TableList{
			HeaderProp: props.TableListContent{Size: 9, GridSizes: []uint{2, 6, 4}},
			ContentProp: props.TableListContent{Size: 9, GridSizes: []uint{2, 6, 4}},
			Align: consts.Left,
			AlternatedBackground: &headerBg,
		})
	}

	pdfBytes, err := m.Output()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("lesson_%s_%d.pdf", userID, SystemTimeNow())
	return pdfBytes.Bytes(), filename, nil
}

func SystemTimeNow() int64 {
	return 0 // Placeholder for timestamp logic handled in main
}