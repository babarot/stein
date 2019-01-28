package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/b4b4r07/stein/lint/internal/policy"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// Parser represents mainly HCL parser
type Parser struct {
	p *hclparse.Parser
}

// NewParser creates Parser instance
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

func visitHCL(files *[]string) filepath.WalkFunc {
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

// getPolicyFiles walks the given path and returns the files ending with HCL
// Also, it returns the path if the path is just a file and a HCL file
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
		return files, filepath.Walk(path, visitHCL(&files))
	}
	switch filepath.Ext(path) {
	case ".hcl":
		files = append(files, path)
	}
	return files, err
}

// SearchPolicyDir searchs the hierarchy of the given path step by step and find the default directory
func SearchPolicyDir(path string) []string {
	var dirs []string
	for {
		if !strings.Contains(path, "/") {
			break
		}
		// search parent dir nextly
		path = filepath.Dir(path)
		policyDirPath := filepath.Join(path, ".policy")
		_, err := os.Stat(policyDirPath)
		if err != nil {
			// skip if policy dir doesn't exist
			continue
		}
		dirs = append(dirs, policyDirPath)
	}
	return dirs
}

// Policy represents the raw data of HCL files decoded by hcl2 package
type Policy struct {
	Body  hcl.Body
	Files map[string]*hcl.File
	// Data represents the raw data decoded based on stein schema
	Data *policy.Policy
}

// Load reads the files and converts them to Policy object
func Load(paths ...string) (Policy, error) {
	parser := NewParser()

	var diags hcl.Diagnostics
	var bodies []hcl.Body
	var policies []string
	var err error

	// paths can take a file path and a dir path
	for _, path := range paths {
		// if c.Option.Policy is empty, in other words, additionals is nothing,
		// paths are likely to contain empty string.
		// if so, skip to run getPolicyFiles
		if path == "" {
			continue
		}
		files, err := getPolicyFiles(path)
		if err != nil {
			return Policy{}, err
		}
		// gather full paths of HCL file to one array
		policies = append(policies, files...)
	}

	// delete duplicate file paths
	// in consideration of the case the same files are read
	//
	// TODO: think if unique is needed (if not needed, just returns error)
	for _, policy := range unique(policies) {
		body, fDiags := parser.loadHCLFile(policy)
		bodies = append(bodies, body)
		diags = append(diags, fDiags...)
	}

	if diags.HasErrors() {
		err = diags
	}

	return Policy{
		Body:  hcl.MergeBodies(bodies),
		Files: parser.p.Files(),
	}, err
}

func unique(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}
