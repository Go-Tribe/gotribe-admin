# Template Bootstrap

Use this project as a template by running the bootstrap script once after copying the repository:

```bash
./scripts/init-project.sh \
  --app-name xxxxx-admin \
  --module-path github.com/your-org/xxxxx-admin
```

What it updates:

- `go.mod` module path
- Go import paths
- entry file name
- build and deployment files such as `Makefile`, `Dockerfile`, `docker-compose.yml`
- README / Swagger / release metadata that still reference the template name

The script is only for project initialization. The application runtime does not depend on it.
