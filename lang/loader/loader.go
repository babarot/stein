package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/b4b4r07/stein/lang"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// // Loaded is
// type Loaded struct {
// 	Policy *lang.Policy
// 	Files  map[string]*hcl.File
// 	Body   hcl.Body
// }

// Parser is
type Parser struct {
	p *hclparse.Parser
}

// NewParser is
// // Loader is a starting point function for loading the HCL file
// // (called this policy in this application).
func NewParser() *Parser {
	return &Parser{hclparse.NewParser()}
}

func (p *Parser) loadHCLFile(path string) (hcl.Body, hcl.Diagnostics) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The file %q could not be read.", path),
			},
		}
	}

	var file *hcl.File
	var diags hcl.Diagnostics
	switch {
	case strings.HasSuffix(path, ".json"):
		file, diags = p.p.ParseJSON(src, path)
	default:
		file, diags = p.p.ParseHCL(src, path)
	}

	// If the returned file or body is nil, then we'll return a non-nil empty
	// body so we'll meet our contract that nil means an error reading the file.
	if file == nil || file.Body == nil {
		return hcl.EmptyBody(), diags
	}

	return file.Body, diags
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		switch filepath.Ext(path) {
		case ".hcl":
			*files = append(*files, path)
		}
		return nil
	}
}

func getPolicyFiles(path string) ([]string, error) {
	var (
		files []string
		err   error
	)
	fi, err := os.Stat(path)
	if err != nil {
		return files, err
	}
	if fi.IsDir() {
		return files, filepath.Walk(path, visit(&files))
	}
	switch filepath.Ext(path) {
	case ".hcl":
		files = append(files, path)
	}
	return files, err
}

// func readBody(path string) (hcl.Body, map[string]*hcl.File, error) {
// 	parser := NewParser()
//
// 	var diags hcl.Diagnostics
// 	var bodies []hcl.Body
// 	var err error
//
// 	files, err := getPolicyFiles(path)
// 	if err != nil {
// 		return nil, map[string]*hcl.File{}, err
// 	}
//
// 	for _, file := range files {
// 		body, fDiags := parser.loadHCLFile(file)
// 		bodies = append(bodies, body)
// 		diags = append(diags, fDiags...)
// 	}
//
// 	if diags.HasErrors() {
// 		err = diags
// 	}
//
// 	return hcl.MergeBodies(bodies), parser.p.Files(), err
// }

// Policy is
type Policy struct {
	Body  hcl.Body
	Files map[string]*hcl.File
	Data  *lang.Policy
}

// Load is
func Load(paths ...string) (Policy, error) {
	parser := NewParser()

	var diags hcl.Diagnostics
	var bodies []hcl.Body
	var err error

	for _, path := range paths {
		files, err := getPolicyFiles(path)
		if err != nil {
			return Policy{}, err
		}

		for _, file := range files {
			body, fDiags := parser.loadHCLFile(file)
			bodies = append(bodies, body)
			diags = append(diags, fDiags...)
		}
	}

	if diags.HasErrors() {
		err = diags
	}

	return Policy{
		Body:  hcl.MergeBodies(bodies),
		Files: parser.p.Files(),
	}, err
}

// // Load is
// func (l *Loader) Load() (*lang.Policy, error) {
// 	// body, files, err := readBody(path)
// 	// if err != nil {
// 	// 	return Loaded{
// 	// 		Policy: nil,
// 	// 		Body:   body,
// 	// 		Files:  files,
// 	// 	}, err
// 	// }
// 	var err error
// 	policy, diags := lang.Parse(l.Body, l.Files)
// 	if diags.HasErrors() {
// 		err = diags
// 	}
// 	return policy, err
// 	// return Loaded{
// 	// 	Policy: policy,
// 	// 	Body:   l.Body,
// 	// 	Files:  l.Files,
// 	// }, err
// }
