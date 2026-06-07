package language

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

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

// Unicode ranges for Indian scripts — each script is unique, no overlap possible.
var scriptRanges = []struct {
	name  string
	start rune
	end   rune
}{
	{"kn", 0x0C80, 0x0CFF}, // Kannada
	{"ta", 0x0B80, 0x0BFF}, // Tamil
	{"te", 0x0C00, 0x0C7F}, // Telugu
	{"ml", 0x0D00, 0x0D7F}, // Malayalam
	{"hi", 0x0900, 0x097F}, // Devanagari (Hindi, Marathi, etc.)
}

var localityScores = map[string]float64{
	"kn": 1.0,
	"ta": 1.0,
	"te": 1.0,
	"ml": 1.0,
	"hi": 1.0,
	"en": 0.3,
}

// Detect identifies the language of text.
//
// Detection order:
// 1. Unicode script analysis — primary method, free, instant, 100% accurate for Indian scripts.
// 2. DeepL API — fallback only for pure-Latin text to distinguish English from Romanised Hindi etc.
// 3. Default to English (locality 0.3) if both methods fail.
func (s *Service) Detect(text string) (*DetectResult, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return &DetectResult{Language: "en", Confidence: 1.0, LocalityScore: 0.3}, nil
	}

	result := s.detectByScript(text)
	if result != nil {
		return result, nil
	}

	// Pure Latin — try DeepL for confirmation (handles Romanised Hindi / Hinglish)
	if config.App.DeepLAPIKey != "" {
		deepLResult, err := s.detectViaDeepL(text)
		if err == nil && deepLResult != nil {
			return deepLResult, nil
		}
	}

	// Fallback
	return &DetectResult{Language: "en", Confidence: 0.5, LocalityScore: 0.3}, nil
}

// detectByScript uses Unicode block analysis.
// Returns nil if text is pure Latin (no Indian script characters found).
func (s *Service) detectByScript(text string) *DetectResult {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	// Count characters per script
	scriptCount := make(map[string]int)
	latinCount := 0
	totalMeaningful := 0

	for _, r := range runes {
		if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsDigit(r) {
			continue
		}
		totalMeaningful++

		matched := false
		for _, sr := range scriptRanges {
			if r >= sr.start && r <= sr.end {
				scriptCount[sr.name]++
				matched = true
				break
			}
		}
		if !matched && r < 0x0300 { // Basic Latin + Latin Extended
			latinCount++
		}
	}

	if totalMeaningful == 0 {
		return nil
	}

	// Find dominant Indian script
	dominantLang := ""
	dominantCount := 0
	for lang, count := range scriptCount {
		if count > dominantCount {
			dominantCount = count
			dominantLang = lang
		}
	}

	indianTotal := totalMeaningful - latinCount
	indianRatio := float64(indianTotal) / float64(totalMeaningful)
	latinRatio := float64(latinCount) / float64(totalMeaningful)

	// Pure Indian script (>85% Indian chars)
	if dominantLang != "" && indianRatio >= 0.85 {
		confidence := float64(dominantCount) / float64(indianTotal)
		if confidence > 1.0 {
			confidence = 1.0
		}
		return &DetectResult{
			Language:      dominantLang,
			Confidence:    confidence,
			LocalityScore: 1.0,
		}
	}

	// Code-mixed: significant Indian script + significant Latin (>15% each)
	if dominantLang != "" && indianRatio >= 0.15 && latinRatio >= 0.15 {
		return &DetectResult{
			Language:      dominantLang + "-mixed",
			Confidence:    0.85,
			LocalityScore: 0.6,
		}
	}

	// Mostly Latin — return nil to trigger DeepL fallback
	if latinRatio > 0.85 {
		return nil
	}

	// Ambiguous — treat as code-mixed if any Indian script found
	if dominantLang != "" {
		return &DetectResult{
			Language:      dominantLang + "-mixed",
			Confidence:    0.6,
			LocalityScore: 0.6,
		}
	}

	return nil
}

// detectViaDeepL calls DeepL translate API to get detected_source_language.
// DeepL is only used for Latin text — it does NOT support Kannada/Tamil/Telugu/Malayalam.
// For Hindi in Devanagari, Unicode detection handles it. This covers Romanised/Hinglish edge cases.
func (s *Service) detectViaDeepL(text string) (*DetectResult, error) {
	type requestBody struct {
		Text       []string `json:"text"`
		TargetLang string   `json:"target_lang"`
	}

	payload, err := json.Marshal(requestBody{
		Text:       []string{text},
		TargetLang: "EN-US",
	})
	if err != nil {
		return nil, fmt.Errorf("deepl marshal: %w", err)
	}

	// Free tier uses api-free.deepl.com — keys end with :fx
	baseURL := "https://api-free.deepl.com"
	if !strings.HasSuffix(config.App.DeepLAPIKey, ":fx") {
		baseURL = "https://api.deepl.com"
	}

	req, err := http.NewRequest(http.MethodPost,
		baseURL+"/v2/translate",
		strings.NewReader(string(payload)),
	)
	if err != nil {
		return nil, fmt.Errorf("deepl request: %w", err)
	}
	req.Header.Set("Authorization", "DeepL-Auth-Key "+config.App.DeepLAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deepl call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deepl status: %d", resp.StatusCode)
	}

	var result struct {
		Translations []struct {
			DetectedSourceLanguage string `json:"detected_source_language"`
			Text                   string `json:"text"`
		} `json:"translations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("deepl decode: %w", err)
	}

	if len(result.Translations) == 0 {
		return nil, fmt.Errorf("deepl: empty response")
	}

	detectedLang := strings.ToLower(result.Translations[0].DetectedSourceLanguage)
	// DeepL returns "EN", "DE", "FR" etc — map to our ISO 639-1 codes
	lang := mapDeepLCode(detectedLang)
	locality := localityScoreForCode(lang)

	return &DetectResult{
		Language:      lang,
		Confidence:    0.90, // DeepL is accurate for Latin scripts
		LocalityScore: locality,
	}, nil
}

// mapDeepLCode normalises DeepL language codes to our app codes.
// DeepL returns codes like "EN", "EN-US", "PT-BR" — we only care about hi vs en.
func mapDeepLCode(code string) string {
	// DeepL uses "HI" for Hindi
	switch strings.ToUpper(strings.Split(code, "-")[0]) {
	case "HI":
		return "hi"
	case "EN":
		return "en"
	default:
		return "en" // treat any other Latin language as English-level locality
	}
}

func localityScoreForCode(lang string) float64 {
	if score, ok := localityScores[lang]; ok {
		return score
	}
	return 0.3
}

// LanguageCodeToName maps ISO code → display name used across the app.
func LanguageCodeToName(code string) string {
	m := map[string]string{
		"kn":       "kannada",
		"ta":       "tamil",
		"te":       "telugu",
		"ml":       "malayalam",
		"hi":       "hindi",
		"en":       "english",
		"kn-mixed": "kannada",
		"ta-mixed": "tamil",
		"te-mixed": "telugu",
		"ml-mixed": "malayalam",
		"hi-mixed": "hindi",
	}
	if name, ok := m[code]; ok {
		return name
	}
	return code
}
