# Horizon OpenAPI 3.0

This repo contains definitions of Horizon API in OpenAPI 3.0.

## Contribution

As ReDoc has some issues with rendering names of the referenced types, all files are combined into one.
Run `make` to get `public/index.html` (requires global installation of `@tokend/redoc-cli`).
Reference all the models as if they were in one file. Example: `$ref: '#/components/schemas/Resources'`.
