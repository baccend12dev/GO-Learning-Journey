# Walkthrough: Implementation of Feature Requests, Documentation, Swagger & User Authentication

This document summarizes all changes made to the Application Knowledge Management System, covering testing setup, feature request development, routing corrections, documentation capabilities, Swagger UI integration, and User Authentication/Role Otorisasi.

---

## 1. Unit Testing Setup for Server Controller

To verify controllers without requiring a running PostgreSQL instance, we:
- Installed the cgo-free pure-Go SQLite GORM driver (`github.com/glebarez/sqlite`) to support compilation and execution in environments with `CGO_ENABLED=0`.
- Created [server_controller_test.go](file:///f:/GO/it-dashboard/backend/controllers/server_controller_test.go) containing comprehensive unit tests for all CRUD actions of `ServerController`.
- Verified successfully:
  ```
  ok  	backend/controllers	0.816s
  ```

---

## 2. Feature Requests Functionality

We implemented the Feature Requests feature allowing developers to record and track desired system features.

### Changes Made:
- **Model Definition** in [feature_request.go](file:///f:/GO/it-dashboard/backend/models/feature_request.go).
- **System Model Linking** in [system.go](file:///f:/GO/it-dashboard/backend/models/system.go) (one-to-many cascading relationship).
- **Database Migration** in [main.go](file:///f:/GO/it-dashboard/backend/main.go) (`config.DB.AutoMigrate(&models.FeatureRequest{})`).
- **Controller Actions** in [feature_request_controller.go](file:///f:/GO/it-dashboard/backend/controllers/feature_request_controller.go).
- **API Routing** in [feature_request_routes.go](file:///f:/GO/it-dashboard/backend/routes/feature_request_routes.go).

---

## 3. Documentation Module

We implemented the Documentation module to allow developers to create, read, update, and delete documentation writeups linked to specific systems.

### Changes Made:
1. **Model Definition** in [documentation.go](file:///f:/GO/it-dashboard/backend/models/documentation.go).
2. **System Model Linking** in [system.go](file:///f:/GO/it-dashboard/backend/models/system.go) (one-to-many cascading relationship).
3. **Database Migration** in [main.go](file:///f:/GO/it-dashboard/backend/main.go) (`config.DB.AutoMigrate(&models.Documentation{})`).
4. **Controller Actions** in [documentation_controller.go](file:///f:/GO/it-dashboard/backend/controllers/documentation_controller.go).
5. **API Routing** in [documentation_routes.go](file:///f:/GO/it-dashboard/backend/routes/documentation_routes.go).

---

## 4. Swagger API Documentation Setup

We integrated automatic Swagger API documentation using `swag` annotations and Gin middleware.

### Changes Made:
1. **Global Swagger Configuration** in [main.go](file:///f:/GO/it-dashboard/backend/main.go).
2. **Controller Handler Annotations**: Added comments above all 23 handler functions in all controller files.
3. **Generated Documentation Assets**: Generated files under `backend/docs/`.

---

## 5. User Authentication & Role-Based Authorization

We added a secure JWT user authentication and role-based route authorization system.

### Changes Made:
1. **Model Definition** in [user.go](file:///f:/GO/it-dashboard/backend/models/user.go):
   - Created `User` GORM model (`ID`, `Username`, `Password` (hashed), `Role`, `CreatedAt`, `UpdatedAt`).
   - Setup `RegisterRequest`, `LoginRequest`, and `LoginResponse` structs.
2. **Password Cryptography** in [bcrypt.go](file:///f:/GO/it-dashboard/backend/utils/bcrypt.go):
   - Used `golang.org/x/crypto/bcrypt` to hash passwords upon registration and verify passwords on login.
3. **JWT Token Utilities** in [jwt.go](file:///f:/GO/it-dashboard/backend/utils/jwt.go):
   - Configured `golang-jwt/jwt/v5` to sign and parse tokens containing User ID, Username, and Role claims.
4. **Auth & Authorization Middlewares** in [auth.go](file:///f:/GO/it-dashboard/backend/middleware/auth.go):
   - `AuthMiddleware()`: Reads the `Authorization: Bearer <token>` header, validates it, and saves claims in the Gin context.
   - `RoleMiddleware(allowedRoles ...string)`: Restricts routes to specific roles, returning `403 Forbidden` if disallowed.
5. **Auth Controller** in [auth_controller.go](file:///f:/GO/it-dashboard/backend/controllers/auth_controller.go) and **Auth Routes** in [auth_routes.go](file:///f:/GO/it-dashboard/backend/routes/auth_routes.go):
    - `POST /api/auth/register`: Public endpoint to register new users (validates role choices).
    - `POST /api/auth/login`: Public endpoint to login and receive a signed JWT token.
    - `GET /api/auth/me`: Secured endpoint to retrieve the current user's profile details.
    - `POST /api/auth/logout`: Secured endpoint to logout the current user.
6. **Admin Seeding** in [main.go](file:///f:/GO/it-dashboard/backend/main.go):
   - Added auto-migration `config.DB.AutoMigrate(&models.User{})`.
   - Enabled startup seeding: If no users exist, automatically creates an Administrator account with credentials `admin` / `admin123`.
7. **Secured Routing Groups**:
   - Wrapped GET endpoints (Systems, Servers, Notes, FeatureRequests, Documentations) with `AuthMiddleware` (accessible by Administrator, Developer, Viewer roles).
   - Wrapped POST, PUT, DELETE endpoints with `RoleMiddleware("Administrator", "Developer")` to prevent Viewer write access.

---

## 6. Routing Adjustments and Bug Fixes

We resolved routing layout issues:
1. **Removed Literal Routing Bug**:
   - Refactored routes to map `PUT` and `DELETE` on the root router using the dynamic parameter `/:id` (e.g., `/api/feature-requests/:id` and `/api/notes/:id`).
2. **Removed Trailing Slash Redirection**:
   - Registered relative paths inside route groups as `""` (empty string) instead of `"/"` to prevent Gin from redirecting URL requests without trailing slashes.
3. **Implemented Missing Note Controllers**:
   - Added `UpdateNote` and `DeleteNote` to `note_controller.go` to support the PUT/DELETE endpoints requested in `note_routes.go`.
4. **Import Cycle Resolution in Tests**:
   - Changed package declaration of integration tests (e.g. `feature_request_controller_test.go`, `note_controller_test.go`, and `documentation_controller_test.go`) to `package controllers_test` to bypass Go's compiler import cycles checks.

---

## 7. Verification & Testing

### Automated Unit Tests
We wrote a new unit test suite [auth_controller_test.go](file:///f:/GO/it-dashboard/backend/controllers/auth_controller_test.go) and refactored all existing tests under `controllers_test` to include the JWT authorization header. We also tested role limits by asserting `403 Forbidden` status codes for Viewers trying to execute POST write endpoints.

Run output for `go test -count=1 -v ./controllers`:
```
=== RUN   TestAuthRegister
--- PASS: TestAuthRegister (0.11s)
=== RUN   TestAuthLogin
--- PASS: TestAuthLogin (0.15s)
=== RUN   TestGetDocumentationsBySystemID
--- PASS: TestGetDocumentationsBySystemID (0.05s)
=== RUN   TestGetDocumentationByID
--- PASS: TestGetDocumentationByID (0.05s)
=== RUN   TestCreateDocumentation
--- PASS: TestCreateDocumentation (0.10s)
=== RUN   TestUpdateDocumentation
--- PASS: TestUpdateDocumentation (0.05s)
=== RUN   TestDeleteDocumentation
--- PASS: TestDeleteDocumentation (0.05s)
=== RUN   TestGetFeatureRequestsBySystemID
--- PASS: TestGetFeatureRequestsBySystemID (0.05s)
=== RUN   TestCreateFeatureRequest
--- PASS: TestCreateFeatureRequest (0.10s)
=== RUN   TestUpdateFeatureRequest
--- PASS: TestUpdateFeatureRequest (0.05s)
=== RUN   TestDeleteFeatureRequest
--- PASS: TestDeleteFeatureRequest (0.05s)
=== RUN   TestGetNotesBySystemID
--- PASS: TestGetNotesBySystemID (0.05s)
=== RUN   TestCreateNote
--- PASS: TestCreateNote (0.11s)
=== RUN   TestUpdateNote
--- PASS: TestUpdateNote (0.06s)
=== RUN   TestDeleteNote
--- PASS: TestDeleteNote (0.06s)
=== RUN   TestGetServers
--- PASS: TestGetServers (0.05s)
=== RUN   TestGetServerByID
--- PASS: TestGetServerByID (0.06s)
=== RUN   TestCreateServer
--- PASS: TestCreateServer (0.12s)
=== RUN   TestUpdateServer
--- PASS: TestUpdateServer (0.05s)
=== RUN   TestDeleteServer
--- PASS: TestDeleteServer (0.05s)
PASS
ok  	backend/controllers	2.294s
```

---

## 8. React Modern Enterprise Frontend Implementation

We successfully built, configured, and compiled the React-TypeScript-Vite frontend client under `/frontend`.

### Key Features Implemented:
1. **Design Token System (`src/index.css`)**:
   - Outfitted with a premium enterprise dark theme (obsidian colors, cobalt gradients, and subtle hover borders).
   - Styled cards, buttons, custom TIMELINE logs for developer notes, and form fields out-of-the-box.
   - Built a custom responsive layout shell featuring a fixed Sidebar and absolute Header.
2. **Global Auth State (`src/context/AuthContext.tsx` & `src/services/api.ts`)**:
   - `AuthContext`: Feeds user credentials, token, authentication status, and logout actions to components. Fetches `GET /api/auth/me` to refresh profiles automatically on page reload.
   - `Axios Interceptors`: Injects `Authorization: Bearer <token>` into outgoing requests. Clears localStorage and forces a logout redirect if the backend returns `401 Unauthorized` (token expired/invalid).
3. **Type-Safe Routing (`src/router.tsx`)**:
   - Set up `@tanstack/react-router` with code-defined route trees (`/login`, `/dashboard`, `/dashboard/servers`, `/dashboard/systems/$systemId`, and `/dashboard/documentations/new`).
   - Secured `/dashboard` routes to require active user authentication.
4. **Data Views**:
   - **Login**: A premium glassmorphic credentials portal.
   - **Systems Table**: Lists systems with tech stack, online status badges, links, and hosts.
   - **Servers Manager**: Lets administrators and developers perform full CRUD on server node configurations.
   - **System Detail View**: Uses Tab headers to group system specifications, linked host details, timeline-style Note logs, feature request cards, and documentation guides.
   - **Global Add Documentation**: Implemented the shortcut requested by the user, enabling direct documentation publishing for any system from a centralized wizard.
5. **Role-Based Client Security**:
   - Disables or completely hides creation/update forms and DELETE buttons if the logged-in user is a `Viewer` (ensuring read-only access in compliance with backend middleware).

### Verification
Successfully compiled and verified production bundler correctness with zero errors:
```bash
vite v8.0.16 building client environment for production...
✓ built in 228ms
dist/assets/index-CPPT6pxN.css   11.11 kB
dist/assets/index-CB0ynm4_.js   455.61 kB
```

