# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
--- FAIL: TestInvalidDeviceSerialDoesNotIssueCertificate (0.02s)
    bug_test.go:27: invalid serial unexpectedly succeeded
FAIL
FAIL	devicecert	0.020s
ok  	devicecert/cli	0.019s
ok  	devicecert/cmd/devicecert	0.002s
ok  	devicecert/cryptoengine	0.004s
ok  	devicecert/domain	0.004s
ok  	devicecert/safelog	0.003s
ok  	devicecert/service	0.026s
ok  	devicecert/store	0.013s
ok  	devicecert/testkit	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/devicecert): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/devicecert): exit `0`
