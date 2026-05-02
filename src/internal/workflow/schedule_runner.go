/*-------------------------------------------------------------------------
 *
 * schedule_runner.go
 *    Polls workflow_schedules and executes due workflows (cron).
 *
 *-------------------------------------------------------------------------
 */

package workflow

import (
	"context"
	"time"

	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/rs/zerolog/log"
)

/* ScheduleRunner fires workflows whose next_run_at is due. */
type ScheduleRunner struct {
	queries  *db.Queries
	engine   *Engine
	interval time.Duration
}

func NewScheduleRunner(queries *db.Queries, engine *Engine, interval time.Duration) *ScheduleRunner {
	if interval < 5*time.Second {
		interval = 5 * time.Second
	}
	return &ScheduleRunner{queries: queries, engine: engine, interval: interval}
}

/* Start runs the polling loop until ctx is cancelled. */
func (r *ScheduleRunner) Start(ctx context.Context) {
	go r.loop(ctx)
}

func (r *ScheduleRunner) loop(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	r.runTick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.runTick(ctx)
		}
	}
}

func (r *ScheduleRunner) runTick(ctx context.Context) {
	due, err := r.queries.ListWorkflowSchedulesByNextRun(ctx, time.Now())
	if err != nil {
		log.Warn().Err(err).Msg("workflow schedule tick: list due schedules failed")
		return
	}
	for i := range due {
		s := &due[i]
		next, err := NextScheduledRun(s.CronExpression, s.Timezone, time.Now())
		if err != nil {
			log.Warn().Err(err).Str("workflow_id", s.WorkflowID.String()).Msg("workflow schedule: skip invalid cron")
			continue
		}
		now := time.Now()
		s.LastRunAt = &now
		s.NextRunAt = &next
		if err := r.queries.UpdateWorkflowSchedule(ctx, s); err != nil {
			log.Warn().Err(err).Str("workflow_id", s.WorkflowID.String()).Msg("workflow schedule: advance next_run failed")
			continue
		}
		trigger := map[string]interface{}{
			"schedule_id":     s.ID.String(),
			"cron_expression": s.CronExpression,
			"timezone":        s.Timezone,
		}
		_, execErr := r.engine.ExecuteWorkflow(ctx, s.WorkflowID, "schedule", trigger, nil)
		if execErr != nil {
			log.Warn().Err(execErr).Str("workflow_id", s.WorkflowID.String()).Msg("workflow schedule execution failed")
		}
	}
}
