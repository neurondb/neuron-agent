# Security

## Reporting a vulnerability

Please report security issues privately to the maintainers (see repository contact or organization security policy). Do not open public issues for undisclosed vulnerabilities.

Include:

- Description and impact
- Steps to reproduce
- Affected versions or commits if known

## Supported versions

Security fixes are applied to the default branch (`main`). Tagged releases will note security-related changes in the changelog.

## Deployment hard practices

- Use strong database passwords; production mode rejects common default passwords when `CONFIG_PROFILE` / `ENV` indicates production.
- Protect API keys; never commit `.env` files.
- Run NeuronAgent behind TLS in production and restrict network access to PostgreSQL.
