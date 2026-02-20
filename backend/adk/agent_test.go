package adk

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestNewCompetitorIntelligenceAgent tests the agent creation
func TestNewCompetitorIntelligenceAgent(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()

	if agent == nil {
		t.Fatal("Expected agent to be created, got nil")
	}

	if agent.Name != "CompetitorIntelligenceAgent" {
		t.Errorf("Expected agent name 'CompetitorIntelligenceAgent', got '%s'", agent.Name)
	}

	if agent.Description == "" {
		t.Error("Expected agent description to be non-empty")
	}
}

// TestCompetitorIntelligenceAgent_MarketResearch tests the market research functionality
func TestCompetitorIntelligenceAgent_MarketResearch(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	tests := []struct {
		name        string
		companyName string
		industry    string
		wantErr     bool
	}{
		{
			name:        "Valid SaaS industry",
			companyName: "TestCorp",
			industry:    "SaaS",
			wantErr:     false,
		},
		{
			name:        "Valid Fintech industry",
			companyName: "FinTech Inc",
			industry:    "Fintech",
			wantErr:     false,
		},
		{
			name:        "Empty company name",
			companyName: "",
			industry:    "Healthcare",
			wantErr:     false,
		},
		{
			name:        "Empty industry",
			companyName: "SomeCorp",
			industry:    "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			competitors, err := agent.MarketResearch(ctx, tt.companyName, tt.industry)

			if (err != nil) != tt.wantErr {
				t.Errorf("MarketResearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Verify we get exactly 3 competitors
			if len(competitors) != 3 {
				t.Errorf("Expected 3 competitors, got %d", len(competitors))
			}

			// Verify each competitor has required fields
			for i, comp := range competitors {
				if comp.Name == "" {
					t.Errorf("Competitor %d: Name is empty", i)
				}
				if comp.Website == "" {
					t.Errorf("Competitor %d: Website is empty", i)
				}
				if comp.Industry != tt.industry {
					t.Errorf("Competitor %d: Industry = %s, want %s", i, comp.Industry, tt.industry)
				}
				if len(comp.Products) == 0 {
					t.Errorf("Competitor %d: Products list is empty", i)
				}
				if comp.Pricing == "" {
					t.Errorf("Competitor %d: Pricing is empty", i)
				}
				if comp.MarketShare <= 0 {
					t.Errorf("Competitor %d: MarketShare should be positive, got %f", i, comp.MarketShare)
				}
				if len(comp.Strengths) == 0 {
					t.Errorf("Competitor %d: Strengths list is empty", i)
				}
				if len(comp.Weaknesses) == 0 {
					t.Errorf("Competitor %d: Weaknesses list is empty", i)
				}
			}
		})
	}
}

// TestCompetitorIntelligenceAgent_Analyze tests the analysis functionality
func TestCompetitorIntelligenceAgent_Analyze(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	tests := []struct {
		name    string
		data    []CompetitorData
		wantErr bool
	}{
		{
			name: "Multiple competitors with different market shares",
			data: []CompetitorData{
				{
					Name:        "High Share Corp",
					Pricing:     "Premium",
					MarketShare: 25.0,
					Strengths:   []string{"Brand", "Innovation"},
					Weaknesses:  []string{"Price"},
				},
				{
					Name:        "Mid Share Corp",
					Pricing:     "Mid-range",
					MarketShare: 15.0,
					Strengths:   []string{"UX", "Support"},
					Weaknesses:  []string{"Features", "Brand"},
				},
				{
					Name:        "Low Share Corp",
					Pricing:     "Budget",
					MarketShare: 5.0,
					Strengths:   []string{"Price"},
					Weaknesses:  []string{"Quality", "Support", "Features"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty competitor list",
			data:    []CompetitorData{},
			wantErr: false,
		},
		{
			name: "Single competitor - Enterprise pricing",
			data: []CompetitorData{
				{
					Name:        "Enterprise Corp",
					Pricing:     "Enterprise",
					MarketShare: 12.0,
					Strengths:   []string{"Security", "Compliance"},
					Weaknesses:  []string{"Complexity"},
				},
			},
			wantErr: false,
		},
		{
			name: "Single competitor - Unknown pricing",
			data: []CompetitorData{
				{
					Name:        "Unknown Corp",
					Pricing:     "Custom",
					MarketShare: 8.0,
					Strengths:   []string{"Flexibility"},
					Weaknesses:  []string{"Transparency"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyses, err := agent.Analyze(ctx, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Verify we get the same number of analyses as input competitors
			if len(analyses) != len(tt.data) {
				t.Errorf("Expected %d analyses, got %d", len(tt.data), len(analyses))
			}

			// Verify each analysis
			for i, analysis := range analyses {
				competitor := tt.data[i]

				if analysis.CompetitorName != competitor.Name {
					t.Errorf("Analysis %d: CompetitorName = %s, want %s", i, analysis.CompetitorName, competitor.Name)
				}

				// Verify threat level based on market share
				expectedThreatLevel := getExpectedThreatLevel(competitor.MarketShare)
				if analysis.ThreatLevel != expectedThreatLevel {
					t.Errorf("Analysis %d: ThreatLevel = %s, want %s", i, analysis.ThreatLevel, expectedThreatLevel)
				}

				// Verify positioning based on pricing
				expectedPositioning := getExpectedPositioning(competitor.Pricing)
				if analysis.Positioning != expectedPositioning {
					t.Errorf("Analysis %d: Positioning = %s, want %s", i, analysis.Positioning, expectedPositioning)
				}

				// Verify key differentiators match strengths
				if len(analysis.KeyDifferentiators) != len(competitor.Strengths) {
					t.Errorf("Analysis %d: KeyDifferentiators length = %d, want %d", i, len(analysis.KeyDifferentiators), len(competitor.Strengths))
				}

				// Verify opportunities are generated from weaknesses
				if len(analysis.Opportunities) != len(competitor.Weaknesses) {
					t.Errorf("Analysis %d: Opportunities length = %d, want %d", i, len(analysis.Opportunities), len(competitor.Weaknesses))
				}

				// Verify risks are generated from strengths
				if len(analysis.Risks) != len(competitor.Strengths) {
					t.Errorf("Analysis %d: Risks length = %d, want %d", i, len(analysis.Risks), len(competitor.Strengths))
				}

				// Verify opportunity format
				for j, opp := range analysis.Opportunities {
					if !strings.Contains(opp, "Capitalize on") {
						t.Errorf("Analysis %d, Opportunity %d: expected to contain 'Capitalize on', got '%s'", i, j, opp)
					}
				}

				// Verify risk format
				for j, risk := range analysis.Risks {
					if !strings.Contains(risk, "Competitor's") {
						t.Errorf("Analysis %d, Risk %d: expected to contain 'Competitor's', got '%s'", i, j, risk)
					}
				}
			}
		})
	}
}

// Helper function to determine expected threat level
func getExpectedThreatLevel(marketShare float64) string {
	switch {
	case marketShare > 20:
		return "High"
	case marketShare > 10:
		return "Medium"
	default:
		return "Low"
	}
}

// Helper function to determine expected positioning
func getExpectedPositioning(pricing string) string {
	switch pricing {
	case "Premium":
		return "Premium market leader"
	case "Mid-range":
		return "Value-focused challenger"
	case "Enterprise":
		return "Enterprise specialist"
	default:
		return "Undifferentiated"
	}
}

// TestCompetitorIntelligenceAgent_GenerateReport tests the report generation
func TestCompetitorIntelligenceAgent_GenerateReport(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	tests := []struct {
		name          string
		targetCompany string
		analyses      []CompetitorAnalysis
		wantErr       bool
	}{
		{
			name:          "Full report with multiple competitors",
			targetCompany: "MyCompany",
			analyses: []CompetitorAnalysis{
				{
					CompetitorName:     "Competitor A",
					ThreatLevel:        "High",
					Positioning:        "Premium market leader",
					KeyDifferentiators: []string{"Brand", "Innovation"},
					Opportunities:      []string{"Capitalize on Price"},
					Risks:              []string{"Competitor's Brand advantage"},
				},
				{
					CompetitorName:     "Competitor B",
					ThreatLevel:        "Medium",
					Positioning:        "Value-focused challenger",
					KeyDifferentiators: []string{"UX", "Support"},
					Opportunities:      []string{"Capitalize on Features"},
					Risks:              []string{"Competitor's UX advantage"},
				},
			},
			wantErr: false,
		},
		{
			name:          "Report with empty analyses",
			targetCompany: "SoloCorp",
			analyses:      []CompetitorAnalysis{},
			wantErr:       false,
		},
		{
			name:          "Report with single competitor",
			targetCompany: "Startup Inc",
			analyses: []CompetitorAnalysis{
				{
					CompetitorName:     "BigCorp",
					ThreatLevel:        "High",
					Positioning:        "Enterprise specialist",
					KeyDifferentiators: []string{"Security"},
					Opportunities:      []string{"Capitalize on Complexity"},
					Risks:              []string{"Competitor's Security advantage"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := agent.GenerateReport(ctx, tt.targetCompany, tt.analyses)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Verify report fields
			if report.TargetCompany != tt.targetCompany {
				t.Errorf("Report.TargetCompany = %s, want %s", report.TargetCompany, tt.targetCompany)
			}

			// Verify timestamp is set and recent
			if report.GeneratedAt.IsZero() {
				t.Error("Report.GeneratedAt is zero")
			}
			timeDiff := time.Since(report.GeneratedAt)
			if timeDiff > time.Minute {
				t.Errorf("Report.GeneratedAt is too old: %v", timeDiff)
			}

			// Verify competitors match analyses
			if len(report.Competitors) != len(tt.analyses) {
				t.Errorf("Report.Competitors length = %d, want %d", len(report.Competitors), len(tt.analyses))
			}

			// Verify market insights is non-empty
			if report.MarketInsights == "" {
				t.Error("Report.MarketInsights is empty")
			}

			// Verify recommendations are present
			if len(report.Recommendations) == 0 {
				t.Error("Report.Recommendations is empty")
			}

			// Verify specific recommendations exist
			expectedRecommendations := []string{
				"Focus on differentiation",
				"Target mid-market segment",
				"Invest in customer support",
				"Develop integrations",
				"Monitor competitor pricing",
			}

			for _, expected := range expectedRecommendations {
				found := false
				for _, rec := range report.Recommendations {
					if strings.Contains(rec, expected) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected recommendation containing '%s' not found", expected)
				}
			}
		})
	}
}

// TestCompetitorIntelligenceAgent_Run tests the full workflow
func TestCompetitorIntelligenceAgent_Run(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	tests := []struct {
		name        string
		companyName string
		industry    string
		wantErr     bool
	}{
		{
			name:        "Complete workflow - SaaS",
			companyName: "TechStartup",
			industry:    "SaaS",
			wantErr:     false,
		},
		{
			name:        "Complete workflow - Fintech",
			companyName: "FinanceCo",
			industry:    "Fintech",
			wantErr:     false,
		},
		{
			name:        "Complete workflow - Healthcare",
			companyName: "HealthTech",
			industry:    "Healthcare",
			wantErr:     false,
		},
		{
			name:        "Empty inputs",
			companyName: "",
			industry:    "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := agent.Run(ctx, tt.companyName, tt.industry)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Verify report structure
			if report == nil {
				t.Fatal("Expected report to be non-nil")
			}

			// Verify target company
			if report.TargetCompany != tt.companyName {
				t.Errorf("Report.TargetCompany = %s, want %s", report.TargetCompany, tt.companyName)
			}

			// Verify we have competitors
			if len(report.Competitors) == 0 {
				t.Error("Expected competitors in report")
			}

			// Verify we have 3 competitors from market research
			if len(report.Competitors) != 3 {
				t.Errorf("Expected 3 competitors, got %d", len(report.Competitors))
			}

			// Verify market insights
			if report.MarketInsights == "" {
				t.Error("Expected market insights in report")
			}

			// Verify recommendations
			if len(report.Recommendations) == 0 {
				t.Error("Expected recommendations in report")
			}

			// Verify timestamp
			if report.GeneratedAt.IsZero() {
				t.Error("Expected GeneratedAt to be set")
			}
		})
	}
}

// TestCompetitorReport_ToJSON tests the JSON serialization
func TestCompetitorReport_ToJSON(t *testing.T) {
	tests := []struct {
		name   string
		report *CompetitorReport
		wantErr bool
	}{
		{
			name: "Valid report with all fields",
			report: &CompetitorReport{
				GeneratedAt:   time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				TargetCompany: "TestCorp",
				Competitors: []CompetitorAnalysis{
					{
						CompetitorName:     "Competitor A",
						ThreatLevel:        "High",
						Positioning:        "Premium market leader",
						KeyDifferentiators: []string{"Brand", "Innovation"},
						Opportunities:      []string{"Capitalize on Price"},
						Risks:              []string{"Competitor's Brand advantage"},
					},
				},
				MarketInsights:  "Market is competitive",
				Recommendations: []string{"Differentiate", "Invest in support"},
			},
			wantErr: false,
		},
		{
			name: "Empty report",
			report: &CompetitorReport{
				GeneratedAt:     time.Now(),
				TargetCompany:   "",
				Competitors:     []CompetitorAnalysis{},
				MarketInsights:  "",
				Recommendations: []string{},
			},
			wantErr: false,
		},
		{
			name: "Report with nil slices",
			report: &CompetitorReport{
				GeneratedAt:     time.Now(),
				TargetCompany:   "TestCorp",
				Competitors:     nil,
				MarketInsights:  "Test",
				Recommendations: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := tt.report.ToJSON()

			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Verify JSON is valid
			var decoded CompetitorReport
			if err := json.Unmarshal(jsonData, &decoded); err != nil {
				t.Errorf("Failed to unmarshal JSON: %v", err)
			}

			// Verify JSON contains expected fields
			jsonStr := string(jsonData)

			if !strings.Contains(jsonStr, tt.report.TargetCompany) && tt.report.TargetCompany != "" {
				t.Error("JSON does not contain target company")
			}

			if !strings.Contains(jsonStr, "generated_at") {
				t.Error("JSON does not contain generated_at field")
			}

			if !strings.Contains(jsonStr, "competitors") {
				t.Error("JSON does not contain competitors field")
			}

			// Verify JSON is properly indented
			if !strings.Contains(jsonStr, "\n") {
				t.Error("JSON should be indented (contain newlines)")
			}
		})
	}
}

// TestCompetitorData_Struct tests the CompetitorData struct
func TestCompetitorData_Struct(t *testing.T) {
	data := CompetitorData{
		Name:        "Test Corp",
		Website:     "https://test.com",
		Industry:    "SaaS",
		Products:    []string{"Product A", "Product B"},
		Pricing:     "Premium",
		MarketShare: 15.5,
		Strengths:   []string{"Innovation", "Support"},
		Weaknesses:  []string{"Price"},
	}

	// Verify all fields are accessible
	if data.Name != "Test Corp" {
		t.Error("Name field incorrect")
	}
	if data.Website != "https://test.com" {
		t.Error("Website field incorrect")
	}
	if data.Industry != "SaaS" {
		t.Error("Industry field incorrect")
	}
	if len(data.Products) != 2 {
		t.Error("Products field incorrect")
	}
	if data.Pricing != "Premium" {
		t.Error("Pricing field incorrect")
	}
	if data.MarketShare != 15.5 {
		t.Error("MarketShare field incorrect")
	}
	if len(data.Strengths) != 2 {
		t.Error("Strengths field incorrect")
	}
	if len(data.Weaknesses) != 1 {
		t.Error("Weaknesses field incorrect")
	}
}

// TestCompetitorAnalysis_Struct tests the CompetitorAnalysis struct
func TestCompetitorAnalysis_Struct(t *testing.T) {
	analysis := CompetitorAnalysis{
		CompetitorName:     "Test Corp",
		ThreatLevel:        "High",
		Positioning:        "Premium market leader",
		KeyDifferentiators: []string{"Brand", "Innovation"},
		Opportunities:      []string{"Capitalize on Price"},
		Risks:              []string{"Competitor's Brand advantage"},
	}

	// Verify all fields are accessible
	if analysis.CompetitorName != "Test Corp" {
		t.Error("CompetitorName field incorrect")
	}
	if analysis.ThreatLevel != "High" {
		t.Error("ThreatLevel field incorrect")
	}
	if analysis.Positioning != "Premium market leader" {
		t.Error("Positioning field incorrect")
	}
	if len(analysis.KeyDifferentiators) != 2 {
		t.Error("KeyDifferentiators field incorrect")
	}
	if len(analysis.Opportunities) != 1 {
		t.Error("Opportunities field incorrect")
	}
	if len(analysis.Risks) != 1 {
		t.Error("Risks field incorrect")
	}
}

// TestCompetitorReport_Struct tests the CompetitorReport struct
func TestCompetitorReport_Struct(t *testing.T) {
	now := time.Now()
	report := CompetitorReport{
		GeneratedAt:   now,
		TargetCompany: "MyCorp",
		Competitors: []CompetitorAnalysis{
			{
				CompetitorName: "Competitor A",
				ThreatLevel:    "High",
			},
		},
		MarketInsights:  "Market is competitive",
		Recommendations: []string{"Differentiate"},
	}

	// Verify all fields are accessible
	if !report.GeneratedAt.Equal(now) {
		t.Error("GeneratedAt field incorrect")
	}
	if report.TargetCompany != "MyCorp" {
		t.Error("TargetCompany field incorrect")
	}
	if len(report.Competitors) != 1 {
		t.Error("Competitors field incorrect")
	}
	if report.MarketInsights != "Market is competitive" {
		t.Error("MarketInsights field incorrect")
	}
	if len(report.Recommendations) != 1 {
		t.Error("Recommendations field incorrect")
	}
}

// TestMarketResearch_ContextCancellation tests context handling
func TestMarketResearch_ContextCancellation(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context before call
	cancel()

	// The current implementation doesn't check context cancellation,
	// but this test ensures it doesn't panic
	_, err := agent.MarketResearch(ctx, "TestCorp", "SaaS")

	// Current implementation doesn't return error on cancelled context
	// This documents the current behavior
	if err != nil {
		t.Logf("MarketResearch returned error with cancelled context: %v", err)
	}
}

// TestAnalyze_ContextCancellation tests context handling in Analyze
func TestAnalyze_ContextCancellation(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx, cancel := context.WithCancel(context.Background())

	data := []CompetitorData{
		{
			Name:        "Test Corp",
			Pricing:     "Premium",
			MarketShare: 20.0,
			Strengths:   []string{"Brand"},
			Weaknesses:  []string{"Price"},
		},
	}

	// Cancel context before call
	cancel()

	// The current implementation doesn't check context cancellation
	_, err := agent.Analyze(ctx, data)

	// Current implementation doesn't return error on cancelled context
	if err != nil {
		t.Logf("Analyze returned error with cancelled context: %v", err)
	}
}

// TestGenerateReport_ContextCancellation tests context handling in GenerateReport
func TestGenerateReport_ContextCancellation(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx, cancel := context.WithCancel(context.Background())

	analyses := []CompetitorAnalysis{
		{
			CompetitorName: "Test Corp",
			ThreatLevel:    "High",
		},
	}

	// Cancel context before call
	cancel()

	// The current implementation doesn't check context cancellation
	_, err := agent.GenerateReport(ctx, "MyCorp", analyses)

	// Current implementation doesn't return error on cancelled context
	if err != nil {
		t.Logf("GenerateReport returned error with cancelled context: %v", err)
	}
}

// TestRun_ContextCancellation tests context handling in Run
func TestRun_ContextCancellation(t *testing.T) {
	agent := NewCompetitorIntelligenceAgent()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context before call
	cancel()

	// The current implementation doesn't check context cancellation
	_, err := agent.Run(ctx, "TestCorp", "SaaS")

	// Current implementation doesn't return error on cancelled context
	if err != nil {
		t.Logf("Run returned error with cancelled context: %v", err)
	}
}

// BenchmarkMarketResearch benchmarks the market research function
func BenchmarkMarketResearch(b *testing.B) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.MarketResearch(ctx, "TestCorp", "SaaS")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAnalyze benchmarks the analyze function
func BenchmarkAnalyze(b *testing.B) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	data := []CompetitorData{
		{
			Name:        "Competitor A",
			Pricing:     "Premium",
			MarketShare: 25.0,
			Strengths:   []string{"Brand", "Innovation", "Support"},
			Weaknesses:  []string{"Price", "Complexity"},
		},
		{
			Name:        "Competitor B",
			Pricing:     "Mid-range",
			MarketShare: 15.0,
			Strengths:   []string{"UX", "Price"},
			Weaknesses:  []string{"Features"},
		},
		{
			Name:        "Competitor C",
			Pricing:     "Enterprise",
			MarketShare: 10.0,
			Strengths:   []string{"Security"},
			Weaknesses:  []string{"Price", "Complexity", "Setup"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.Analyze(ctx, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerateReport benchmarks the report generation function
func BenchmarkGenerateReport(b *testing.B) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	analyses := []CompetitorAnalysis{
		{
			CompetitorName:     "Competitor A",
			ThreatLevel:        "High",
			Positioning:        "Premium market leader",
			KeyDifferentiators: []string{"Brand", "Innovation"},
			Opportunities:      []string{"Capitalize on Price"},
			Risks:              []string{"Competitor's Brand advantage"},
		},
		{
			CompetitorName:     "Competitor B",
			ThreatLevel:        "Medium",
			Positioning:        "Value-focused challenger",
			KeyDifferentiators: []string{"UX"},
			Opportunities:      []string{"Capitalize on Features"},
			Risks:              []string{"Competitor's UX advantage"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.GenerateReport(ctx, "TestCorp", analyses)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRun benchmarks the full workflow
func BenchmarkRun(b *testing.B) {
	agent := NewCompetitorIntelligenceAgent()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.Run(ctx, "TestCorp", "SaaS")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkToJSON benchmarks the JSON serialization
func BenchmarkToJSON(b *testing.B) {
	report := &CompetitorReport{
		GeneratedAt:   time.Now(),
		TargetCompany: "TestCorp",
		Competitors: []CompetitorAnalysis{
			{
				CompetitorName:     "Competitor A",
				ThreatLevel:        "High",
				Positioning:        "Premium market leader",
				KeyDifferentiators: []string{"Brand", "Innovation", "Support"},
				Opportunities:      []string{"Capitalize on Price", "Capitalize on Complexity"},
				Risks:              []string{"Competitor's Brand advantage", "Competitor's Innovation advantage"},
			},
			{
				CompetitorName:     "Competitor B",
				ThreatLevel:        "Medium",
				Positioning:        "Value-focused challenger",
				KeyDifferentiators: []string{"UX", "Price"},
				Opportunities:      []string{"Capitalize on Features"},
				Risks:              []string{"Competitor's UX advantage", "Competitor's Price advantage"},
			},
		},
		MarketInsights:  "The competitive landscape shows 2 major players.",
		Recommendations: []string{"Differentiate", "Invest in support", "Target mid-market"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := report.ToJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
}
