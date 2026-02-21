package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func GeneratePDF(userID string, content string) (string, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	// 1. Branding Header
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("VAELIA ESL LESSON PLAN", props.Text{
				Top:   5,
				Size:  20,
				Style: consts.Bold,
				Align: consts.Center,
				Color: color.Color{Red: 44, Green: 62, Blue: 80},
			})
		})
	})

	// 2. Parser Logic
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle Headers
		if strings.HasPrefix(line, "##") {
			headerText := strings.TrimSpace(strings.TrimPrefix(line, "##"))
			m.Row(12, func() {
				m.Col(12, func() {
					m.Text(headerText, props.Text{
						Size:  14,
						Style: consts.Bold,
						Align: consts.Left,
					})
				})
			})
			m.Line(1.0) // Decorative separator
		} else if strings.HasPrefix(line, "*") || strings.HasPrefix(line, "-") {
			// Handle Bullet Points
			bulletText := strings.TrimSpace(line[1:])
			m.Row(8, func() {
				m.Col(12, func() {
					m.Text("• "+bulletText, props.Text{
						Size:  11,
						Align: consts.Left,
						Left:  5,
					})
				})
			})
		} else {
			// Standard Paragraph
			m.Row(10, func() {
				m.Col(12, func() {
					m.Text(line, props.Text{
						Size:  11,
						Align: consts.Left,
					})
				})
			})
		}
	}

	// 3. File Persistence
	dir := "./output"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("lesson_%s_%d.pdf", userID, os.Getpid())
	path := filepath.Join(dir, filename)

	err := m.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	return path, nil
}