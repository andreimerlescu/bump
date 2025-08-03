# Bump

![Bump](/bump.jpg)

This package is designed to take a string, like `v1.0.0` and allow you to `bump` it with the command line.

The `-in` argument is the **Input File** and it defaults to `./VERSION` from the _current working directory_ of where `bump` is being invoked.

The `bump` binary can intelligently bump `-in` files like `go.mod`, `package.json`, `pom.xml`, `Chart.yml` and `Dockerfile`. 

The `bump` binary can have its default runtime manipulated using **Environment Variables**. 

The `bump` binary leverages the `Exit Code 0` or `Exit Code 1` in order to use `bump` in a DevOps pipeline.

The `bump` binary offers you `-json` for _JSON_ Encoded output. 

This repository was built using _end to end testing_ including unit tests, _fuzz_ testing, benchmark and integration tests. 

## Demo of Bump

When you run `make test-cli`, you will get a generated `test-results/results.cli.md` that will contain the output captured
during this demo of the `bump` binary. Truly, I hope you enjoy using bump! It's a fun name, and I enjoyed writing it. 
I hope you find good use for it!

<video width="100%" controls poster="https://raw.githubusercontent.com/andreimerlescu/bump/refs/heads/master/bump.jpg">
  <source src="https://raw.githubusercontent.com/andreimerlescu/bump/refs/heads/master/demo.mp4" type="video/mp4">
  Your browser does not support the video tag.
</video>

### Version Format Priority

The `bump` package can parse multiple version formats. To ensure accuracy, it checks for formats in a specific order, from most complex to least complex. This prevents a detailed pre-release version from being incorrectly identified as a simpler one.

The priority is as follows:

1.  `v1.2.3-beta.4-alpha.5` (Beta with Alpha)
2.  `v1.2.3-alpha.1` (Alpha)
3.  `v1.2.3-beta.2` (Beta)
4.  `v1.2.3-rc.1` (Release Candidate)
5.  `v1.2.3-preview.7` (Preview)
6.  `v1.2.3` (Standard)

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
  -about
        show about
  -alpha
        alpha version bump
  -beta
        beta version bump
  -check
        check version file
  -env
        show env
  -in string
        input file (default "VERSION")
  -json
        use json version bump
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
  -set string
        set env to new value
  -v    show version
  -write
        writeInput version file
```

## About

```bash
bump -about
```

```log
Bump Version: v1.0.6
Usage:
  bump -check [-in=FILE]
  bump -fix [-write] [-in=FILE]
  bump -[major|minor|patch|alpha|beta|rc|preview] [-write] [-in=FILE] [-json]
Supported File Types:
  VERSION
  package.json
  pom.xml
  Chart.yaml
  Dockerfile
  go.mod
Defaults: 
  -in=VERSION [default: VERSION]
Environment Variables:
  BUMP_NO_PREVIEW=false
  BUMP_INIT_ON_NOT_FOUND=false
  BUMP_DEFAULT_INPUT=VERSION
  BUMP_NO_ALPHA=false
  BUMP_NO_BETA=false
  BUMP_NO_RC=false
  BUMP_ALWAYS_FIX=false
  BUMP_ALWAYS_WRITE=false
  BUMP_NEVER_FIX=false
  BUMP_NO_ALPHA_BETA=false

```

## Environment

You can use your **Environment Variables** to control the default runtime of the `bump` application. 

| `ENV`                |   Type   | Default   | Action                                                                   | 
|----------------------|:--------:|-----------|--------------------------------------------------------------------------|
| `BUMP_ALWAYS_WRITE`  |  `Bool`  | `false`   | When `true`, `-write` is `true` automatically and `-in` gets modified.   |
| `BUMP_DEFAULT_INPUT` | `String` | `<blank>` | When defined, a path to your default `VERSION` file should be used here. |
| `BUMP_NO_BETA`       |  `Bool`  | `false`   | When `true`, `-beta` will have no effect.                                |
| `BUMP_NO_ALPHA`      |  `Bool`  | `false`   | When `true`, `-alpha` will have no effect.                               |
| `BUMP_NO_ALPHA_BETA` |  `Bool`  | `false`   | When `true`, `-alpha` and `-beta` will have no effect.                   | 
| `BUMP_NO_RC`         |  `Bool`  | `false`   | When `true`, `-rc` will have no effect.                                  | 
| `BUMP_NO_PREVIEW`    |  `Bool`  | `false`   | When `true`, `-preview` will have no effect.                             |

It may be useful to enable to this on your environment. 

> *NOTE*: The `-json` mode _does NOT_ write to the `-in` file as JSON, but as plain text.

By disabling options like `BUMP_NO_ALPHA_BETA`, you can avoid having versions in your history that look like 
`v1.0.1-beta.3-alpha-3` from getting into your pipelines due to any invocations that combine the allowed `-beta` and 
`-alpha` flags during runtime.

You can use the `-env` argument to show results of the environment: 

```bash
bump -env
```

```log
BUMP_NO_PREVIEW=false
BUMP_ALWAYS_FIX=false
BUMP_DEFAULT_INPUT=VERSION
BUMP_NO_ALPHA=false
BUMP_NO_BETA=false
BUMP_NO_ALPHA_BETA=false
BUMP_NO_RC=false
BUMP_INIT_ON_NOT_FOUND=false
BUMP_ALWAYS_WRITE=false
BUMP_NEVER_FIX=false
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

The bump package has a comprehensive unit, fuzz, and benchmark test associated in the [version_test.go](/bump/version_test.go) 
file.

```bash
make test
# OR
make test-cli
make test-unit
make test-bench
make test-fuzz      # 3s run
# and not part of `make test`
make test-fuzz-long # 30s run
```

```log
Testing CLI... SUCCESS! Took 3 (s)! Wrote results.cli.md (size: 12K )
Testing Unit... SUCCESS! Took 1 (s)! Wrote results.unit.md ( size: 4.0K )
Testing Benchmark... SUCCESS! Took 2 (s)! Wrote results.benchmark.md ( size: 4.0K )
Testing Fuzz... SUCCESS! Took 33 (s)! Wrote results.fuzz.md ( size: 4.0K )
```

### Results

