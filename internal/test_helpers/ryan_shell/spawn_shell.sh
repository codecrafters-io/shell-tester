#!/bin/sh
go build -o /Users/ryang/Developer/work/course-testers/shell-tester/internal/test_helpers/ryan_shell/ryan_shell /Users/ryang/Developer/work/course-testers/shell-tester/internal/test_helpers/ryan_shell/main.go
echo "Welcome to Ryan's shell!"
exec /Users/ryang/Developer/work/course-testers/shell-tester/internal/test_helpers/ryan_shell/ryan_shell
