# Death From Below

Death from below is a CLI tool to search and replace lines from files that
matches the glob pattern provided in the directory provided.

## Build
Run `go build` to build the program

## Test
Run `go test` to run the unit tests

## How to use
The tool takes four arguments:
  1. Source directory to search files for
  2. Glob pattern to match to files to search and replace lines for
  3. The line you want to replace
  4. The line you want to replace to

**NOTE:** The glob pattern paramters must be in double quotes

```sh
death-from-below example-dir "*.txt" Hello Goodbye
```

## Limitations
Death from below will recursively look through all folders, there is not way to
exclude folders.

## Note
There are existing solutions but this was made for learning and personal purposes.

I name projects that reference League of Legends. The name "Death From Below"
comes from the champion Pyke's ultimate name.
