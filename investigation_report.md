# Investigation Report: Bash Job Number Recycling in PR #294

## Hypothesis Being Tested
The hypothesis was that the bash version in the CI environment (Alpine Linux) is too low, causing the fixture regeneration failure for the job number recycling test.

## Key Findings

### 1. Bash Versions Involved

- **Alpine Linux (from CI logs):** bash 5.3.3-r1
- **Current test environment:** bash 5.2.21
- **macOS (from Slack):** bash 3.2.57

### 2. The Hypothesis is INCORRECT

The bash version in Alpine Linux (5.3.3) is actually quite modern and recent. It is NOT too low.

### 3. Actual Problem: Test Expectations vs Bash Behavior

The test in `testBg9Recycle` expects bash to immediately recycle job number 2 after it's reaped, but **bash does not work this way**.

#### Test Scenario:
1. Start job [1]: `sleep 100 &`
2. Start job [2]: `cat /tmp/fifo &`
3. Start job [3]: `sleep 50 &`
4. Complete job 2 by writing to the FIFO
5. Run `echo apple` (which displays the reaped job 2)
6. Start new job: `sleep 10 &`
   - **Test expects:** Job number [2]
   - **Bash assigns:** Job number [4]

#### Bash's Actual Behavior:

Bash maintains terminated jobs in the jobs table with "Done" status. These jobs are only completely removed (and their numbers made available for recycling) after:
- Running the `jobs` builtin (which removes them after display)
- Or in POSIX mode, after certain other operations

**Evidence from testing:**
```
[1]   Running                 sleep 100 &
[2]   Done                    cat "$FIFO"     ← Still in table
[3]-  Running                 sleep 50 &
[4]+  Running                 sleep 10 &      ← Got [4], not [2]
```

But if we run `jobs` to completely remove the terminated job first:
```
$ jobs  # Shows job 2 as Done
$ jobs  # Removes job 2 completely
$ sleep 10 &
[2]+  Running                 sleep 10 &      ← NOW it gets [2]
```

### 4. Consistency Across Versions

Both bash 5.3.3 (Alpine) and bash 5.2.21 (test environment) exhibit the SAME behavior - neither recycles job numbers immediately after reaping. This suggests the behavior is consistent across modern bash versions.

### 5. Why Did macOS bash 3.2 Also Not Recycle?

This actually supports the finding that job number recycling doesn't work the way the test expects, even in older versions. It's not a version-specific bug but a misunderstanding of bash's job management behavior.

## Conclusion

**The hypothesis that bash version is too low is FALSE.**

The real issue is that the test has incorrect expectations about when bash recycles job numbers. Bash does recycle job numbers, but only after terminated jobs are fully removed from the jobs table, not immediately after they are reaped and displayed.

## Technical Details

### CI Environment
- **Platform:** GitHub Actions, ubuntu-latest (Ubuntu 24.04)
- **Docker Image:** `golang:1.24-alpine`
- **Bash installed via:** `apk add bash` (Alpine package)
- **Bash version:** 5.3.3-r1 (modern, recent release from Aug 2025)

### Test Failure Location
- **File:** `internal/stage_bg9.go`
- **Function:** `testBg9Recycle`
- **Line:** ~186 (expectation of job number 2 for recycled job)

### Relevant CI Logs
From the failed run at https://github.com/codecrafters-io/shell-tester/actions/runs/27519938967:
```
[your-program] $ sleep 10 &
[your-program] [4] 1736
[tester::#FY4] ^ Line does not match expected value.
[tester::#FY4] Expected: "[2] <PID>"
[tester::#FY4] Received: "[4] 1736"
```

## Recommendation

The test in PR #294 (testBg9Recycle) needs to be revised to match bash's actual behavior:

### Option 1: Modify test to explicitly reap the job
Add a `jobs` command after reaping job 2 to fully remove it from the table:
```go
// After echoing to reap job 2, add:
jobsTestCase := test_cases.JobsBuiltinResponseTestCase{
    ExpectedOutputEntries: []test_cases.BackgroundJobStatusEntry{
        {JobNumber: 1, ...},
        {JobNumber: 3, ...},
    },
    SuccessMessage: "✓ Job 2 fully reaped",
}
if err := jobsTestCase.Run(asserter, shell, logger); err != nil {
    return err
}
// NOW start sleep 10 & and expect job [2]
```

### Option 2: Change expected job number
Simply change the expectation from job [2] to job [4].

### Option 3: Remove this specific test
Reconsider whether testing immediate gap-filling is necessary, since it doesn't match bash's documented or actual behavior.

## References

- Alpine Linux bash package: https://pkgs.alpinelinux.org/package/v3.23/main/x86_64/bash
- Bash 5.3 release notes: https://tiswww.case.edu/php/chet/bash/NEWS
- PR under investigation: https://github.com/codecrafters-io/shell-tester/pull/294
- CI failure: https://github.com/codecrafters-io/shell-tester/actions/runs/27519938967/job/81335755853

## Test Artifacts

Test scripts demonstrating the behavior are available at:
- `/tmp/test_job_recycling.sh` - Initial reproduction
- `/tmp/test_simple_recycle.sh` - Demonstrates recycling after full removal  
- `/tmp/test_exact_scenario.sh` - Exact reproduction of test scenario
