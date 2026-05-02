/*-------------------------------------------------------------------------
 *
 * cron_schedule.go
 *    Workflow cron next-run computation (robfig/cron)
 *
 *-------------------------------------------------------------------------
 */

package workflow

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

/* NextScheduledRun returns the next instant matching cronExpr after `after`, interpreted in tz (IANA name, e.g. UTC, America/New_York). */
func NextScheduledRun(cronExpr, tz string, after time.Time) (time.Time, error) {
	if cronExpr == "" {
		return time.Time{}, fmt.Errorf("cron expression is empty")
	}
	loc := time.UTC
	if tz != "" && tz != "UTC" {
		l, err := time.LoadLocation(tz)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid timezone %q: %w", tz, err)
		}
		loc = l
	}
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse cron %q: %w", cronExpr, err)
	}
	t := after.In(loc)
	return schedule.Next(t), nil
}
