### `/Users/andrei/work/bump/test-results/results.cli.md` 

 Test results captured at 2025-08-01 21:59:55. 

```log
Preparing test env...
andrei@bump:test.sh ⚡ Test #1 ⇒  echo "--- SCENARIO ONE ---"
--- SCENARIO ONE ---
andrei@bump:test.sh ⚡ Test #2 ⇒  echo "v1.0.0" > VERSION
andrei@bump:test.sh ⚡ Test #3 ⇒  bump -check
v1.0.0
andrei@bump:test.sh ⚡ Test #4 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #5 ⇒  bump -alpha
Bumped v1.0.0 → v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #6 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #7 ⇒  bump -alpha -write
Bumped v1.0.0 → v1.0.0-alpha.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #8 ⇒  cat VERSION
v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #9 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #10 ⇒  echo "--- SCENARIO TWO ---"
--- SCENARIO TWO ---
andrei@bump:test.sh ⚡ Test #11 ⇒  echo "v1.0.0-alpha.0" > VERSION
andrei@bump:test.sh ⚡ Test #12 ⇒  bump -check
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #13 ⇒  cat VERSION
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #14 ⇒  bump -alpha
Bumped v1.0.0-alpha.0 → v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #15 ⇒  cat VERSION
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #16 ⇒  bump -alpha -write
Bumped v1.0.0-alpha.0 → v1.0.0-alpha.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #17 ⇒  cat VERSION
v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #18 ⇒  bump -patch
Bumped v1.0.0-alpha.1 → v1.0.1
andrei@bump:test.sh ⚡ Test #19 ⇒  cat VERSION
v1.0.0-alpha.1
andrei@bump:test.sh ⚡ Test #20 ⇒  bump -patch -write
Bumped v1.0.0-alpha.1 → v1.0.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #21 ⇒  cat VERSION
v1.0.1
andrei@bump:test.sh ⚡ Test #22 ⇒  bump -major -write
Bumped v1.0.1 → v2.0.0 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #23 ⇒  cat VERSION
v2.0.0
andrei@bump:test.sh ⚡ Test #24 ⇒  bump -preview -write
Bumped v2.0.0 → v2.0.0-preview.1 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #25 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #26 ⇒  echo "--- SCENARIO THREE ---"
--- SCENARIO THREE ---
andrei@bump:test.sh ⚡ Test #27 ⇒  echo "1.25" > VERSION
andrei@bump:test.sh ⚡ Test #28 ⇒  bump -fix
Bumped v1.25.0 → v1.25.0
andrei@bump:test.sh ⚡ Test #29 ⇒  bump -fix -write
Parsed v1.25.0 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #30 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #31 ⇒  echo "--- SCENARIO FOUR ---"
--- SCENARIO FOUR ---
andrei@bump:test.sh ⚡ Test #32 ⇒  echo "v1.17.7-beta.6" > VERSION
andrei@bump:test.sh ⚡ Test #33 ⇒  bump -check -fix
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #34 ⇒  cat VERSION
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #35 ⇒  bump -check -fix -write
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #36 ⇒  cat VERSION
v1.17.7-beta.6
andrei@bump:test.sh ⚡ Test #37 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #38 ⇒  echo "--- SCENARIO FIVE ---"
--- SCENARIO FIVE ---
andrei@bump:test.sh ⚡ Test #39 ⇒  echo "module testApp-bump-test-passes" > go.mod
andrei@bump:test.sh ⚡ Test #40 ⇒  echo "" >> go.mod
andrei@bump:test.sh ⚡ Test #41 ⇒  echo "go 1.24" >> go.mod
andrei@bump:test.sh ⚡ Test #42 ⇒  cat go.mod
module testApp-bump-test-passes

go 1.24
andrei@bump:test.sh ⚡ Test #43 ⇒  bump -in go.mod -fix
Bumped 1.24 → 1.24
andrei@bump:test.sh ⚡ Test #44 ⇒  bump -in go.mod -fix -write
Parsed 1.24 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #45 ⇒  cat go.mod
module testApp-bump-test-passes

go 1.24.0
andrei@bump:test.sh ⚡ Test #46 ⇒  rm go.mod
andrei@bump:test.sh ⚡ Test #47 ⇒  echo "--- SCENARIO SIX ---"
--- SCENARIO SIX ---
andrei@bump:test.sh ⚡ Test #48 ⇒  echo v1.0.0 > VERSION
andrei@bump:test.sh ⚡ Test #49 ⇒  bump -json -check
{
  "version": "v1.0.0"
}
andrei@bump:test.sh ⚡ Test #50 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #51 ⇒  bump -json -beta
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
andrei@bump:test.sh ⚡ Test #52 ⇒  cat VERSION
v1.0.0
andrei@bump:test.sh ⚡ Test #53 ⇒  bump -json -beta -write
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
andrei@bump:test.sh ⚡ Test #54 ⇒  cat VERSION
v1.0.0-alpha.0
andrei@bump:test.sh ⚡ Test #55 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #56 ⇒  echo "--- SCENARIO SEVEN: package.json ---"
--- SCENARIO SEVEN: package.json ---
andrei@bump:test.sh ⚡ Test #57 ⇒  echo eyJuYW1lIjoidGVzdCIsInZlcnNpb24iOiIxLjIuMyJ9 | base64 -d | tee package.json > /dev/null
andrei@bump:test.sh ⚡ Test #58 ⇒  cat package.json
{"name":"test","version":"1.2.3"}
andrei@bump:test.sh ⚡ Test #59 ⇒  bump -in package.json -fix
Bumped 1.2.3 → 1.2.3
andrei@bump:test.sh ⚡ Test #60 ⇒  bump -in package.json -fix -write
Parsed 1.2.3 (saved to package.json)
andrei@bump:test.sh ⚡ Test #61 ⇒  bump -in package.json -patch
Bumped 1.2.3 → v1.2.4
andrei@bump:test.sh ⚡ Test #62 ⇒  bump -in package.json -patch -write
Bumped 1.2.3 → v1.2.4 (saved to package.json)
andrei@bump:test.sh ⚡ Test #63 ⇒  grep '"version": "1.2.4"' package.json
  "version": "1.2.4"
andrei@bump:test.sh ⚡ Test #64 ⇒  bump -in package.json -json -minor
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
andrei@bump:test.sh ⚡ Test #65 ⇒  bump -in package.json -minor -write
Bumped 1.2.4 → 1.3.0 (saved to package.json)
andrei@bump:test.sh ⚡ Test #66 ⇒  grep '"version": "1.3.0"' package.json
  "version": "1.3.0"
andrei@bump:test.sh ⚡ Test #67 ⇒  rm package.json
andrei@bump:test.sh ⚡ Test #68 ⇒  echo "--- SCENARIO EIGHT: go.mod ---"
--- SCENARIO EIGHT: go.mod ---
andrei@bump:test.sh ⚡ Test #69 ⇒  echo 'module myapp' > go.mod && echo 'go 1.21' >> go.mod
andrei@bump:test.sh ⚡ Test #70 ⇒  bump -in go.mod -check
1.21
andrei@bump:test.sh ⚡ Test #71 ⇒  bump -in go.mod -minor
Bumped 1.21 → 1.22
andrei@bump:test.sh ⚡ Test #72 ⇒  bump -in go.mod -minor -write
Bumped 1.21 → 1.22 (saved to go.mod)
andrei@bump:test.sh ⚡ Test #73 ⇒  grep 'go 1.22' go.mod
go 1.22.0
andrei@bump:test.sh ⚡ Test #74 ⇒  rm go.mod
andrei@bump:test.sh ⚡ Test #75 ⇒  echo "--- SCENARIO NINE: Dockerfile ---"
--- SCENARIO NINE: Dockerfile ---
andrei@bump:test.sh ⚡ Test #76 ⇒  echo 'LABEL version="v3.2.1"' > Dockerfile
andrei@bump:test.sh ⚡ Test #77 ⇒  bump -in Dockerfile -check
v3.2.1
andrei@bump:test.sh ⚡ Test #78 ⇒  bump -in Dockerfile -patch
Bumped v3.2.1 → v3.2.2
andrei@bump:test.sh ⚡ Test #79 ⇒  bump -in Dockerfile -patch -write
Bumped v3.2.1 → v3.2.2 (saved to Dockerfile)
andrei@bump:test.sh ⚡ Test #80 ⇒  grep 'LABEL version="v3.2.2"' Dockerfile
LABEL version="v3.2.2"
andrei@bump:test.sh ⚡ Test #81 ⇒  rm Dockerfile
andrei@bump:test.sh ⚡ Test #82 ⇒  echo "--- SCENARIO TEN: Chart.yaml (Helm) ---"
--- SCENARIO TEN: Chart.yaml (Helm) ---
andrei@bump:test.sh ⚡ Test #83 ⇒  echo 'YXBpVmVyc2lvbjogdjIKbmFtZTogbXljaGFydAp2ZXJzaW9uOiAwLjEuMAo=' | base64 -d | tee Chart.yaml > /dev/null
andrei@bump:test.sh ⚡ Test #84 ⇒  bump -in Chart.yaml -check
0.1.0
andrei@bump:test.sh ⚡ Test #85 ⇒  bump -in Chart.yaml -patch
Bumped 0.1.0 → v0.1.1
andrei@bump:test.sh ⚡ Test #86 ⇒  bump -in Chart.yaml -patch -write
Bumped 0.1.0 → v0.1.1 (saved to Chart.yaml)
andrei@bump:test.sh ⚡ Test #87 ⇒  grep 'version: 0.1.1' Chart.yaml
version: 0.1.1
andrei@bump:test.sh ⚡ Test #88 ⇒  rm Chart.yaml
andrei@bump:test.sh ⚡ Test #89 ⇒  echo "--- SCENARIO ELEVEN: pom.xml (Maven) ---"
--- SCENARIO ELEVEN: pom.xml (Maven) ---
andrei@bump:test.sh ⚡ Test #90 ⇒  echo '<project><version>2.2.2</version></project>' > pom.xml
andrei@bump:test.sh ⚡ Test #91 ⇒  bump -in pom.xml -check
2.2.2
andrei@bump:test.sh ⚡ Test #92 ⇒  bump -in pom.xml -patch
Bumped 2.2.2 → v2.2.3
andrei@bump:test.sh ⚡ Test #93 ⇒  bump -in pom.xml -patch -write
Bumped 2.2.2 → v2.2.3 (saved to pom.xml)
andrei@bump:test.sh ⚡ Test #94 ⇒  grep '<version>2.2.3</version>' pom.xml
<project><version>2.2.3</version></project>
andrei@bump:test.sh ⚡ Test #95 ⇒  rm pom.xml
andrei@bump:test.sh ⚡ Test #96 ⇒  echo "--- SCENARIO TWELVE: Environment Variables ---"
--- SCENARIO TWELVE: Environment Variables ---
andrei@bump:test.sh ⚡ Test #97 ⇒  echo "v5.5.5" > VERSION
andrei@bump:test.sh ⚡ Test #98 ⇒  BUMP_ALWAYS_WRITE=true bump -env
BUMP_NEVER_FIX=false
BUMP_NO_ALPHA=false
BUMP_NO_BETA=false
BUMP_NO_ALPHA_BETA=false
BUMP_ALWAYS_FIX=false
BUMP_ALWAYS_WRITE=true
BUMP_DEFAULT_INPUT=VERSION
BUMP_NO_RC=false
BUMP_NO_PREVIEW=false
BUMP_INIT_ON_NOT_FOUND=false
andrei@bump:test.sh ⚡ Test #99 ⇒  BUMP_ALWAYS_WRITE=true bump -patch
Bumped v5.5.5 → v5.5.6 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #100 ⇒  grep 'v5.5.6' VERSION
v5.5.6
andrei@bump:test.sh ⚡ Test #101 ⇒  BUMP_DEFAULT_INPUT=VERSION bump -minor
Bumped v5.5.6 → v5.6.0
andrei@bump:test.sh ⚡ Test #102 ⇒  BUMP_DEFAULT_INPUT=VERSION bump -minor -write
Bumped v5.5.6 → v5.6.0 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #103 ⇒  grep 'v5.6.0' VERSION
v5.6.0
andrei@bump:test.sh ⚡ Test #104 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #105 ⇒  echo "--- TEST THIRTEEN ---"
--- TEST THIRTEEN ---
andrei@bump:test.sh ⚡ Test #106 ⇒  bump -parse v1.2.3-alpha.4 -init
Initialized v1.2.3-alpha.4
andrei@bump:test.sh ⚡ Test #107 ⇒  cat VERSION
v1.2.3-alpha.4
andrei@bump:test.sh ⚡ Test #108 ⇒  bump -parse v2.3.4-alpha.5 -write
Parsed v2.3.4-alpha.5 (saved to VERSION)
andrei@bump:test.sh ⚡ Test #109 ⇒  cat VERSION
v2.3.4-alpha.5
andrei@bump:test.sh ⚡ Test #110 ⇒  rm VERSION
andrei@bump:test.sh ⚡ Test #111 ⇒  bump -parse v3.4.5-alpha.6 -init
Initialized v3.4.5-alpha.6
All 111 tests PASS!
```

