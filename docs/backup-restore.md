# Backup and restore

NeuronAgent state lives in **PostgreSQL** — backup the database like any production Postgres workload.

## Logical backup

```bash
pg_dump -Fc -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f neuronagent.dump
```

## Restore

```bash
pg_restore -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" --clean neuronagent.dump
```

Use your organization’s retention, encryption, and PITR policies. Point-in-time recovery follows PostgreSQL/WAL configuration on your infrastructure.
