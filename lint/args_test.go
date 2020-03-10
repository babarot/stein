package lint

import (
	"bufio"
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

func Test_handleYAML_LineTooLong(t *testing.T) {
	_, err := handleYAML("testdata/line-too-long.yaml")
	if err != bufio.ErrTooLong {
		t.Fatalf("want %q, got %q", bufio.ErrTooLong, err)
	}
}
