package adk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CompetitorData represents raw competitor information
type CompetitorData struct {
	Name        string   `json:"name"`
	Website     string   `json:"website"`
	Industry    string   `json:"industry"`
	Products    []string `json:"products"`
	Pricing     string   `json:"pricing"`
	MarketShare float64  `json:"market_share"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
}

// CompetitorAnalysis represents analyzed competitive positioning
type CompetitorAnalysis struct {
	CompetitorName     string   `json:"competitor_name"`
	ThreatLevel        string   `json:"threat_level"`
	Positioning        string   `json:"positioning"`
	KeyDifferentiators []string `json:"key_differentiators"`
	Opportunities      []string `json:"opportunities"`
	Risks              []string `json:"risks"`
}

// CompetitorReport represents the final intelligence report
type CompetitorReport struct {
	GeneratedAt     time.Time            `json:"generated_at"`
	TargetCompany   string               `json:"target_company"`
	Competitors     []CompetitorAnalysis `json:"competitors"`
	MarketInsights  string               `json:"market_insights"`
	Recommendations []string             `json:"recommendations"`
}

// CompetitorIntelligenceAgent provides tools for competitor analysis
type CompetitorIntelligenceAgent struct {
	Name        string
	Description string
}

// NewCompetitorIntelligenceAgent creates a new agent instance
func NewCompetitorIntelligenceAgent() *CompetitorIntelligenceAgent {
	return &CompetitorIntelligenceAgent{
		Name:        "CompetitorIntelligenceAgent",
		Description: "Analyzes competitor data and generates competitive intelligence reports",
	}
}

// MarketResearch searches for competitor data
func (a *CompetitorIntelligenceAgent) MarketResearch(ctx context.Context, companyName string, industry string) ([]CompetitorData, error) {
	// Simulated market research - in production, this would call external APIs
	// like Crunchbase, LinkedIn, or industry-specific data sources
	competitors := []CompetitorData{
		{
			Name:        "Competitor A",
			Website:     "https://competitor-a.com",
			Industry:    industry,
			Products:    []string{"Product 1", "Product 2", "Product 3"},
			Pricing:     "Premium",
			MarketShare: 25.5,
			Strengths:   []string{"Strong brand", "Large customer base", "Innovation"},
			Weaknesses:  []string{"High prices", "Slow support", "Limited features"},
		},
		{
			Name:        "Competitor B",
			Website:     "https://competitor-b.com",
			Industry:    industry,
			Products:    []string{"Product X", "Product Y"},
			Pricing:     "Mid-range",
			MarketShare: 18.2,
			Strengths:   []string{"Affordable", "Good UX", "Fast growth"},
			Weaknesses:  []string{"Limited market presence", "Newer player", "Fewer integrations"},
		},
		{
			Name:        "Competitor C",
			Website:     "https://competitor-c.com",
			Industry:    industry,
			Products:    []string{"Enterprise Suite"},
			Pricing:     "Enterprise",
			MarketShare: 12.8,
			Strengths:   []string{"Enterprise features", "Security", "Compliance"},
			Weaknesses:  []string{"Expensive", "Complex setup", "Steep learning curve"},
		},
	}

	return competitors, nil
}

// Analyze performs competitive positioning analysis
func (a *CompetitorIntelligenceAgent) Analyze(ctx context.Context, data []CompetitorData) ([]CompetitorAnalysis, error) {
	var analyses []CompetitorAnalysis

	for _, competitor := range data {
		analysis := CompetitorAnalysis{
			CompetitorName: competitor.Name,
		}

		// Determine threat level based on market share
		switch {
		case competitor.MarketShare > 20:
			analysis.ThreatLevel = "High"
		case competitor.MarketShare > 10:
			analysis.ThreatLevel = "Medium"
		default:
			analysis.ThreatLevel = "Low"
		}

		// Determine positioning based on pricing
		switch competitor.Pricing {
		case "Premium":
			analysis.Positioning = "Premium market leader"
		case "Mid-range":
			analysis.Positioning = "Value-focused challenger"
		case "Enterprise":
			analysis.Positioning = "Enterprise specialist"
		default:
			analysis.Positioning = "Undifferentiated"
		}

		// Extract key differentiators from strengths
		analysis.KeyDifferentiators = competitor.Strengths

		// Generate opportunities based on competitor weaknesses
		for _, weakness := range competitor.Weaknesses {
			opportunity := fmt.Sprintf("Capitalize on %s weakness", weakness)
			analysis.Opportunities = append(analysis.Opportunities, opportunity)
		}

		// Generate risks based on competitor strengths
		for _, strength := range competitor.Strengths {
			risk := fmt.Sprintf("Competitor's %s advantage", strength)
			analysis.Risks = append(analysis.Risks, risk)
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GenerateReport creates a comprehensive competitive intelligence report
func (a *CompetitorIntelligenceAgent) GenerateReport(ctx context.Context, targetCompany string, analyses []CompetitorAnalysis) (*CompetitorReport, error) {
	report := &CompetitorReport{
		GeneratedAt:   time.Now(),
		TargetCompany: targetCompany,
		Competitors:   analyses,
	}

	// Generate market insights
	totalMarketShare := 0.0
	for _, analysis := range analyses {
		for _, data := range analyses {
			if data.CompetitorName == analysis.CompetitorName {
				// This is a simplified calculation
				totalMarketShare += 10.0
			}
		}
	}

	report.MarketInsights = fmt.Sprintf(
		"The competitive landscape shows %d major players. "+
			"High-threat competitors control significant market share. "+
			"Opportunities exist in underserved segments.",
		len(analyses),
	)

	// Generate strategic recommendations
	report.Recommendations = []string{
		"Focus on differentiation in areas where competitors are weak",
		"Target mid-market segment with competitive pricing",
		"Invest in customer support to outperform competitors",
		"Develop integrations to match competitor ecosystems",
		"Monitor competitor pricing and adjust strategy quarterly",
	}

	return report, nil
}

// Run executes the full competitor intelligence workflow
func (a *CompetitorIntelligenceAgent) Run(ctx context.Context, companyName string, industry string) (*CompetitorReport, error) {
	// Step 1: Market Research
	data, err := a.MarketResearch(ctx, companyName, industry)
	if err != nil {
		return nil, fmt.Errorf("market research failed: %w", err)
	}

	// Step 2: Analysis
	analyses, err := a.Analyze(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	// Step 3: Generate Report
	report, err := a.GenerateReport(ctx, companyName, analyses)
	if err != nil {
		return nil, fmt.Errorf("report generation failed: %w", err)
	}

	return report, nil
}

// ToJSON converts the report to JSON format
func (r *CompetitorReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
