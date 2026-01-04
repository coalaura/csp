package main

import (
	"encoding/json"
	"io"
)

type Report struct {
	ViolatedDirective string
	DocumentURL       string
	BlockedURL        string
}

type ReportUriBody struct {
	DocumentURI        string `json:"document-uri"`
	BlockedURI         string `json:"blocked-uri"`
	ViolatedDirective  string `json:"violated-directive"`
	EffectiveDirective string `json:"effective-directive"`
	OriginalPolicy     string `json:"original-policy"`
	Disposition        string `json:"disposition"`
	StatusCode         int    `json:"status-code"`
}

type ReportUriPayload struct {
	Report ReportUriBody `json:"csp-report"`
}

type ReportToBody struct {
	DocumentURL       string `json:"documentURL"`
	BlockedURL        string `json:"blockedURL"`
	ViolatedDirective string `json:"violatedDirective"`
}

type ReportToPayload struct {
	Type string       `json:"type"`
	Age  int          `json:"age"`
	URL  string       `json:"url"`
	Body ReportToBody `json:"body"`
}

func ParseReport(rd io.Reader) []Report {
	body, err := io.ReadAll(io.LimitReader(rd, 64*1024))
	if err != nil {
		return nil
	}

	var legacy ReportUriPayload

	err = json.Unmarshal(body, &legacy)
	if err == nil && legacy.Report.DocumentURI != "" && legacy.Report.ViolatedDirective != "" {
		report := legacy.Report

		return []Report{
			{
				ViolatedDirective: report.ViolatedDirective,
				DocumentURL:       report.DocumentURI,
				BlockedURL:        report.BlockedURI,
			},
		}
	}

	var modern []ReportToPayload

	err = json.Unmarshal(body, &modern)
	if err == nil {
		reports := make([]Report, 0, len(modern))

		for _, report := range modern {
			body := report.Body

			if body.DocumentURL == "" || body.ViolatedDirective == "" {
				continue
			}

			reports = append(reports, Report{
				ViolatedDirective: body.ViolatedDirective,
				DocumentURL:       body.DocumentURL,
				BlockedURL:        body.BlockedURL,
			})
		}

		return reports
	}

	return nil
}
