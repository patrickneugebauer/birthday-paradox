# Pharo Language Research

## Overview
Pharo is a modern Smalltalk implementation with a live programming environment. For CLI scripts, Pharo can be run in headless mode where command-line arguments are accessible through system attributes.
- Source: [Pharo Official Site](https://pharo.org/)

## Command-Line Arguments
- `Smalltalk getSystemAttribute: n` retrieves the nth system attribute/command-line argument (1-indexed)
- `ifNil:` checks if the result is nil (missing argument)
- In headless mode: attribute 1 is the first command-line argument, attribute 0 is the image path
- Source: [Pharo Manual - System](https://github.com/pharo-open-documentation/pharo-wiki/blob/master/General/CommandLine.md)

## Argument Parsing
- `asNumber` converts a String to a Number (raises Error if invalid)
- `[block] on: Error do: [:err | handler]` is the exception handling syntax
- Any number parsing error is caught and handled with descriptive output
- Source: [Pharo Manual - Exception Handling](https://github.com/pharo-open-documentation/pharo-wiki)

## Error Output
- `Stdio stderr` is the standard error stream
- `nextPutAll: 'string'` writes a string to the stream
- `lf` appends a line feed (newline)
- Source: [Pharo Class - Stdio](https://github.com/pharo-open-documentation/pharo-wiki)

## Program Termination
- `Smalltalk snapshot: false andQuit: true` exits the Pharo image
- `snapshot: false` prevents saving the image state
- `andQuit: true` terminates the process
- Source: [Pharo Manual - Image Management](https://github.com/pharo-open-documentation/pharo-wiki)

## Timing
- `DateAndTime now` returns the current date and time
- `.floor` truncates to whole seconds for simpler calculation
- Note: More precise timing would require `System getTickCount` or similar
- Source: [Pharo Class - DateAndTime](https://github.com/pharo-open-documentation/pharo-wiki)

## Rounding
- `rounded` rounds a number to the nearest integer
- `(value * 10^n) rounded / 10^n` rounds to n decimal places
- For 2 decimal places: `(percent * 100) rounded / 100.0`
- For 6 decimal places: `(seconds * 1000000) rounded / 1000000.0`
- Source: [Pharo Manual - Number Methods](https://github.com/pharo-open-documentation/pharo-wiki)

## Number Formatting for Output
- `printShowingDecimalPlaces: n` formats a number with exactly n decimal places
- `value asString` converts to basic string representation
- String concatenation: `, ` operator joins strings
- `String lf` is a newline character constant
- Source: [Pharo Class - Number](https://github.com/pharo-open-documentation/pharo-wiki)

## Type Conversions
- `asFloat` converts to floating-point number
- `asNumber` converts string to appropriate numeric type
- `asString` converts any object to string representation
- Source: [Pharo Manual - Object Protocol](https://github.com/pharo-open-documentation/pharo-wiki)

## Variable Declaration
- Local variables are declared in `| var1 var2 |` syntax
- `:=` is the assignment operator
- Smalltalk is dynamically typed; no explicit type declarations needed
- Source: [Pharo Manual - Language Syntax](https://github.com/pharo-open-documentation/pharo-wiki)
