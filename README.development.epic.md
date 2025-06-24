# EPIC
## User Story: SQL-Backed REST API for Product and Category Management

### Description

As a developer, I want to create a **SQL-backed REST API** that allows clients to perform **CRUD operations** on **products** and **categories**, so that product-related data can be managed efficiently through a standardized interface.

This API will:

- Provide endpoints for managing products and categories
- Persist data in a SQL database (e.g., PostgreSQL)
- Return structured and meaningful error responses
- Log all errors at the point they occur with contextual metadata

---

### Acceptance Criteria

- [ ] API exposes the following endpoints:
  - `POST /products`
  - `GET /products`
  - `GET /products/{id}`
  - `PUT /products/{id}`
  - `DELETE /products/{id}`
  - `POST /categories`
  - `GET /categories`
  - `GET /categories/{id}`
  - `PUT /categories/{id}`
  - `DELETE /categories/{id}`

- [ ] Product entity includes the following fields:
  - `id`, `name`, `description`, `price`, `category_id`, `created_at`, `updated_at`

- [ ] Category entity includes the following fields:
  - `id`, `name`, `description`, `created_at`, `updated_at`

- [ ] API is backed by a SQL database with properly indexed tables for performance

- [ ] Validation errors and server errors are logged with:
  - Operation name
  - Request ID (if available)
  - Error message and metadata

- [ ] API returns appropriate HTTP status codes and error messages in a consistent format

- [ ] Unit and integration tests cover all endpoints and edge cases

- [ ] OpenAPI (Swagger) documentation is available and up-to-date

---

### Technical Notes

- **Language/Framework:** Go with Gorilla Mux or similar
- **Logging:** Use structured logging (e.g., zerolog) with contextual metadata
- **Middleware:** 
  - Request ID injection
  - Centralized error handling and logging
- **Database:** PostgreSQL (or any SQL-compliant DB)
- **Migrations:** Use version-controlled migration tools (e.g., `golang-migrate`)

---

### Tasks

- [ ] Design SQL schema for products and categories
- [ ] Implement CRUD endpoints for products
- [ ] Implement CRUD endpoints for categories
- [ ] Add centralized error logging and request tracing
- [ ] Write unit and integration tests
- [ ] Generate and expose API documentation
