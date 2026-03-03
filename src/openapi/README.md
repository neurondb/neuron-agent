# NeuronAgent OpenAPI Specification

This directory contains the OpenAPI (formerly Swagger) specification for the NeuronAgent REST API.

## Files

- `openapi.yaml` - OpenAPI 3.0.3 specification file

## Versioning

The OpenAPI specification is versioned with each release. Snapshot files are named according to the NeuronAgent version:

- `v1.0.0.yaml` - API snapshot for version 1.0.0
- `v1.1.0.yaml` - API snapshot for version 1.1.0
- etc.

The main `openapi.yaml` file always represents the current development version.

## Usage

### Viewing the API Documentation

You can view the API documentation using various tools:

#### Swagger UI

```bash
# Using Docker
docker run -p 8080:8080 -e SWAGGER_JSON=/openapi.yaml -v $(pwd):/docs swaggerapi/swagger-ui

# Then open http://localhost:8080 in your browser
```

#### Redoc

```bash
# Install Redoc CLI
npm install -g redoc-cli

# Generate HTML documentation
redoc-cli serve openapi.yaml

# Or build static HTML
redoc-cli build openapi.yaml -o api-docs.html
```

#### Online Tools

- [Swagger Editor](https://editor.swagger.io/) - Paste the YAML content to view and edit
- [Redocly](https://redocly.github.io/redoc/) - Online API documentation generator

### Validating the Specification

```bash
# Install swagger-codegen or openapi-generator
npm install -g @apidevtools/swagger-cli

# Validate the specification
swagger-cli validate openapi.yaml
```

### Generating Client Libraries

You can generate client libraries from the OpenAPI specification:

```bash
# Using openapi-generator
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
  -i /local/openapi.yaml \
  -g python \
  -o /local/clients/python

# Using swagger-codegen
swagger-codegen generate -i openapi.yaml -l python -o clients/python
```

### Testing with Generated Clients

Example using generated Python client:

```python
import openapi_client
from openapi_client.api import agents_api
from openapi_client.model.create_agent_request import CreateAgentRequest

configuration = openapi_client.Configuration(
    host="http://localhost:8080/api/v1",
    api_key={'BearerAuth': 'your-api-key-here'}
)

with openapi_client.ApiClient(configuration) as api_client:
    api_instance = agents_api.AgentsApi(api_client)
    
    agent_request = CreateAgentRequest(
        name="my-agent",
        system_prompt="You are a helpful assistant",
        model_name="gpt-4",
        enabled_tools=["sql", "http"],
        config={"temperature": 0.7}
    )
    
    api_response = api_instance.create_agent(create_agent_request=agent_request)
    print(api_response)
```

## Specification Maintenance

### Updating the Specification

1. Update `openapi.yaml` with API changes
2. Validate the specification: `swagger-cli validate openapi.yaml`
3. Test with Swagger UI or Redoc
4. Commit changes

### Release Process

When creating a new release:

1. Copy `openapi.yaml` to `v{X.Y.Z}.yaml` (e.g., `v1.0.0.yaml`)
2. Update the version in the snapshot file's `info.version` field
3. Commit the snapshot file
4. Update `openapi.yaml` for the next development version

## Specification Coverage

The OpenAPI specification covers:

- ✅ All REST API endpoints
- ✅ Request/response schemas
- ✅ Authentication (Bearer token)
- ✅ Error responses
- ✅ Query parameters

**Note:** WebSocket endpoints (`/ws`) are not included in the OpenAPI spec as they use a different protocol. Refer to the API documentation for WebSocket usage.

## Related Documentation

- [API Documentation](../docs/API.md) - Human-readable API documentation
- [Architecture Guide](../docs/architecture.md) - System architecture
- [Deployment Guide](../docs/deployment.md) - Deployment instructions

