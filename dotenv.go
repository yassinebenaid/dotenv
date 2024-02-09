// package "dotenv" provides functionality to consume .env files
package dotenv

import (
	"fmt"
	"os"
)

// Read reads all environment variables from the given files and return them as map.
//
//   - if `files` is empty, Read will default to loading .env in the current working directory.
//
// you can pass multiple files like this :
//
//	dotenv.Read("/path/to/file-1","/path/to/file-2")
//
// in this case, Read will read them in the given order, and you can reference keys from
// other files , just like they were the same file.
//
// file-1:
//
//	KEY_1=value
//
// file-2:
//
//	KEY_2=$KEY_1
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

// Load will read all environment variables from files and writes them to the ENV for this process.
// (typically using  os.Setenv()).
//
//   - if `override` is false, Load WILL NOT OVERRIDE the variables that are already set.
//   - if `files` is empty, Read will default to loading .env in the current working directory.
//
// you can pass multiple files like this :
//
//	dotenv.Load(false,"/path/to/file-1","/path/to/file-2")
//
// in this case, Load will read them in the given order, and you can reference keys from
// other files , just like they were the same file.
//
// file-1:
//
//	KEY_1=value
//
// file-2:
//
//	KEY_2=$KEY_1
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

// Unmarshal reads environment variables from b and returns them as map.
func Unmarshal(b []byte) (map[string]string, error) {
	p := parser{env: make(map[string]string), input: b}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p.env, nil
}
