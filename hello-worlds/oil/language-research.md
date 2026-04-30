# Oil Shell Language Research

## Overview
Oil is a modern shell language that aims to improve upon Bash. It's compatible with Bash but adds cleaner syntax and better error handling.
- Source: [Oil Shell Official Site](https://www.oilshell.org/)

## Enabling Oil Features
- `shopt -s oil:all` enables all Oil language features in a script
- This must appear before Oil-specific syntax is used
- Source: [Oil Language Manual - Dialect](https://www.oilshell.org/release/latest/doc/language-manual.html)

## Command-Line Arguments
- `${#@}` returns the number of positional arguments (same as Bash `$#`)
- `$1` accesses the first argument as a string
- Check argument count with `(( ${#@} < 1 ))` - double parentheses for arithmetic evaluation
- Source: [Oil Shell Manual - Parameters](https://www.oilshell.org/release/latest/doc/oil-language-builtins.html)

## Argument Validation
- `[[ $var =~ ^[0-9]+$ ]]` uses regex to validate string is numeric
- `!` before `[[` negates the condition (if NOT matching pattern)
- This is Bash-compatible syntax, works in Oil with `oil:all` enabled
- Source: [Bash Manual - Conditional Constructs](https://www.gnu.org/software/bash/manual/html_node/Conditional-Constructs.html)

## Error Output
- `>&2` redirects stdout to stderr (2 is file descriptor for stderr)
- `echo 'message' >&2` prints error message to stderr
- Source: [Bash Manual - Redirections](https://www.gnu.org/software/bash/manual/html_node/Redirections.html)

## Exit Status
- `exit 1` terminates script with exit code 1 (non-zero indicates error)
- Source: [Bash Manual - Exit Status](https://www.gnu.org/software/bash/manual/html_node/Exit-Status.html)

## Timing
- `date +%s%N` outputs seconds (%s) and nanoseconds (%N) since Unix epoch
- Arithmetic: `$(($(date +%s%N)/1000))` converts nanoseconds to microseconds
- Source: [GNU Coreutils - date manual](https://www.gnu.org/software/coreutils/manual/html_node/date-invocation.html)

## Variable Assignment and Arithmetic
- `set variable = value` is Oil syntax for assignment
- `set count = $((count + 1))` performs arithmetic with `$(( ))`
- Both traditional `$1` and Oil `set` syntax can be mixed in the same file
- Source: [Oil Shell Manual - Variable Assignment](https://www.oilshell.org/release/latest/doc/oil-language-builtins.html)

## String Formatting
- `printf '%06d'` formats a number with leading zeros (6 digits total)
- Used for nanosecond padding in timestamp formatting
- Source: [GNU Coreutils - printf manual](https://www.gnu.org/software/coreutils/manual/html_node/printf-invocation.html)
