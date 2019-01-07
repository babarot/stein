package lint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl"

	"k8s.io/apimachinery/pkg/util/yaml"
)

// File represents the files to be linted
// It's converted from the arguments
type File struct {
	Path string
	Data []byte

	// Meta field means the annotation
	Meta string
}

// Args converts from its arguments to the collection of File object
func Args(paths []string) (files []File, err error) {
	for _, path := range paths {
		ext := filepath.Ext(path)
		switch ext {
		case ".yaml", ".yml":
			yamlFiles, err := handleYAML(path)
			if err != nil {
				return files, err
			}
			files = append(files, yamlFiles...)
		case ".json":
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return files, err
			}
			files = append(files, File{
				Path: path,
				Data: data,
			})
		case ".hcl", ".tf":
			contents, err := ioutil.ReadFile(path)
			if err != nil {
				return files, err
			}
			var v interface{}
			err = hcl.Unmarshal(contents, &v)
			if err != nil {
				return files, fmt.Errorf("unable to parse HCL: %s", err)
			}
			data, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return files, fmt.Errorf("unable to marshal json: %s", err)
			}
			files = append(files, File{
				Path: path,
				Data: data,
			})
		default:
			return files, fmt.Errorf("%q (%s): unsupported file type", path, ext)
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
			// If more than one block is defined in one YAML file,
			// notes the numbering of the block as Meta field
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
