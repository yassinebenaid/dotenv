# dotenv

Go Library to read environment variables from .env file.

[![Tests](https://github.com/yassinebenaid/dotenv/actions/workflows/tests.yaml/badge.svg)](https://github.com/yassinebenaid/dotenv/actions/workflows/tests.yaml)
[![Version](https://badge.fury.io/gh/yassinebenaid%2Fdotenv.svg)](https://badge.fury.io/gh/yassinebenaid%2Fdotenv)
[![AGPL License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENCE)

# Usage

install the library :

```bash
go get -u github.com/yassinebenaid/dotenv
```

_.env_

```bash
# comments are supported

PORT=8080
HOST="example.com:$PORT" # inline comments too

```

then in the beginning of your program, preferrably in the `main` function:

_main.go_

```go
package main

import (
    "os"
    "fmt"

    "github.com/yassinebenaid/dotenv"
)

func main(){
    err := dotenv.Load(false,".env")
    if err != nil{
        panic(err)
    }

    fmt.Println(os.Getenv("PORT")) // 8080
    fmt.Println(os.Getenv("HOST")) // example.com:8080


    // if no file path provided, Load will default
    // to loading .env in current working directory
    err = dotenv.Load(false)
    if err != nil{
        panic(err)
    }

    fmt.Println(os.Getenv("PORT")) // 8080
    fmt.Println(os.Getenv("HOST")) // example.com:8080

    // ...
}
``
```

# API

## dotenv.Load

_.env_

```bash
PORT=8080
HOST="example.com:$PORT"

```

_main.go_

```go
package main

import (
    "os"
    "fmt"

    "github.com/yassinebenaid/dotenv"
)

func main(){
    os.Setenv("PORT","5000")

    err := dotenv.Load(false) // will not override
    if err != nil{
        panic(err)
    }

    fmt.Println(os.Getenv("PORT")) // 5000


    // if you want to override
     os.Setenv("PORT","5000")

    err = dotenv.Load(true) // will  override
    if err != nil{
        panic(err)
    }

    fmt.Println(os.Getenv("PORT")) // 8080

    // you can pass multiple files
    err = dotenv.Load(true,"./path/to/file-1","./path/to/file-2")
    if err != nil{
        panic(err)
    }
}
``
```

## dotenv.Read

_main.go_

```go
package main

import (
	"fmt"

	"github.com/yassinebenaid/dotenv"
)

func main() {
	env, err := dotenv.Read("/path/to/.env")
	if err != nil {
		panic(err)
	}

	fmt.Println(env["PORT"])

	// you can pass multiple files
	env, err = dotenv.Read("/path/to/.env-1", "/path/to/.env-2")
	if err != nil {
		panic(err)
	}

	// if not path provided, Read will default to loading .env in current
	// working directory
	env, err = dotenv.Read()
	if err != nil {
		panic(err)
	}

}

``
```

`

## dotenv.Unmarshal

_main.go_

```go
package main

import (
    "fmt"

    "github.com/yassinebenaid/dotenv"
)

func main(){
    input := []byte(`
        KEY_1=value
        KEY_2=value
    `)
    env,err := dotenv.Unmarshal(input)
    if err != nil{
        panic(err)
    }

    fmt.Println(env["KEY_1"])
}
``
```

for more details check [go doc](https://pkg.go.dev/pkg/github.com/yassinebenaid/dotenv)

# Contributing

Contributions are welcome, make sure to include tests or update exising ones to cover your changes.

_code changes without tests and references to peer dotenv implementations will not be accepted_

- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Added some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

# Licence

Check the [Licence file](./LICENCE) for details.
