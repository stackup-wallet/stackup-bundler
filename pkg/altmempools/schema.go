package altmempools

import (
	"embed"
	"io/fs"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed *AltMempoolSchema.json
var files embed.FS

func newJSONSchema() (*jsonschema.Schema, error) {
	var s string
	err := fs.WalkDir(files, "AltMempoolSchema.json", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(files, path)
		if err != nil {
			return err
		}

		s = string(b)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sch, err := jsonschema.CompileString("AltMempoolSchema.json", s)
	if err != nil {
		return nil, err
	}
	return sch, nil
}
