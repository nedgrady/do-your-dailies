---
name: create-endpoint
description: "Use when adding or updating a CRUD API resource in this repo (OpenAPI-first, dedicated store per entity, and consistent 404/422/500 semantics)."
---

# CRUD Endpoint

Use this workflow when implementing or extending a CRUD API resource.

## Purpose

- Add CRUD endpoints consistently with this codebase's patterns.
- Keep entity boundaries explicit (separate resource means separate store contract).
- Enforce observable behavior tests.

## Required Policy

1. For a new entity/resource, create a dedicated store interface and implementation (do not append to unrelated stores).
2. Start with OpenAPI contract updates before handler/store wiring.

## Resource Boundary Rules

1. If the request says the feature is its own resource/entity, treat it as independent in API routes, handler methods, and store contracts.
2. Do not nest behavior into an existing resource store unless the user explicitly asks for that coupling.
3. Keep route naming and operation IDs resource-specific.

## CRUD Contract Rules

1. Add or update OpenAPI paths and request/response schemas first.
2. Define operation IDs for all needed actions: list/get/create/update/delete.
3. Use explicit request models for mutation operations (for example `Create<Resource>Request`, `Update<Resource>Request`).
4. Use generated OpenAPI types in handlers.
5. Keep JSON casing and strict decode behavior aligned with existing API conventions.

## Date Query Rules

1. When an endpoint accepts a day-based filter, use ISO 8601 date-only format (`YYYY-MM-DD`) in query parameters.
2. Do not use time portions in date filters unless the user explicitly requests timestamp filtering.
3. For list endpoints with optional date filters, default to today when omitted.
4. Treat invalid date formats as bad request errors.

## Error Semantics

1. Invalid JSON or unknown fields: return 422.
2. Missing entity (resource ID not found) or missing related parent entity: return 404.
3. Successful delete with no body: return 204.
4. Unexpected persistence/runtime failures: return 500.

## Default Implementation Sequence

1. Update OpenAPI spec for resource paths and operation schemas.
2. Regenerate OpenAPI types if generation tooling is available.
3. Add or update domain request model(s) only as needed.
4. Add or extend dedicated store interface and concrete implementation for the entity.
5. Wire store dependency into application composition.
6. Implement handlers for list/get/create/update/delete using generated request types and API error mapping.
7. Add route registration for collection and entity ID paths.
8. Confirm behavior and error semantics match the contract.

## Test Expectations for CRUD Endpoints

Include tests for at least:

1. List returns 200 and expected collection payload.
2. Get returns 200 for existing ID and 404 for missing ID.
3. Create returns 201 and expected created payload.
4. Update returns 200 for existing ID and 404 for missing ID.
5. Delete returns 204 for existing ID and 404 for missing ID.
6. 422 on malformed JSON and unknown fields for create/update.
7. 404 when required parent entity does not exist.
8. 500 on store/internal error paths.
