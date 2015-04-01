# file2const

Generates a Go file containing the given files as strings. Useful for packaging in static files.

```
Usage: file2const <flags> [<inputFile:constantName>] <outputFile.go>
  -package="": name of package to give to file. Default to directory name
```

Given a file like follows:

```javascript
console.log("HEY")
```

This command: `file2const --package="example" file.js:value out.go`

Will write the following to out.go

```go
package example

const (
	value = "console.log(\"HEY\")\n"
)
```
