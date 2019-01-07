package lint

import (
	"bytes"
	"errors"
	"fmt"
	tt "html/template"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/b4b4r07/stein/lang"
	"github.com/b4b4r07/stein/lang/loader"
	"github.com/fatih/color"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/k0kubun/pp"
)

const (
	// RulePrefix is a prefix of rule name
	RulePrefix = "rule."

	// DefaultFormat is a default format string
	DefaultFormat = "[{{.Level}}]  {{.Rule}}  {{.Message}}"
)

// Status is
type Status int

const (
	// Success is
	Success Status = iota
	// Failure is
	Failure
)

// Linter is a linter structure
type Linter struct {
	stdout io.Writer
	stderr io.Writer

	config *Config
	cache  cache

	// policy schema
	policy *lang.Policy
	// all configuration body decoded by HCL
	body hcl.Body

	// files map[string]*hcl.File
}

type cache struct {
	// policies represents all policy files loaded by Run method
	policies map[string]Policy

	// policy represent the structure of Policy corresponding to currently loaded YAML
	policy Policy

	// filepath represents the file path of current YAML loaded in Run method
	filepath string
}

// NewLinter is
func NewLinter(policy loader.Policy) *Linter {
	return &Linter{
		stdout: os.Stdout,
		stderr: os.Stderr,
		cache: cache{
			policies: map[string]Policy{},
			policy:   Policy{},
			filepath: "",
		},
		body:   policy.Body,
		policy: policy.Data,
		// files:  loaded.Files,
	}
}

func (l *Linter) decodePolicy(yamlfile string) (Policy, error) {
	var policy Policy

	ctx, diags := l.policy.BuildContext(l.body, yamlfile)
	if diags.HasErrors() {
		return policy, diags
	}

	decodeDiags := gohcl.DecodeBody(l.body, ctx, &policy)
	diags = append(diags, decodeDiags...)
	if diags.HasErrors() {
		return policy, diags
	}

	// policy.Config can be nil. If so, it should be set to default value
	if policy.Config == nil {

		// Default config setting
		policy.Config = &Config{
			Report: ReportConfig{
				Format: DefaultFormat,
				Style:  "console",
				Color:  true,
			},
		}
	}

	return policy, nil
}

// Result represents
type Result struct {
	Filename string
	Items    []Item
	OK       bool
}

// Item is
type Item struct {
	Name    string
	Message string
	Status  Status
	Level   string
}

// Run runs the linter
func (l *Linter) Run(yamlfile string) (Result, error) {
	policy, err := l.decodePolicy(yamlfile)
	if err != nil {
		return Result{}, err
	}

	if err := policy.Validate(); err != nil {
		return Result{}, err
	}

	l.cache.policies[yamlfile] = policy
	l.cache.policy = policy
	l.cache.filepath = yamlfile
	l.config = policy.Config

	if err := l.Validate(); err != nil {
		return Result{}, err
	}

	result := Result{
		Filename: yamlfile,
		Items:    []Item{},
		OK:       true,
	}

	length := l.calcReportLength()
	for _, rule := range policy.Rules {
		if err := rule.Validate(); err != nil {
			return result, err
		}

		message, err := rule.BuildMessage(policy.Config.Report, length)
		if err != nil {
			return result, err
		}

		result.Items = append(result.Items, Item{
			Name:    rule.Name,
			Message: message,
			Status:  rule.getStatus(),
			Level:   rule.Report.Level,
		})

		// this linter will fail if it has even one failed rule.
		if rule.getStatus() != Success {
			result.OK = false
		}
		for _, debug := range rule.Debugs {
			pp.Println(debug)
		}
	}

	return result, nil
}

func (r *Rule) getStatus() Status {
	if r.SkipCase() {
		return Success
	}
	if r.TrueCase() {
		return Success
	}
	return Failure
}

// Print is
func (l *Linter) Print(result Result) {
	var (
		out            = l.stderr
		style          = l.config.Report.Style
		consolePadding = "  "
	)

	// setup print
	switch style {
	case "console":
		color.New(color.Underline).Fprintln(out, result.Filename)
	}

	// main print
	for _, rule := range result.Items {
		// Do not print successful items
		if rule.Status == Success {
			continue
		}
		switch style {
		case "console":
			fmt.Fprintf(out, consolePadding)
		}
		fmt.Fprintln(out, rule.Message)
	}

	// teardown print
	switch style {
	case "console":
		if result.OK {
			fmt.Fprintln(out, consolePadding+"No violated rules")
		}
		fmt.Fprintln(out)
	}
}

// Status is
func (l *Linter) Status(results ...Result) int {
	for _, result := range results {
		if !result.OK {
			return 1
		}
	}
	return 0
}

// PrintSummary is [TODO]
func (l *Linter) PrintSummary(results ...Result) {
	s := struct {
		warns  int
		errors int
	}{}
	for _, result := range results {
		for _, item := range result.Items {
			// switch item.Status {
			// case Success:
			// case Failure:
			// 	s.errors++
			// }
			if item.Status == Success {
				continue
			}
			switch item.Level {
			case "ERROR":
				s.errors++
			case "WARN":
				s.warns++
			}
		}
	}
	result := fmt.Sprintf("%d error(s), %d warn(s)", s.errors, s.warns)
	fmt.Fprintln(l.stderr, strings.Repeat("=", len(result)))
	fmt.Fprintln(l.stderr, result)
}

// SkipCase returns if there is even one IgnoreCases.
func (r *Rule) SkipCase() bool {
	for _, ignore := range r.IgnoreCases {
		if ignore {
			return true
		}
	}
	return false
}

// TrueCase returns true if all expressions in a rule are true.
func (r *Rule) TrueCase() bool {
	for _, expr := range r.Expressions {
		if !expr {
			return false
		}
	}
	return true
}

// Validate validates linter configuration
func (l *Linter) Validate() error {
	validations := []struct {
		rule    bool
		message string
	}{}

	for _, validation := range validations {
		if !validation.rule {
			return fmt.Errorf(validation.message)
		}
	}

	return nil
}

// Validate validates rule syntax
func (r *Rule) Validate() error {
	validations := []struct {
		rule    bool
		message string
	}{
		{
			r.Report.Level == "ERROR" || r.Report.Level == "WARN",
			"report level accepts only ERROR or WARN",
		},
		{
			len(r.Report.Message) > 0,
			fmt.Sprintf("%s: report message should be written", r.Name),
		},
	}

	for _, validation := range validations {
		if !validation.rule {
			return fmt.Errorf(validation.message)
		}
	}

	return nil
}

// ReportLength is the information of the max length of each format strings
type ReportLength struct {
	// Max length of RulePrefix + {{.Rule}}
	MaxRuleName int

	// Max length of {{.Level}}
	MaxLevel int

	// Max length of {{.Message}}
	MaxMessage int
}

// calcReportLength is a method that measures how long length each placeholder
// used in a template actually occurred the maximum length.
//
// Example:
//   [ERROR] rule.one_resource_per_one_file  Only 1 resource should be defined in a YAML file
//   [WARN ] rule.yaml_separator             YAML separator "---" should be removed
//
// In this case, calcReportLength below will be returned
//   max level length: 5, max rule name length: 30, max message length: 48
//
func (l *Linter) calcReportLength() ReportLength {
	var length ReportLength

	for _, rule := range l.cache.policy.Rules {
		if len(rule.Name) > length.MaxRuleName {
			length.MaxRuleName = len(rule.Name)
		}
		if len(rule.Report.Level) > length.MaxLevel {
			length.MaxLevel = len(rule.Report.Level)
		}
		if len(rule.Report.Message) > length.MaxMessage {
			length.MaxMessage = len(rule.Report.Message)
		}
	}
	return length
}

// BuildMessage formats the results reported by linter.
func (r *Rule) BuildMessage(cfg ReportConfig, length ReportLength) (string, error) {
	format := DefaultFormat
	if len(cfg.Format) > 0 {
		format = cfg.Format
	}

	renderedFormat := new(bytes.Buffer)
	tpl, err := tt.New("").Parse(format)
	if err != nil {
		return "", err
	}

	var (
		ruleName = r.Name
		level    = r.Report.Level
		message  = r.Report.Message
	)

	var (
		rulePadding    = strings.Repeat(" ", length.MaxRuleName-len(ruleName))
		levelPadding   = strings.Repeat(" ", length.MaxLevel-len(level))
		messagePadding = strings.Repeat(" ", length.MaxMessage-len(message))
	)

	if cfg.Color {
		switch level {
		case "ERROR":
			level = color.RedString(level)
		case "WARN":
			level = color.YellowString(level)
		}
		// Colorize by default in case of only no advance color specification
		if !containsANSI(message) && !containsANSI(format) {
			message = color.WhiteString(message)
		}
	}

	err = tpl.Execute(renderedFormat, map[string]interface{}{
		"Rule":    RulePrefix + ruleName + rulePadding,
		"Level":   level + levelPadding,
		"Message": tt.HTML(message + messagePadding),
	})
	if err != nil {
		return "", err
	}

	switch renderedFormat.Len() {
	case 0:
		// something wrong TODO
		return "", errors.New("error happen")
	default:
		format = renderedFormat.String()
	}

	return format, nil
}

func stripANSI(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}

func containsANSI(str string) bool {
	return str != stripANSI(str)
}
