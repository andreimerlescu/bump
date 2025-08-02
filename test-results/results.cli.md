### `/Users/andrei/work/bump/test-results/results.cli.md` 

 Test results captured at 2025-08-02 11:40:05. 

```log
Preparing test env...
andrei@bump:test.sh ⚡ Test #1 ⇒  echo "v1.0.0" > VERSION
andrei@bump:test.sh ⚡ Test #2 ⇒  bump -check
v1.0.0
andrei@bump:test.sh ⚡ Test #3 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #4 ⇒  bump -alpha
Bumped v1.0.0 → v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #5 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #6 ⇒  bump -alpha -write
Bumped v1.0.0 → v1.0.0-alpha.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #7 ⇒  cat VERSION
v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #8 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #9 ⇒  echo "v1.0.0-alpha.0" > VERSION
andrei@bump:test.sh ⚡ Test #10 ⇒  bump -check
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #11 ⇒  cat VERSION
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #12 ⇒  bump -alpha
Bumped v1.0.0-alpha.0 → v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #13 ⇒  cat VERSION
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #14 ⇒  bump -alpha -write
Bumped v1.0.0-alpha.0 → v1.0.0-alpha.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #15 ⇒  cat VERSION
v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #16 ⇒  bump -patch
Bumped v1.0.0-alpha.1 → v1.0.1
andrei@bump:test.sh ⚡ Test #17 ⇒  cat VERSION
v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #18 ⇒  bump -patch -write
Bumped v1.0.0-alpha.1 → v1.0.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #19 ⇒  cat VERSION
v1.0.1
andrei@bump:test.sh ⚡ Test #20 ⇒  bump -major -write
Bumped v1.0.1 → v2.0.0 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #21 ⇒  cat VERSION
v2.0.0
andrei@bump:test.sh ⚡ Test #22 ⇒  bump -preview -write
Bumped v2.0.0 → v2.0.0-preview.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #23 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #24 ⇒  echo "1.25" > VERSION
andrei@bump:test.sh ⚡ Test #25 ⇒  bump -fix
Bumped v1.25.0 → v1.25.0
andrei@bump:test.sh ⚡ Test #26 ⇒  bump -fix -write
Bumped v1.25.0 → v1.25.0 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #27 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #28 ⇒  echo "v1.17.7-beta.6" > VERSION
andrei@bump:test.sh ⚡ Test #29 ⇒  bump -check -fix
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #30 ⇒  cat VERSION
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #31 ⇒  bump -check -fix -write
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #32 ⇒  cat VERSION
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #33 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #34 ⇒  echo "module testApp-" > go.mod
andrei@bump:test.sh ⚡ Test #35 ⇒  echo "" >> go.mod
andrei@bump:test.sh ⚡ Test #36 ⇒  echo "go 1.24" >> go.mod
andrei@bump:test.sh ⚡ Test #37 ⇒  cat go.mod
module testApp-

go 1.24
andrei@bump:test.sh ⚡ Test #38 ⇒  bump -in go.mod -fix
Bumped 1.24 → 1.24
andrei@bump:test.sh ⚡ Test #39 ⇒  bump -in go.mod -fix -write
Bumped 1.24 → 1.24 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #40 ⇒  cat go.mod
module testApp-

go 1.24.0
andrei@bump:test.sh ⚡ Test #41 ⇒  rm go.mod
andrei@bump:test.sh ⚡ Test #42 ⇒  echo v1.0.0 > VERSION
andrei@bump:test.sh ⚡ Test #43 ⇒  bump -json -check
{
  "version": "v1.0.0"
}
andrei@bump:test.sh ⚡ Test #44 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #45 ⇒  bump -json -beta
{
  "major": 1,
  "minor": 0,
  "patch": 0,
  "alpha": 0,
  "beta": 1,
  "rc": 0,
  "preview": 0,
  "version": "v1.0.0-alpha.0"
}
andrei@bump:test.sh ⚡ Test #46 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #47 ⇒  bump -json -beta -write
{
  "major": 1,
  "minor": 0,
  "patch": 0,
  "alpha": 0,
  "beta": 1,
  "rc": 0,
  "preview": 0,
  "version": "v1.0.0-alpha.0"
}
andrei@bump:test.sh ⚡ Test #48 ⇒  cat VERSION
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #49 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #50 ⇒  echo 'eyJuYW1lIjoidGVzdCIsInZlcnNpb24iOiIxLjIuMyJ9' | base64 -d | tee package.json > /dev/null
andrei@bump:test.sh ⚡ Test #51 ⇒  cat package.json
{"name":"test","version":"1.2.3"}
andrei@bump:test.sh ⚡ Test #52 ⇒  bump -in package.json -fix
Bumped 1.2.3 → 1.2.3
andrei@bump:test.sh ⚡ Test #53 ⇒  bump -in package.json -fix -write
Bumped 1.2.3 → 1.2.3 (saved to package.json)
andrei@bump:test.sh ⚡ Test #54 ⇒  bump -in package.json -patch
Bumped 1.2.3 → 1.2.4
andrei@bump:test.sh ⚡ Test #55 ⇒  bump -in package.json -patch -write
Bumped 1.2.3 → 1.2.4 (saved to package.json)
andrei@bump:test.sh ⚡ Test #56 ⇒  grep '"version": "1.2.4"' package.json
  "version": "1.2.4"
andrei@bump:test.sh ⚡ Test #57 ⇒  bump -in package.json -json -minor
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
andrei@bump:test.sh ⚡ Test #58 ⇒  bump -in package.json -minor -write
Bumped 1.2.4 → 1.3.0 (saved to package.json)
andrei@bump:test.sh ⚡ Test #59 ⇒  grep '"version": "1.3.0"' package.json
  "version": "1.3.0"
andrei@bump:test.sh ⚡ Test #60 ⇒  rm package.json
andrei@bump:test.sh ⚡ Test #61 ⇒  echo 'module myapp' > go.mod && echo 'go 1.24' >> go.mod
andrei@bump:test.sh ⚡ Test #62 ⇒  bump -in go.mod -check
1.24
andrei@bump:test.sh ⚡ Test #63 ⇒  bump -in go.mod -patch
Bumped 1.24 → v1.24.1
andrei@bump:test.sh ⚡ Test #64 ⇒  bump -in go.mod -patch -write
Bumped 1.24 → v1.24.1 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #65 ⇒  bump -in go.mod -patch -write
Bumped 1.24.1 → 1.24.2 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #66 ⇒  bump -in go.mod -patch -write
Bumped 1.24.2 → 1.24.3 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #67 ⇒  bump -in go.mod -patch -write
Bumped 1.24.3 → 1.24.4 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #68 ⇒  bump -in go.mod -patch -write
Bumped 1.24.4 → 1.24.5 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #69 ⇒  grep 'go 1.24.5' go.mod
go 1.24.5
andrei@bump:test.sh ⚡ Test #70 ⇒  rm go.mod
andrei@bump:test.sh ⚡ Test #71 ⇒  echo 'LABEL version="v3.2.1"' > Dockerfile
andrei@bump:test.sh ⚡ Test #72 ⇒  bump -in Dockerfile -check
v3.2.1
andrei@bump:test.sh ⚡ Test #73 ⇒  bump -in Dockerfile -patch
Bumped v3.2.1 → v3.2.2
andrei@bump:test.sh ⚡ Test #74 ⇒  bump -in Dockerfile -patch -write
Bumped v3.2.1 → v3.2.2 (saved to Dockerfile)
andrei@bump:test.sh ⚡ Test #75 ⇒  grep 'LABEL version="v3.2.2"' Dockerfile
LABEL version="v3.2.2"
andrei@bump:test.sh ⚡ Test #76 ⇒  rm Dockerfile
andrei@bump:test.sh ⚡ Test #77 ⇒  echo 'YXBpVmVyc2lvbjogdjIKbmFtZTogbXljaGFydAp2ZXJzaW9uOiAwLjEuMAo=' | base64 -d | tee Chart.yaml > /dev/null
andrei@bump:test.sh ⚡ Test #78 ⇒  bump -in Chart.yaml -check
0.1.0
andrei@bump:test.sh ⚡ Test #79 ⇒  bump -in Chart.yaml -patch
Bumped 0.1.0 → 0.1.1
andrei@bump:test.sh ⚡ Test #80 ⇒  bump -in Chart.yaml -patch -write
Bumped 0.1.0 → 0.1.1 (saved to Chart.yaml)
andrei@bump:test.sh ⚡ Test #81 ⇒  grep 'version: 0.1.1' Chart.yaml
version: 0.1.1
andrei@bump:test.sh ⚡ Test #82 ⇒  rm Chart.yaml
andrei@bump:test.sh ⚡ Test #83 ⇒  echo '<project><version>2.2.2</version></project>' > pom.xml
andrei@bump:test.sh ⚡ Test #84 ⇒  bump -in pom.xml -check
2.2.2
andrei@bump:test.sh ⚡ Test #85 ⇒  bump -in pom.xml -patch
Bumped 2.2.2 → 2.2.3
andrei@bump:test.sh ⚡ Test #86 ⇒  bump -in pom.xml -patch -write
Bumped 2.2.2 → 2.2.3 (saved to pom.xml)
andrei@bump:test.sh ⚡ Test #87 ⇒  grep '<version>2.2.3</version>' pom.xml
<project><version>2.2.3</version></project>
andrei@bump:test.sh ⚡ Test #88 ⇒  rm pom.xml
andrei@bump:test.sh ⚡ Test #89 ⇒  echo "v5.5.5" > VERSION
andrei@bump:test.sh ⚡ Test #90 ⇒  BUMP_ALWAYS_WRITE=true bump -env
BUMP_NO_PREVIEW=false
BUMP_INIT_ON_NOT_FOUND=false
BUMP_ALWAYS_FIX=false
BUMP_DEFAULT_INPUT=VERSION
BUMP_NEVER_FIX=false
BUMP_NO_BETA=false
BUMP_ALWAYS_WRITE=true
BUMP_NO_ALPHA=false
BUMP_NO_ALPHA_BETA=false
BUMP_NO_RC=false
andrei@bump:test.sh ⚡ Test #91 ⇒  BUMP_ALWAYS_WRITE=true bump -patch
Bumped v5.5.5 → v5.5.6 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #92 ⇒  grep 'v5.5.6' VERSION
v5.5.6
andrei@bump:test.sh ⚡ Test #93 ⇒  BUMP_DEFAULT_INPUT=VERSION bump -minor
Bumped v5.5.6 → v5.6.0
andrei@bump:test.sh ⚡ Test #94 ⇒  BUMP_DEFAULT_INPUT=VERSION bump -minor -write
Bumped v5.5.6 → v5.6.0 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #95 ⇒  grep 'v5.6.0' VERSION
v5.6.0
andrei@bump:test.sh ⚡ Test #96 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #97 ⇒  bump -parse v1.2.3-alpha.4 -init
Initialized v1.2.3-alpha.4
andrei@bump:test.sh ⚡ Test #98 ⇒  cat VERSION
v1.2.3-alpha.4
andrei@bump:test.sh ⚡ Test #99 ⇒  bump -parse v2.3.4-alpha.5 -write
Parsed v2.3.4-alpha.5 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #100 ⇒  cat VERSION
v2.3.4-alpha.5
andrei@bump:test.sh ⚡ Test #101 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #102 ⇒  bump -parse v3.4.5-alpha.6 -init
Initialized v3.4.5-alpha.6
andrei@bump:test.sh ⚡ Test #103 ⇒  cat VERSION
v3.4.5-alpha.6
andrei@bump:test.sh ⚡ Test #104 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #105 ⇒  echo 'ewogICJuYW1lIjogIm15X3BhY2thZ2UiLAogICJkZXNjcmlwdGlvbiI6ICJtYWtlIHlvdXIgcGFja2FnZSBlYXNpZXIgdG8gZmluZCBvbiB0aGUgbnBtIHdlYnNpdGUiLAogICJ2ZXJzaW9uIjogIjEuMC4wIiwKICAic2NyaXB0cyI6IHsKICAgICJ0ZXN0IjogImVjaG8gXCJFcnJvcjogbm8gdGVzdCBzcGVjaWZpZWRcIiAmJiBleGl0IDEiCiAgfSwKICAicmVwb3NpdG9yeSI6IHsKICAgICJ0eXBlIjogImdpdCIsCiAgICAidXJsIjogImh0dHBzOi8vZ2l0aHViLmNvbS9tb25hdGhlb2N0b2NhdC9teV9wYWNrYWdlLmdpdCIKICB9LAogICJrZXl3b3JkcyI6IFtdLAogICJhdXRob3IiOiAiIiwKICAibGljZW5zZSI6ICJJU0MiLAogICJidWdzIjogewogICAgInVybCI6ICJodHRwczovL2dpdGh1Yi5jb20vbW9uYXRoZW9jdG9jYXQvbXlfcGFja2FnZS9pc3N1ZXMiCiAgfSwKICAiaG9tZXBhZ2UiOiAiaHR0cHM6Ly9naXRodWIuY29tL21vbmF0aGVvY3RvY2F0L215X3BhY2thZ2UiCn0=' | base64 -d | tee package.json > /dev/null
andrei@bump:test.sh ⚡ Test #106 ⇒  cat package.json
{
  "name": "my_package",
  "description": "make your package easier to find on the npm website",
  "version": "1.0.0",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "repository": {
    "type": "git",
    "url": "https://github.com/monatheoctocat/my_package.git"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "bugs": {
    "url": "https://github.com/monatheoctocat/my_package/issues"
  },
  "homepage": "https://github.com/monatheoctocat/my_package"
}
andrei@bump:test.sh ⚡ Test #107 ⇒  bump -in package.json -check
1.0.0
andrei@bump:test.sh ⚡ Test #108 ⇒  bump -in package.json -fix
Bumped 1.0.0 → 1.0.0
andrei@bump:test.sh ⚡ Test #109 ⇒  bump -in package.json -fix -write
Bumped 1.0.0 → 1.0.0 (saved to package.json)
andrei@bump:test.sh ⚡ Test #110 ⇒  bump -in package.json -patch
Bumped 1.0.0 → 1.0.1
andrei@bump:test.sh ⚡ Test #111 ⇒  bump -in package.json -patch -write
Bumped 1.0.0 → 1.0.1 (saved to package.json)
andrei@bump:test.sh ⚡ Test #112 ⇒  grep '"version": "1.0.1"' package.json
  "version": "1.0.1"
andrei@bump:test.sh ⚡ Test #113 ⇒  bump -in package.json -json -minor
{
  "major": 1,
  "minor": 1,
  "patch": 0,
  "alpha": 0,
  "beta": 0,
  "rc": 0,
  "preview": 0,
  "version": "1.1.0"
}
andrei@bump:test.sh ⚡ Test #114 ⇒  bump -in package.json -minor -write
Bumped 1.0.1 → 1.1.0 (saved to package.json)
andrei@bump:test.sh ⚡ Test #115 ⇒  grep '"version": "1.1.0"' package.json
  "version": "1.1.0"
andrei@bump:test.sh ⚡ Test #116 ⇒  cat package.json
{
  "author": "",
  "bugs": {
    "url": "https://github.com/monatheoctocat/my_package/issues"
  },
  "description": "make your package easier to find on the npm website",
  "homepage": "https://github.com/monatheoctocat/my_package",
  "keywords": [],
  "license": "ISC",
  "name": "my_package",
  "repository": {
    "type": "git",
    "url": "https://github.com/monatheoctocat/my_package.git"
  },
  "scripts": {
    "test": "echo \"Error: no test specified\" \u0026\u0026 exit 1"
  },
  "version": "1.1.0"
}
andrei@bump:test.sh ⚡ Test #117 ⇒  rm package.json
All 117 tests PASS!
```

