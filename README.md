# Bump

This package is designed to take a string, like `v1.0.0` and allow you to `bump` it with the command line.

## Installation

```bash
go install github.com:andreimerlescu/bump@latest
which bump
```

Alternative method: 

### macOS Silicon

```bash
[ ! -d ~/bin ] && mkdir -p ~/bin
curl -sL https://github.com/andreimerlescu/bump/releases/download/v1.0.0/bump-darwin-arm64 -o ~/bin/bump
chmod +x ~/bin/bump
```

### macOS Intel

```bash
[ ! -d ~/bin ] && mkdir -p ~/bin
curl -sL https://github.com/andreimerlescu/bump/releases/download/v1.0.0/bump-darwin-amd64 -o ~/bin/bump
chmod +x ~/bin/bump
```

### Linux amd64

```bash
[ ! -d ~/bin ] && mkdir -p ~/bin
curl -sL https://github.com/andreimerlescu/bump/releases/download/v1.0.0/bump-linux-amd64 -o ~/bin/bump
chmod +x ~/bin/bump
```

### Windows (bash)

```bash
[ ! -d ~/bin ] && mkdir -p ~/bin
curl -sL https://github.com/andreimerlescu/bump/releases/download/v1.0.0/bump.exe -o ~/bin/bump.exe
```

## Usage

```bash
Usage of bump:
  -alpha
        alpha version bump
  -beta
        beta version bump
  -check
        check version file
  -in string
        input file (default "VERSION")
  -major
        major version bump
  -minor
        minor version bump
  -patch
        patch version bump
  -preview
        preview version bump
  -rc
        rc version bump
  -v    show version
  -write
        writeInput version file
```

## Examples

### Bumping Patch 

```bash
cd ~/work
git clone git@github.com:ProjectApario/genwordpass.git
cd genwordpass
bump -check # check version (error code indicates validity, stdout always raw file contents of -i)
v1.0.2

bump -patch # see new version first
Bumped v1.0.2 → v1.0.3

bump -patch -w # save to file
Bumped v1.0.2 → v1.0.3 (saved to VERSION)

bump -check
v1.0.3
```

## Development

You can clone the repository if you want to. 

```bash
cd ~/work
git clone git@github.com:andreimerlescu/bump.git
cd bump
```

## Building

```bash
make all
# OR
make darwin-arm64
make darwin-amd64 
make linux-arm64
make linux-amd64 
make windows-amd64 
```

```log
Summary generated: summaries/summary.2025.08.17.17.17.17.UTC.md
Building binary target: darwin/amd64
Building binary target: darwin/arm64
Building binary target: linux/amd64
Building binary target: linux/arm64
Building binary target: windows/amd64
NEW: /Users/andrei/go/bin/bump
```

## Testing

The bump package has a comprehensive unit, fuzz, and benchmark test associated in the [bump_test.go](/bump/bump_test.go) 
file.

```bash
make test
```

```log
Testing Unit... SUCCESS! Took 1 (s)! Wrote results.unit.md ( size: 4.0K )
Testing Benchmark... SUCCESS! Took 2 (s)! Wrote results.benchmark.md ( size: 4.0K )
Testing Fuzz... SUCCESS! Took 33 (s)! Wrote results.fuzz.md ( size: 4.0K )
```

