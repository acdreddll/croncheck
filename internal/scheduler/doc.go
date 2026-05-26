// Package scheduler provides utilities for computing scheduled execution
// times from parsed cron expressions.
//
// Given an [audit.Entry], NextN returns the next n wall-clock times at
// which the job would fire.  Overlaps checks whether two entries share
// any scheduled minute within a sliding time window, which can be used
// as a richer conflict-detection strategy than the field-overlap
// heuristic in the audit package.
package scheduler
