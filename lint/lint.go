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

	"github.com/b4b4r07/stein/lint/internal/policy"
	"github.com/b4b4r07/stein/lint/internal/policy/loader"
	"github.com/b4b4r07/stein/lint/internal/topological"
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

	// LevelError represents the error level reported by lint
	LevelError = "ERROR"
	// LevelWarning represents the warning level reported by lint
	LevelWarning = "WARN"
)

// Status represents the status code of Lint
type Status int

const (
	// Success is the success code of Lint
	Success Status = iota
	// Failure is the failure code of Lint
	Failure
)

// Linter is a linter structure
type Linter struct {
	stdout io.Writer
	stderr io.Writer

	config *Config
	cache  cache

	// policy is a policy schema
	policy *policy.Policy
	// body is a decoded body that all policies is merged
	body hcl.Body
	// files are converted from given arguments (file paths are assumed)
	files []File
}

type cache struct {
	// policies represents all policy files loaded by Run method
	policies map[string]Policy

	// policy represent the structure of Policy corresponding to currently loaded YAML
	policy Policy

	// filepath represents the file path of current YAML loaded in Run method
	filepath string
}

// NewLinter creates Linter object based on Lint Policy
func NewLinter(args []string, additionals ...string) (*Linter, error) {
	files, err := filesFromArgs(args, additionals...)
	if err != nil {
		return &Linter{}, err
	}
	return &Linter{
		stdout: os.Stdout,
		stderr: os.Stderr,
		cache: cache{
			policies: map[string]Policy{},
			policy:   Policy{},
			filepath: "",
		},
		policy: nil,
		body:   nil,
		files:  files,
	}, nil
}

func (l *Linter) setPolicy(policy loader.Policy) {
	l.body = policy.Body
	l.policy = policy.Data
}

func (l *Linter) decodePolicy(file File) (Policy, error) {
	var policy Policy

	if l.policy == nil {
		return policy, errors.New("no DAM")
	}

	ctx, diags := l.policy.BuildContext(l.body, file.Path, file.Data)
	if diags.HasErrors() {
		return policy, diags
	}

	decodeDiags := gohcl.DecodeBody(l.body, ctx, &policy)
	diags = append(diags, decodeDiags...)
	if diags.HasErrors() {
		return policy, diags
	}

	// policy.Config can be nil
	// In that case it should be set to default value
	if policy.Config == nil {
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

// Result represents the execution result of Lint
// It's represented against one argument
// The result of each rules for the argument is stored in Items
type Result struct {
	Path  string
	Items Items
	OK    bool

	// RulesNotFound is
	RulesNotFound bool
	// Metadata is something notes related to Result
	Metadata string
}

// Item represents the result of a rule
type Item struct {
	Name    string
	Message string
	Status  Status
	Level   string
}

// Items is the collenction of Item object
type Items []Item

// Files returns []File object converted from given arguments
func (l *Linter) Files() []File {
	return l.files
}

// Run runs the linter against a file of an argument
func (l *Linter) Run(file File) (Result, error) {
	// Set a policy to read for each files
	l.setPolicy(file.Policy)

	policy, err := l.decodePolicy(file)
	if err != nil {
		return Result{}, err
	}

	if err := policy.Validate(); err != nil {
		return Result{}, err
	}

	l.cache.policies[file.Path] = policy
	l.cache.policy = policy
	l.cache.filepath = file.Path
	l.config = policy.Config

	if err := l.Validate(); err != nil {
		return Result{}, err
	}

	result := Result{
		Path:          file.Path,
		Items:         []Item{},
		OK:            true,
		RulesNotFound: len(policy.Rules) == 0,
		Metadata:      file.Meta,
	}

	// sort rules by depends_on
	policy.Rules.Sort()

	length := l.calcReportLength()
	for _, rule := range policy.Rules {
		if err := rule.Validate(); err != nil {
			return result, err
		}

		if rule.hasDependencies() {
			// Skip execution of the rule if the dependent rule fails
			ok := rule.checkDependenciesAreOK(result)
			if !ok {
				continue
			}
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

// Sort sorts the rules based on its own dependencies
//
// For example, in the case that these rules are defined like below,
// the order which the rules are executed should be as follows:
//
//   rule.a --- rule.b --- rule.c
//                      `- rule.d
//
//   rule "a" {
//     ...
//   }
//   rule "b" {
//     depends_on = ["rule.a"]
//   }
//   rule "c" {
//     depends_on = ["rule.b"]
//   }
//   rule "d" {
//     depends_on = ["rule.a"]
//   }
//
// This implementation is based on the algorithm of topological sort
//
func (r *Rules) Sort() {
	graph := topological.NewGraph(len(*r))
	for _, rule := range *r {
		graph.AddNode(rule.Name)
	}

	for _, rule := range *r {
		if !rule.hasDependencies() {
			continue
		}
		for _, dependency := range rule.Dependencies {
			dependency = strings.TrimPrefix(dependency, RulePrefix)
			graph.AddEdge(dependency, rule.Name)
		}
	}

	orderedRuleNames, ok := graph.Sort()
	if !ok {
		// TODO: Handle this pattern
		//   For now it can be ignored.
		//
		// return
	}

	var sortedRules Rules
	for _, orderedRuleName := range orderedRuleNames {
		sortedRules = append(sortedRules, r.getOneByName(orderedRuleName))
	}
	*r = sortedRules
}

func (r Rules) getOneByName(name string) Rule {
	for _, rule := range r {
		if rule.Name == name {
			return rule
		}
	}
	return Rule{}
}

func (r *Rule) hasDependencies() bool {
	return len(r.Dependencies) > 0
}

func (r *Rule) checkDependenciesAreOK(result Result) bool {
	for _, dependency := range r.Dependencies {
		depRule := strings.TrimPrefix(dependency, RulePrefix)
		item := result.Items.getOneByName(depRule)
		switch item.Status {
		case Success:
			// even if this item succeeds,
			// checks all other items
			continue
		case Failure:
			return false
		}
	}
	return true
}

func (i Items) getOneByName(name string) Item {
	for _, item := range i {
		if item.Name == name {
			return item
		}
	}
	return Item{}
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

// Print prints the result of Lint based on Result reported by Run
func (l *Linter) Print(result Result) {
	const consolePadding = "  "

	var (
		out = l.stderr
		cfg = l.config.Report
	)

	// Setup Print method
	switch cfg.Style {
	case "console":
		color.New(color.Underline).Fprintf(out, result.Path)
		if len(result.Metadata) > 0 {
			metadata := fmt.Sprintf(" (%s)", result.Metadata)
			if cfg.Color {
				metadata = color.CyanString(metadata)
			}
			fmt.Fprintf(out, metadata)
		}
		fmt.Fprintln(out)
	}

	// Main logic of Print
	for _, rule := range result.Items {
		// Do not print successful items
		if rule.Status == Success {
			continue
		}
		switch cfg.Style {
		case "console":
			fmt.Fprintf(out, consolePadding)
		}
		fmt.Fprintln(out, rule.Message)
	}

	// Teardown Print method
	switch cfg.Style {
	case "console":
		switch {
		case result.RulesNotFound:
			fmt.Fprintln(out, consolePadding+"No rules found")
		case result.OK:
			fmt.Fprintln(out, consolePadding+"No violated rules")
		}
		fmt.Fprintln(out)
	}
}

// Status indicates execution result of Lint by the status code
func (l *Linter) Status(results ...Result) Status {
	for _, result := range results {
		if !result.OK {
			return Failure
		}
	}
	return Success
}

// PrintSummary prints the summary of all results of the entire Lint
func (l *Linter) PrintSummary(results ...Result) {
	s := struct {
		warns  int
		errors int
	}{}
	for _, result := range results {
		for _, item := range result.Items {
			if item.Status == Success {
				continue
			}
			switch item.Level {
			case LevelError:
				s.errors++
			case LevelWarning:
				s.warns++
			}
		}
	}
	result := fmt.Sprintf("%d error(s), %d warn(s)", s.errors, s.warns)
	fmt.Fprintln(l.stderr, strings.Repeat("=", len(result)))
	fmt.Fprintln(l.stderr, result)
}

// SkipCase returns true if cases of a precondition block includes one or more failed cases
func (r *Rule) SkipCase() bool {
	// don't skip if a precondition block is not specified
	if r.Precondition == nil {
		return false
	}
	for _, ok := range r.Precondition.Cases {
		if !ok {
			return true
		}
	}
	return false
}

// TrueCase returns true if all conditions in a rule are true.
func (r *Rule) TrueCase() bool {
	for _, expr := range r.Conditions {
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
			r.Report.Level == LevelError || r.Report.Level == LevelWarning,
			fmt.Sprintf("report level accepts only %s or %s", LevelError, LevelWarning),
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
		case LevelError:
			level = color.RedString(level)
		case LevelWarning:
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
