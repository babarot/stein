package lint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"

	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/b4b4r07/stein/lint/hclconvert"
	"github.com/b4b4r07/stein/lint/internal/policy"
	"github.com/b4b4r07/stein/lint/internal/policy/loader"
)

// File represents the files to be linted
// It's converted from the arguments
type File struct {
	Path string
	Data []byte

	// Meta field means the annotation
	Meta string

	// File struct has Policy data
	// because policy applied to the file should be determined by each file
	Policy loader.Policy

	Diagnostics hcl.Diagnostics
}

// filesFromArgs converts from given arguments to the collection of File object
func filesFromArgs(args []string, additionals ...string) (files []File, err error) {
	log.Printf("[TRACE] converting from args to lint.Files\n")

	for _, arg := range args {
		var diags hcl.Diagnostics
		log.Printf("[INFO] converting lint.File: %s\n", arg)
		policies := loader.SearchPolicyDir(arg)
		policies = append(policies, additionals...)
		log.Printf("[INFO] policies: %#v\n", policies)

		loadedPolicy, err := loader.Load(policies...)
		if err != nil {
			switch e := err.(type) {
			case hcl.Diagnostics:
				log.Printf("[DEBUG] diags reported from loader.Load in fileFromArg\n")
				diags = append(diags, e...)
			case error:
				return files, err
			}
		}

		data, decodeDiags := policy.Decode(loadedPolicy.Body)
		diags = append(diags, decodeDiags...)

		loadedPolicy.Data = data

		ext := filepath.Ext(arg)
		switch ext {
		case ".yaml", ".yml":
			yamlFiles, err := handleYAML(arg)
			if err != nil {
				return files, err
			}
			log.Printf("[TRACE] %d block(s) found in YAML: %s\n", len(yamlFiles), arg)
			for _, file := range yamlFiles {
				file.Policy = loadedPolicy
				file.Diagnostics = diags
				files = append(files, file)
			}
		case ".json":
			data, err := ioutil.ReadFile(arg)
			if err != nil {
				return files, err
			}
			files = append(files, File{
				Path:        arg,
				Data:        data,
				Policy:      loadedPolicy,
				Diagnostics: diags,
			})
		case ".hcl", ".tf":
			contents, err := ioutil.ReadFile(arg)
			if err != nil {
				return files, err
			}
			f, d := hclsyntax.ParseConfig(contents, arg, hcl.Pos{Line: 1, Column: 1})
			if d.HasErrors() {
				return files, fmt.Errorf("unable to parse HCL: %s", d.Error())
			}

			v, err := hclconvert.ConvertFile(f)
			if err != nil {
				return files, fmt.Errorf("unable to convert HCL: %s", d.Error())
			}

			data, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return files, fmt.Errorf("unable to marshal json: %s", err)
			}

			files = append(files, File{
				Path:        arg,
				Data:        data,
				Policy:      loadedPolicy,
				Diagnostics: diags,
			})
		default:
			return files, fmt.Errorf("%q (%s): unsupported file type", arg, ext)
		}
	}

	return files, nil
}

func handleYAML(path string) (files []File, err error) {
	file, err := os.Open(path)
	if err != nil {
		return files, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return files, err
	}

	dd := yaml.NewDocumentDecoder(file)
	defer dd.Close()
	var documents [][]byte
	for {
		res := make([]byte, fi.Size())
		_, err := dd.Read(res)
		if err == io.EOF {
			break
		}
		documents = append(documents, bytes.Trim(res, "\x00"))
	}

	for idx, document := range documents {
		data, err := yaml.ToJSON(document)
		if err != nil {
			return files, err
		}
		meta := ""
		if len(documents) > 1 {
			// If one or more blocks are defined in one YAML file,
			// records the numbering of the block in Meta field
			meta = fmt.Sprintf("Block %d", idx+1)
		}
		files = append(files, File{
			Path: path,
			Data: data,
			Meta: meta,
		})
	}

	return files, err
}
