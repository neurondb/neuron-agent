/*-------------------------------------------------------------------------
 *
 * profiles.go
 *    Compliance profiles: standard, gdpr, hipaa, sox.
 *
 * Each profile toggles memory retention, export restrictions, audit verbosity,
 * and data anonymization. Wire into memory promotion TTL, audit logging, export endpoints.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/compliance/profiles.go
 *
 *-------------------------------------------------------------------------
 */

package compliance

import "time"

/* ProfileName is the compliance profile identifier. */
const (
	ProfileStandard = "standard"
	ProfileGDPR     = "gdpr"
	ProfileHIPAA    = "hipaa"
	ProfileSOX      = "sox"
)

/* ProfileSettings holds per-profile compliance settings. */
type ProfileSettings struct {
	MemoryRetentionDays       int           /* Max days to retain memory chunks; 0 = no automatic limit */
	ExportRestricted          bool          /* If true, restrict bulk export endpoints */
	AuditVerbose              bool          /* If true, log full request/response in audit where applicable */
	AnonymizePII              bool          /* If true, redact PII in audit and exports */
	EncryptionRequired        bool          /* If true, require TLS and encryption at rest */
	ApprovalRequiredWorkflows bool          /* If true, all workflow executions require approval gate */
	DefaultRetention          time.Duration /* Default TTL for ephemeral data */
}

/* Profile returns settings for the given profile name. */
func Profile(name string) ProfileSettings {
	switch name {
	case ProfileGDPR:
		return ProfileSettings{
			MemoryRetentionDays:       90,
			ExportRestricted:          true,
			AuditVerbose:              true,
			AnonymizePII:              true,
			EncryptionRequired:        false,
			ApprovalRequiredWorkflows: false,
			DefaultRetention:          90 * 24 * time.Hour,
		}
	case ProfileHIPAA:
		return ProfileSettings{
			MemoryRetentionDays:       365,
			ExportRestricted:          true,
			AuditVerbose:              true,
			AnonymizePII:             false, /* HIPAA may require specific handling; configurable */
			EncryptionRequired:        true,
			ApprovalRequiredWorkflows: false,
			DefaultRetention:          365 * 24 * time.Hour,
		}
	case ProfileSOX:
		return ProfileSettings{
			MemoryRetentionDays:       2555, /* 7 years in days */
			ExportRestricted:          true,
			AuditVerbose:              true,
			AnonymizePII:              false,
			EncryptionRequired:        true,
			ApprovalRequiredWorkflows: true,
			DefaultRetention:          2555 * 24 * time.Hour,
		}
	case ProfileStandard:
		fallthrough
	default:
		return ProfileSettings{
			MemoryRetentionDays:       0,
			ExportRestricted:          false,
			AuditVerbose:              false,
			AnonymizePII:              false,
			EncryptionRequired:        false,
			ApprovalRequiredWorkflows: false,
			DefaultRetention:          0,
		}
	}
}

/* MemoryRetentionDays returns the memory retention in days for the profile (0 = no limit). */
func (p ProfileSettings) MemoryRetentionDaysValue() int {
	return p.MemoryRetentionDays
}

/* IsExportRestricted returns whether bulk export is restricted. */
func (p ProfileSettings) IsExportRestricted() bool {
	return p.ExportRestricted
}

/* IsAuditVerbose returns whether to log verbose audit (e.g. full payloads). */
func (p ProfileSettings) IsAuditVerbose() bool {
	return p.AuditVerbose
}

/* IsAnonymizePII returns whether to anonymize PII in audit/exports. */
func (p ProfileSettings) IsAnonymizePII() bool {
	return p.AnonymizePII
}

/* IsEncryptionRequired returns whether TLS/encryption is required. */
func (p ProfileSettings) IsEncryptionRequired() bool {
	return p.EncryptionRequired
}

/* IsApprovalRequiredWorkflows returns whether all workflows require approval. */
func (p ProfileSettings) IsApprovalRequiredWorkflows() bool {
	return p.ApprovalRequiredWorkflows
}
