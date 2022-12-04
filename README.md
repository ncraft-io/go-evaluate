# go-evaluate

this toolkit can evaluate and update the value in value declaration based on go generate.

## Reason

if we need to define some value that should never ever change during the application runtime, 
currently, we can set the value at linking stage by passing -ldflags "-X main.varname=varval". However this method separates the declaration and the evaluation in different files.

The method will cause mistake easily, and also not friendly for refactoring the value name. So, I think making the value declaration and the evaluation command in the same place will be a good idea.

## HowTo

1. add the fixed go:generate directive first

```go
package sample

//go:generate go run github.com/ncraft-io/go-evaluate/cmd/evaluate
```

2. add const value declaration with an empty initial value, and add the go:evaluate directive above it

```go
package sample

//go:evaluate date "+%Y-%m-%d %H:%M:%S %Z"
const BuildTime = ""
```
