package ai

import (
	"fmt"
	"strings"

	"GolangChatAdvisor_Project/internal/models"
)

type Responder struct{}

func NewResponder() *Responder {
	return &Responder{}
}

func (r *Responder) GenerateReply(companion models.Advisor, history []models.Message, userMessage string) string {
	lower := strings.ToLower(userMessage)

	switch {
	case strings.Contains(lower, "outfit"):
		return fmt.Sprintf("%s here — I’d suggest a bold jacket and clean accessories for your avatar.", companion.Name)
	case strings.Contains(lower, "color"):
		return fmt.Sprintf("%s here — try pairing one bright accent color with a neutral base.", companion.Name)
	case strings.Contains(lower, "style"):
		return fmt.Sprintf("%s here — given my %s style, I’d lean toward something polished but expressive.", companion.Name, companion.Style)
	default:
		return fmt.Sprintf("%s says: As a %s, I think that sounds like a fun direction to explore.", companion.Name, companion.Persona)
	}
}
