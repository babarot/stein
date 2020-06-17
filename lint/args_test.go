package lint

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_filesFromArgs(t *testing.T) {
	tests := []struct {
		Filename string
		Want     string
	}{
		{
			"testdata/01.tf",
			`{
  "provider": [
    {
      "google": [
        {
          "project": "my-project-id",
          "region": "us-central1"
        }
      ]
    }
  ]
}`,
		},
		{
			"testdata/02.tf",
			`{
  "provider": [
    {
      "google": [
        {
          "project": "my-project-id",
          "region": "us-central1"
        }
      ]
    },
    {
      "aws": [
        {
          "region": "us-east-1",
          "version": "~\u003e 2.0"
        }
      ]
    },
    {
      "google": [
        {
          "alias": "west",
          "project": "my-project-id",
          "region": "us-west1"
        }
      ]
    }
  ]
}`,
		},
		{
			"testdata/03.tf",
			`{
  "module": [
    {
      "foo": [
        {
          "array": [
            "val1",
            "val2"
          ],
          "bar": "baz"
        }
      ]
    }
  ]
}`,
		},
		{
			"testdata/04.tf",
			`{
  "resource": [
    {
      "google_project": [
        {
          "my_project": [
            {
              "name": "My Project",
              "org_id": "1234567",
              "project_id": "your-project-id"
            }
          ]
        }
      ]
    }
  ]
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Filename, func(t *testing.T) {
			files, err := filesFromArgs([]string{test.Filename})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if len(files) != 1 {
				t.Errorf("unexpected files length: got %d, want 1", len(files))
			}
			got := string(files[0].Data)
			if got != test.Want {
				t.Errorf("wrong result: got: %#v, want: %#v", got, test.Want)
			}
		})
	}
}

// This test is based on https://github.com/kubernetes/apimachinery/blob/d8530e6c952f75365336be8ea29cfd758ce49ee8/pkg/util/yaml/decoder_test.go#L57-L82
func Test_handleYAML_LineTooLong(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "Test_handleYAML_LineTooLong")
	if err != nil {
		t.Fatalf("failed to create temporary file: %s", err)
	}
	defer os.Remove(tmpfile.Name())

	d := `
stuff: 1
`
	//  maxLen 5 M
	dd := strings.Repeat(d, 512*1024)
	if _, err := tmpfile.WriteString(dd); err != nil {
		t.Fatalf("failed to write to temporary file: %s", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temporary file: %s", err)
	}

	_, err = handleYAML(tmpfile.Name())
	if err != bufio.ErrTooLong {
		t.Fatalf("want %q, got %q", bufio.ErrTooLong, err)
	}
}
