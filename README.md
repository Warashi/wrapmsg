# wrapmsg
[![Coverage Status](https://coveralls.io/repos/github/Warashi/wrapmsg/badge.svg?branch=main)](https://coveralls.io/github/Warashi/wrapmsg?branch=main)

wrapmsg is Go code linter.
this enforces fmt.Errorf's message when you wrap error.

## Example
```go
// OK 👍🏻
if err := pkg.Cause(); err != nil {
  return fmt.Errorf("pkg.Cause: %w", err)
}

// NG 🙅
if err := pkg.Cause(); err != nil {
  return fmt.Errorf("cause failed: %w", err)
}
```

## Install
```sh
go install github.com/Warashi/wrapmsg/cmd/wrapmsg@latest
```

## Usage
You can use wrapmsg as vettool.
```sh
go vet -vettool=$(which wrapmsg) ./...
```

You can also build your linter with singlechecker or multichecker.
In this way, you can use `--fix` option to autocorrect.
```go
package main

import (
	"github.com/Warashi/wrapmsg"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(wrapmsg.Analyzer)
}
```
