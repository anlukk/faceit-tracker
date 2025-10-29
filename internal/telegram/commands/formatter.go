package commands

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/anlukk/faceit-tracker/internal/faceit/pkg/go-faceit"
)

func formatPlayerCard(p *faceit.Player, matches []faceit.MatchStats) string {
	if p == nil {
		return "❌ Player not found"
	}

	var b strings.Builder

	flag := countryFlag(p.Country)
	if flag == "" {
		flag = "🌍"
	}
	verified := ""
	if p.Verified {
		verified = " ✅"
	}
	fmt.Fprintf(&b, "<b>%s</b> %s%s\n", p.Nickname, flag, verified)

	urlNick := url.PathEscape(p.Nickname)
	fmt.Fprintf(&b, "<a href=\"https://www.faceit.com/players/%s\">🔗 Faceit Profile</a>\n", urlNick)

	if !p.ActivatedAt.IsZero() {
		fmt.Fprintf(&b, "🕓 Joined: %s\n\n", p.ActivatedAt.Format("02 Jan 2006"))
	}

	if cs2, ok := p.Games["cs2"]; ok {
		fmt.Fprintf(&b, "<b>🎮 CS2:</b>\n")
		fmt.Fprintf(&b, "• 🧠 Skill level: <b>%d</b>\n", cs2.SkillLevel)
		if cs2.SkillLevelLabel != "" {
			fmt.Fprintf(&b, "• 🏆 Skill label: %s\n", cs2.SkillLevelLabel)
		}
		fmt.Fprintf(&b, "• 🧮 ELO: <b>%d</b>\n", cs2.FaceitElo)
		fmt.Fprintf(&b, "• 🧭 Region: %s\n", cs2.Region)
		if cs2.GamePlayerName != "" {
			fmt.Fprintf(&b, "• 🎮 In-game name: %s\n", cs2.GamePlayerName)
		}
		b.WriteString("\n")
	}

	if len(p.Memberships) > 0 {
		b.WriteString("<b>🏅 Membership:</b>\n")
		for _, m := range p.Memberships {
			fmt.Fprintf(&b, "• %s\n", strings.Title(m))
		}
		b.WriteString("\n")
	}

	if len(matches) > 0 {
		limit := 5
		if len(matches) < limit {
			limit = len(matches)
		}

		var wins, totalKills, totalDeaths int

		for i, m := range matches {
			if i >= limit {
				break
			}

			for _, round := range m.Rounds {
				for _, team := range round.Teams {
					for _, player := range team.Players {
						if strings.EqualFold(player.Nickname.(string), p.Nickname) {
							ps := player.PlayerStats
							k := atoi(getString(ps, "Kills"))
							d := atoi(getString(ps, "Deaths"))
							res := getString(ps, "Result")

							totalKills += k
							totalDeaths += d
							if res == "1" {
								wins++
							}
						}
					}
				}
			}

			if wins >= limit {
				break
			}
		}

		winrate := float64(wins) / float64(limit) * 100
		avgKills := float64(totalKills) / float64(limit)
		avgDeaths := float64(totalDeaths) / float64(limit)
		avgKD := float64(totalKills) / math.Max(1, float64(totalDeaths))

		b.WriteString("<b>📈 Summary (Last Matches):</b>\n")
		fmt.Fprintf(&b, "• ✅ Winrate: <b>%.1f%%</b> (%d/%d)\n", winrate, wins, limit)
		fmt.Fprintf(&b, "• 🔫 Avg Kills: <b>%.1f</b>\n", avgKills)
		fmt.Fprintf(&b, "• ☠️ Avg Deaths: <b>%.1f</b>\n", avgDeaths)
		fmt.Fprintf(&b, "• ⚖️ K/D Ratio: <b>%.2f</b>\n\n", avgKD)

		b.WriteString("<b>📊 Last 5 matches:</b>\n")

		for i, m := range matches {
			if i >= 5 {
				break
			}

			for _, round := range m.Rounds {
				mapName := getString(round.RoundStats, "Map")
				team1Score := getString(round.Teams[0].TeamStats, "Final Score")
				team2Score := getString(round.Teams[1].TeamStats, "Final Score")

				for _, team := range round.Teams {
					for _, player := range team.Players {
						if strings.EqualFold(player.Nickname.(string), p.Nickname) {
							ps := player.PlayerStats

							kills := getString(ps, "Kills")
							deaths := getString(ps, "Deaths")
							kd := getString(ps, "K/D Ratio")
							adr := getString(ps, "ADR")

							hs := getString(ps, "Headshots")
							hsPerc := getString(ps, "Headshots %")
							mvps := getString(ps, "MVPs")

							flashCount := getString(ps, "Flash Count")
							flashSuccess := getString(ps, "Flash Successes")
							utilityDamage := getString(ps, "Utility Damage")
							entryWins := getString(ps, "Entry Wins")
							entryCount := getString(ps, "Entry Count")

							clutch := getString(ps, "Clutch Kills")
							oneVOne := getString(ps, "1v1Wins")
							oneVTwo := getString(ps, "1v2Wins")

							result := getString(ps, "Result")
							switch result {
							case "1":
								result = "✅ Win"
							case "0":
								result = "❌ Loss"
							default:
								result = "❔ Unknown"
							}

							fmt.Fprintf(&b, "🗺 <b>%s</b> — %s | Score: %s:%s\n", mapName, result, team1Score, team2Score)
							fmt.Fprintf(&b, "• 💀 K/D: %s/%s | Ratio: %s | ADR: %s\n", kills, deaths, kd, adr)
							fmt.Fprintf(&b, "• 🎯 HS: %s (%s%%) | MVPs: %s\n", hs, hsPerc, mvps)
							fmt.Fprintf(&b, "• 💥 Utility: Flash %s/%s, Damage %s | Entry Wins %s/%s\n", flashSuccess, flashCount, utilityDamage, entryWins, entryCount)
							fmt.Fprintf(&b, "• 🏆 Clutch %s | 1v1 %s | 1v2 %s\n\n", clutch, oneVOne, oneVTwo)
						}
					}
				}
			}
		}
	}

	return b.String()
}

func formatSearchCommandPlayerCard(p *faceit.Player, matches []faceit.MatchStats) string {
	if p == nil {
		return "❌ Player not found"
	}

	var b strings.Builder

	flag := countryFlag(p.Country)
	if flag == "" {
		flag = "🌍"
	}

	verified := ""
	if p.Verified {
		verified = " ✅"
	}

	fmt.Fprintf(&b, "%s %s%s\n", p.Nickname, flag, verified)

	urlNick := url.PathEscape(p.Nickname)
	fmt.Fprintf(&b, "🔗 <a href=\"https://www.faceit.com/players/%s\">Faceit Profile</a>\n", urlNick)

	if cs2, ok := p.Games["cs2"]; ok {
		fmt.Fprintf(&b, "🎮 CS2: Lvl %d | ELO %d | 🌍 %s\n", cs2.SkillLevel, cs2.FaceitElo, strings.ToUpper(cs2.Region))
	}

	if len(matches) > 0 {
		total := len(matches)
		var wins, kills, deaths int

		for _, m := range matches {
			for _, round := range m.Rounds {
				for _, team := range round.Teams {
					for _, player := range team.Players {
						if strings.EqualFold(player.Nickname.(string), p.Nickname) {
							ps := player.PlayerStats
							k := atoi(getString(ps, "Kills"))
							d := atoi(getString(ps, "Deaths"))
							res := getString(ps, "Result")

							kills += k
							deaths += d
							if res == "1" {
								wins++
							}
						}
					}
				}
			}
		}

		winrate := float64(wins) / float64(total) * 100
		avgKills := float64(kills) / float64(total)
		avgKD := float64(kills) / math.Max(1, float64(deaths))

		fmt.Fprintf(&b, "📈 Winrate: %.0f%% | Avg K/D: %.2f | Avg Kills: %.0f", winrate, avgKD, avgKills)
	}

	return b.String()
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func getString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func countryFlag(code string) string {
	if len(code) != 2 {
		return ""
	}
	code = strings.ToUpper(code)
	runes := []rune(code)
	return string(rune(runes[0]-'A'+0x1F1E6)) + string(rune(runes[1]-'A'+0x1F1E6))
}
