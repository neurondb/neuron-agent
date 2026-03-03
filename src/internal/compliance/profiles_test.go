/*-------------------------------------------------------------------------
 *
 * profiles_test.go
 *    Tests for compliance profile settings.
 *
 *-------------------------------------------------------------------------
 */

package compliance

import (
	"testing"
	"time"
)

func TestProfile_Standard(t *testing.T) {
	p := Profile(ProfileStandard)
	if p.MemoryRetentionDaysValue() != 0 {
		t.Errorf("standard MemoryRetentionDays = %d", p.MemoryRetentionDaysValue())
	}
	if p.IsExportRestricted() {
		t.Error("standard should not restrict export")
	}
	if p.IsEncryptionRequired() {
		t.Error("standard should not require encryption")
	}
}

func TestProfile_GDPR(t *testing.T) {
	p := Profile(ProfileGDPR)
	if p.MemoryRetentionDaysValue() != 90 {
		t.Errorf("gdpr MemoryRetentionDays = %d", p.MemoryRetentionDaysValue())
	}
	if !p.IsExportRestricted() {
		t.Error("gdpr should restrict export")
	}
	if !p.IsAuditVerbose() {
		t.Error("gdpr should have verbose audit")
	}
	if !p.IsAnonymizePII() {
		t.Error("gdpr should anonymize PII")
	}
}

func TestProfile_HIPAA(t *testing.T) {
	p := Profile(ProfileHIPAA)
	if p.MemoryRetentionDaysValue() != 365 {
		t.Errorf("hipaa MemoryRetentionDays = %d", p.MemoryRetentionDaysValue())
	}
	if !p.IsEncryptionRequired() {
		t.Error("hipaa should require encryption")
	}
}

func TestProfile_SOX(t *testing.T) {
	p := Profile(ProfileSOX)
	if p.MemoryRetentionDaysValue() != 2555 {
		t.Errorf("sox MemoryRetentionDays = %d", p.MemoryRetentionDaysValue())
	}
	if !p.IsApprovalRequiredWorkflows() {
		t.Error("sox should require approval for workflows")
	}
}

func TestProfile_Unknown(t *testing.T) {
	p := Profile("unknown")
	if p.MemoryRetentionDaysValue() != 0 {
		t.Errorf("unknown profile should have 0 retention, got %d", p.MemoryRetentionDaysValue())
	}
}

func TestProfile_DefaultRetention(t *testing.T) {
	p := Profile(ProfileGDPR)
	if p.DefaultRetention != 90*24*time.Hour {
		t.Errorf("gdpr DefaultRetention = %v", p.DefaultRetention)
	}
}
