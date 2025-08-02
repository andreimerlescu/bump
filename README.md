# Bump

This package is designed to take a string, like `v1.0.0` and allow you to `bump` it with the command line.

The `-in` argument is the **Input File** and it defaults to `./VERSION` from the _current working directory_ of where `bump` is being invoked.

The `bump` binary can intelligently bump `-in` files like `go.mod`, `package.json`, `pom.xml`, `Chart.yml` and `Dockerfile`. 

The `bump` binary can have its default runtime manipulated using **Environment Variables**. 

The `bump` binary leverages the `Exit Code 0` or `Exit Code 1` in order to use `bump` in a DevOps pipeline.

The `bump` binary offers you `-json` for _JSON_ Encoded output. 

This repository was built using _end to end testing_ including unit tests, _fuzz_ testing, benchmark and integration tests. 

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
Bump Version: v1.0.5
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

The bump package has a comprehensive unit, fuzz, and benchmark test associated in the [bump_test.go](/bump/bump_test.go) 
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
Testing CLI... SUCCESS! Took 2 (s)! Wrote results.cli.md (size: 12K )
Testing Unit... SUCCESS! Took 1 (s)! Wrote results.unit.md ( size: 4.0K )
Testing Benchmark... SUCCESS! Took 2 (s)! Wrote results.benchmark.md ( size: 4.0K )
Testing Fuzz... SUCCESS! Took 33 (s)! Wrote results.fuzz.md ( size: 4.0K )
```

### Results

### `/Users/andrei/work/bump/test-results/results.cli.md`

Test results captured at 2025-08-01 12:30:40.

```log
Preparing test env...
run.01(echo "--- SCENARIO ONE ---") - SUCCESS
    output: --- SCENARIO ONE ---
run.02(echo "v1.0.0" > VERSION) - SUCCESS
run.03(bump -check) - SUCCESS
    output: v1.0.0
run.04(cat VERSION) - SUCCESS
    output: v1.0.0
run.05(bump -alpha) - SUCCESS
    output: Bumped v1.0.0 → v1.0.0-alpha.1
run.06(cat VERSION) - SUCCESS
    output: v1.0.0
run.07(bump -alpha -write) - SUCCESS
    output: Bumped v1.0.0 → v1.0.0-alpha.1 (saved to VERSION)
run.08(cat VERSION) - SUCCESS
    output: v1.0.0-alpha.1
run.09(rm VERSION) - SUCCESS
run.10(echo "--- SCENARIO TWO ---") - SUCCESS
    output: --- SCENARIO TWO ---
run.11(echo "v1.0.0-alpha.0" > VERSION) - SUCCESS
run.12(bump -check) - SUCCESS
    output: v1.0.0
run.13(cat VERSION) - SUCCESS
    output: v1.0.0-alpha.0
run.14(bump -alpha) - SUCCESS
    output: Bumped v1.0.0 → v1.0.0-alpha.1
run.15(cat VERSION) - SUCCESS
    output: v1.0.0-alpha.0
run.16(bump -alpha -write) - SUCCESS
    output: Bumped v1.0.0 → v1.0.0-alpha.1 (saved to VERSION)
run.17(cat VERSION) - SUCCESS
    output: v1.0.0-alpha.1
run.18(bump -patch) - SUCCESS
    output: Bumped v1.0.0-alpha.1 → v1.0.1
run.19(cat VERSION) - SUCCESS
    output: v1.0.0-alpha.1
run.20(bump -patch -write) - SUCCESS
    output: Bumped v1.0.0-alpha.1 → v1.0.1 (saved to VERSION)
run.21(cat VERSION) - SUCCESS
    output: v1.0.1
run.22(bump -major -write) - SUCCESS
    output: Bumped v1.0.1 → v2.0.0 (saved to VERSION)
run.23(cat VERSION) - SUCCESS
    output: v2.0.0
run.24(bump -preview -write) - SUCCESS
    output: Bumped v2.0.0 → v2.0.0-preview.1 (saved to VERSION)
run.25(rm VERSION) - SUCCESS
run.26(echo "--- SCENARIO THREE ---") - SUCCESS
    output: --- SCENARIO THREE ---
run.27(echo "1.25" > VERSION) - SUCCESS
run.28(bump -fix) - SUCCESS
    output: No bump operation specified. Use -major, -minor, -patch, etc., to bump the version.
Current version is: v1.25.0
run.29(bump -fix -write) - SUCCESS
    output: Fixed and saved version v1.25.0 to VERSION
run.30(rm VERSION) - SUCCESS
run.31(echo "--- SCENARIO FOUR ---") - SUCCESS
    output: --- SCENARIO FOUR ---
run.32(echo "v1.17.7-beta.6" > VERSION) - SUCCESS
run.33(bump -check -fix) - SUCCESS
    output: v1.17.7-beta.6
run.34(cat VERSION) - SUCCESS
    output: v1.17.7-beta.6
run.35(bump -check -fix -write) - SUCCESS
    output: v1.17.7-beta.6
run.36(cat VERSION) - SUCCESS
    output: v1.17.7-beta.6
run.37(rm VERSION) - SUCCESS
run.38(echo "--- SCENARIO FIVE ---") - SUCCESS
    output: --- SCENARIO FIVE ---
run.39(echo "module testApp-bump-test-passes" > go.mod) - SUCCESS
run.40(echo "" >> go.mod) - SUCCESS
run.41(echo "go 1.24" >> go.mod) - SUCCESS
run.42(cat go.mod) - SUCCESS
    output: module testApp-bump-test-passes

go 1.24
run.43(bump -in go.mod -fix) - SUCCESS
    output: No bump operation specified. Use -major, -minor, -patch, etc., to bump the version.
Current version is: 1.24
run.44(bump -in go.mod -fix -write) - SUCCESS
    output: Fixed and saved version 1.24 to go.mod
run.45(cat go.mod) - SUCCESS
    output: module testApp-bump-test-passes

go 1.24
run.46(rm go.mod) - SUCCESS
run.47(echo "--- SCENARIO SIX ---") - SUCCESS
    output: --- SCENARIO SIX ---
run.48(echo v1.0.0 > VERSION) - SUCCESS
run.49(bump -json -check) - SUCCESS
    output: {
  "version": "v1.0.0"
}
run.50(cat VERSION) - SUCCESS
    output: v1.0.0
run.51(bump -json -beta) - SUCCESS
    output: {
  "major": 1,
  "minor": 0,
  "patch": 0,
  "alpha": 0,
  "beta": 1,
  "rc": 0,
  "preview": 0,
  "version": "v1.0.0-beta.1"
}
run.52(cat VERSION) - SUCCESS
    output: v1.0.0
run.53(bump -json -beta -write) - SUCCESS
    output: {
  "major": 1,
  "minor": 0,
  "patch": 0,
  "alpha": 0,
  "beta": 1,
  "rc": 0,
  "preview": 0,
  "version": "v1.0.0-beta.1"
}
run.54(cat VERSION) - SUCCESS
    output: v1.0.0-beta.1
run.55(rm VERSION) - SUCCESS
run.56(echo "--- SCENARIO SEVEN: package.json ---") - SUCCESS
    output: --- SCENARIO SEVEN: package.json ---
run.57(echo eyJuYW1lIjoidGVzdCIsInZlcnNpb24iOiIxLjIuMyJ9 | base64 -d | tee package.json > /dev/null) - SUCCESS
run.58(cat package.json) - SUCCESS
    output: {"name":"test","version":"1.2.3"}
run.59(bump -in package.json -fix) - SUCCESS
    output: i can see this package.json
i can see this package.json
No bump operation specified. Use -major, -minor, -patch, etc., to bump the version.
Current version is: 1.2.3
run.60(bump -in package.json -fix -write) - SUCCESS
    output: i can see this package.json
i can see this package.json
Fixed and saved version 1.2.3 to package.json
run.61(bump -in package.json -patch) - SUCCESS
    output: i can see this package.json
Bumped 1.2.3 → 1.2.4
run.62(bump -in package.json -patch -write) - SUCCESS
    output: i can see this package.json
Bumped 1.2.3 → 1.2.4 (saved to package.json)
run.63(grep '"version": "1.2.4"' package.json) - SUCCESS
    output:   "version": "1.2.4"
run.64(bump -in package.json -json -minor) - SUCCESS
    output: i can see this package.json
{
  "major": 1,
  "minor": 3,
  "patch": 0,
  "alpha": 0,
  "beta": 0,
  "rc": 0,
  "preview": 0,
  "version": "1.3.0"
}
run.65(bump -in package.json -minor -write) - SUCCESS
    output: i can see this package.json
Bumped 1.2.4 → 1.3.0 (saved to package.json)
run.66(grep '"version": "1.3.0"' package.json) - SUCCESS
    output:   "version": "1.3.0"
run.67(rm package.json) - SUCCESS
run.68(echo "--- SCENARIO EIGHT: go.mod ---") - SUCCESS
    output: --- SCENARIO EIGHT: go.mod ---
run.69(echo 'module myapp' > go.mod && echo 'go 1.21' >> go.mod) - SUCCESS
run.70(bump -in go.mod -check) - SUCCESS
    output: 1.21
run.71(bump -in go.mod -minor) - SUCCESS
    output: Bumped 1.21 → 1.22
run.72(bump -in go.mod -minor -write) - SUCCESS
    output: Bumped 1.21 → 1.22 (saved to go.mod)
run.73(grep 'go 1.22' go.mod) - SUCCESS
    output: go 1.22
run.74(rm go.mod) - SUCCESS
run.75(echo "--- SCENARIO NINE: Dockerfile ---") - SUCCESS
    output: --- SCENARIO NINE: Dockerfile ---
run.76(echo 'LABEL version="v3.2.1"' > Dockerfile) - SUCCESS
run.77(bump -in Dockerfile -check) - SUCCESS
    output: v3.2.1
run.78(bump -in Dockerfile -patch) - SUCCESS
    output: Bumped v3.2.1 → v3.2.2
run.79(bump -in Dockerfile -patch -write) - SUCCESS
    output: Bumped v3.2.1 → v3.2.2 (saved to Dockerfile)
run.80(grep 'LABEL version="v3.2.2"' Dockerfile) - SUCCESS
    output: LABEL version="v3.2.2"
run.81(rm Dockerfile) - SUCCESS
run.82(echo "--- SCENARIO TEN: Chart.yaml (Helm) ---") - SUCCESS
    output: --- SCENARIO TEN: Chart.yaml (Helm) ---
run.83(echo 'YXBpVmVyc2lvbjogdjIKbmFtZTogbXljaGFydAp2ZXJzaW9uOiAwLjEuMAo=' | base64 -d | tee Chart.yaml > /dev/null) - SUCCESS
run.84(bump -in Chart.yaml -check) - SUCCESS
    output: 0.1.0
run.85(bump -in Chart.yaml -patch) - SUCCESS
    output: Bumped 0.1.0 → 0.1.1
run.86(bump -in Chart.yaml -patch -write) - SUCCESS
    output: Bumped 0.1.0 → 0.1.1 (saved to Chart.yaml)
run.87(grep 'version: 0.1.1' Chart.yaml) - SUCCESS
    output: version: 0.1.1
run.88(rm Chart.yaml) - SUCCESS
run.89(echo "--- SCENARIO ELEVEN: pom.xml (Maven) ---") - SUCCESS
    output: --- SCENARIO ELEVEN: pom.xml (Maven) ---
run.90(echo '<project><version>2.2.2</version></project>' > pom.xml) - SUCCESS
run.91(bump -in pom.xml -check) - SUCCESS
    output: 2.2.2
run.92(bump -in pom.xml -patch) - SUCCESS
    output: Bumped 2.2.2 → 2.2.3
run.93(bump -in pom.xml -patch -write) - SUCCESS
    output: Bumped 2.2.2 → 2.2.3 (saved to pom.xml)
run.94(grep '<version>2.2.3</version>' pom.xml) - SUCCESS
    output: <project><version>2.2.3</version></project>
run.95(rm pom.xml) - SUCCESS
run.96(echo "--- SCENARIO TWELVE: Environment Variables ---") - SUCCESS
    output: --- SCENARIO TWELVE: Environment Variables ---
run.97(echo "v5.5.5" > VERSION) - SUCCESS
run.98(BUMP_ALWAYS_WRITE=true bump -env) - SUCCESS
    output: BUMP_NO_BETA=false
BUMP_NO_ALPHA_BETA=false
BUMP_NO_RC=false
BUMP_INIT_ON_NOT_FOUND=false
BUMP_NO_PREVIEW=false
BUMP_ALWAYS_FIX=false
BUMP_ALWAYS_WRITE=true
BUMP_DEFAULT_INPUT=VERSION
BUMP_NEVER_FIX=false
BUMP_NO_ALPHA=false
run.99(BUMP_ALWAYS_WRITE=true bump -patch) - SUCCESS
    output: Bumped v5.5.5 → v5.5.6 (saved to VERSION)
run.100(grep 'v5.5.6' VERSION) - SUCCESS
    output: v5.5.6
run.101(BUMP_DEFAULT_INPUT=VERSION bump -minor) - SUCCESS
    output: Bumped v5.5.6 → v5.6.0
run.102(BUMP_DEFAULT_INPUT=VERSION bump -minor -write) - SUCCESS
    output: Bumped v5.5.6 → v5.6.0 (saved to VERSION)
run.103(grep 'v5.6.0' VERSION) - SUCCESS
    output: v5.6.0
run.104(rm VERSION) - SUCCESS
All 104 tests PASS!
```

### `/Users/andrei/work/bump/test-results/results.fuzz.md`

Test results captured at 2025-07-28 07:53:25.

```log
fuzz: elapsed: 0s, gathering baseline coverage: 0/182 completed
fuzz: elapsed: 0s, gathering baseline coverage: 182/182 completed, now fuzzing with 14 workers
fuzz: elapsed: 3s, execs: 657266 (219012/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 6s, execs: 1147687 (163521/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 9s, execs: 1543842 (132017/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 12s, execs: 1724503 (60219/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 15s, execs: 1798638 (24712/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 18s, execs: 1840749 (14037/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 21s, execs: 1840749 (0/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 24s, execs: 1840749 (0/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 27s, execs: 1840749 (0/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 30s, execs: 1840749 (0/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 33s, execs: 1840749 (0/sec), new interesting: 0 (total: 182)
fuzz: elapsed: 34s, execs: 1840749 (0/sec), new interesting: 0 (total: 182)
PASS
ok  	github.com/andreimerlescu/bump/bump	34.297s
```

### `/Users/andrei/work/bump/test-results/results.unit.md`

Test results captured at 2025-08-01 12:30:41.

```log
ok  	github.com/andreimerlescu/bump/bump	0.162s
```



