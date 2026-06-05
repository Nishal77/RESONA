package language

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strings"
	"time"

	"github.com/Nishal77/resona/backend/pkg/config"
)

type DetectResult struct {
	Language      string  `json:"language"`
	Confidence    float64 `json:"confidence"`
	LocalityScore float64 `json:"locality_score"`
}

type Service struct {
	httpClient *http.Client
}

func NewService() *Service {
	return &Service{httpClient: &http.Client{Timeout: 5 * time.Second}}
}

var localityScores = map[string]float64{
	"kn": 1.0, // Kannada
	"ta": 1.0, // Tamil
	"te": 1.0, // Telugu
	"ml": 1.0, // Malayalam
	"hi": 1.0, // Hindi
	"en": 0.3, // English
}

func (s *Service) Detect(text string) (*DetectResult, error) {
	if strings.TrimSpace(text) == "" {
		return &DetectResult{Language: "en", Confidence: 1.0, LocalityScore: 0.3}, nil
	}

	apiURL := fmt.Sprintf(
		"https://translation.googleapis.com/language/translate/v2/detect?key=%s",
		config.App.GoogleTranslateAPIKey,
	)

	payload := fmt.Sprintf(`{"q":%s}`, mustJSON(text))
	resp, err := s.httpClient.Post(apiURL, "application/json", strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("translate api: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Detections [][]struct {
				Language   string  `json:"language"`
				Confidence float64 `json:"confidence"`
			} `json:"detections"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Data.Detections) == 0 || len(result.Data.Detections[0]) == 0 {
		return &DetectResult{Language: "en", Confidence: 0, LocalityScore: 0.3}, nil
	}

	det := result.Data.Detections[0][0]
	locality := localityScoreForCode(det.Language, det.Confidence)

	return &DetectResult{
		Language:      det.Language,
		Confidence:    det.Confidence,
		LocalityScore: locality,
	}, nil
}

func localityScoreForCode(lang string, confidence float64) float64 {
	// Code-mixed: detected as "en" but low confidence → 0.6
	if lang == "en" && confidence < 0.70 {
		return 0.6
	}
	if score, ok := localityScores[lang]; ok {
		return score
	}
	return 0.3
}

func mustJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// LanguageCodeToName maps ISO code → display name
func LanguageCodeToName(code string) string {
	m := map[string]string{
		"kn": "kannada",
		"ta": "tamil",
		"te": "telugu",
		"ml": "malayalam",
		"hi": "hindi",
		"en": "english",
	}
	if name, ok := m[code]; ok {
		return name
	}
	return code
}

