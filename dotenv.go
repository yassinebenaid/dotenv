package dotenv

import (
	"fmt"
	"os"
)

func Read(files ...string) (map[string]string, error) {
	if len(files) == 0 {
		files = append(files, "./.env")
	}

	p := parser{env: make(map[string]string)}
	for _, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading file [%s], %v", path, err)
		}
		p.input = content
		if err := p.parse(); err != nil {
			return nil, fmt.Errorf("error parsing file [%s], %v", path, err)
		}
		p.reset()
	}

	return p.env, nil
}

func Load(override bool, files ...string) error {
	env, err := Read(files...)
	if err != nil {
		return err
	}

	for k, v := range env {
		if !override && os.Getenv(k) != "" {
			continue
		}

		os.Setenv(k, v)
	}
	return nil
}

func Unmarshal(b []byte) (map[string]string, error) {
	p := parser{env: make(map[string]string), input: b}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p.env, nil
}
