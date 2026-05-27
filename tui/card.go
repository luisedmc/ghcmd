package tui

import (
	"fmt"
	"strings"

	"github.com/luisedmc/ghcmd/domain"
)

// RenderRepoCard returns a card for a repository search result
func RenderRepoCard(data domain.Repository, width int) string {
	var sb strings.Builder

	// Header -> 'repo name by owner'
	header := CardHeaderStyle.Render(data.Name) +
		CardOwnerStyle.Render("  by "+data.Owner)
	sb.WriteString(header + "\n")

	// Divider
	dividerWidth := 50
	if width > 0 && width-10 < dividerWidth {
		dividerWidth = width - 10
	}
	if dividerWidth < 20 {
		dividerWidth = 20
	}
	sb.WriteString(CardDividerStyle.Render(strings.Repeat("─", dividerWidth)) + "\n")

	// Description
	if data.Description != "" {
		sb.WriteString("\n" + CardDescStyle.Render(data.Description) + "\n")
	}

	// Stats row (stars, forks, open issues & language)
	stats := []string{
		CardStatStyle.Render("★ ") + CardValueStyle.Render(fmt.Sprintf("%d", data.Stars)),
		CardStatStyle.Render("⑂ ") + CardValueStyle.Render(fmt.Sprintf("%d", data.Forks)),
		CardStatStyle.Render("⚑ ") + CardValueStyle.Render(fmt.Sprintf("%d", data.OpenIssues)),
	}
	if data.Language != "" {
		stats = append(stats, CardStatStyle.Render("● ")+CardValueStyle.Render(data.Language))
	}
	sb.WriteString("\n" + strings.Join(stats, "   ") + "\n")

	// License & CreatedAt
	sb.WriteString("\n")
	if data.License != "" {
		sb.WriteString(CardLabelStyle.Render("License: ") + CardValueStyle.Render(data.License) + "\n")
	}

	if data.CreatedAt != "" {
		sb.WriteString(CardLabelStyle.Render("Created: ") + CardValueStyle.Render(data.CreatedAt) + "\n")
	}

	// URL
	sb.WriteString("\n" + CardURLStyle.Render(data.URL))

	// Card width
	cardWidth := 60
	if width > 0 && width-6 < cardWidth {
		cardWidth = width - 6
	}
	if cardWidth < 30 {
		cardWidth = 30
	}

	return CardStyle.Width(cardWidth).Render(sb.String())
}
