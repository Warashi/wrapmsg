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
