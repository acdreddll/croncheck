// Package parser provides utilities for parsing and validating standard
// 5-field cron expressions of the form:
//
//	<minute> <hour> <day-of-month> <month> <day-of-week>
//
// Each field may contain:
//   - A wildcard (*)
//   - A single integer value
//   - A comma-separated list of values (e.g. 1,3,5)
//   - A range (e.g. 1-5)
//   - A step expression (e.g. */5 or 1-10/2)
//
// Use [Parse] to convert a raw string into a [CronExpression] struct.
// Validation is performed automatically during parsing.
package parser
