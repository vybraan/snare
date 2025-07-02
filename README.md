# snare
Snare is a modular authentication package for Go web applications built with the Gin framework. It provides session-based authentication with server-rendered HTML responses using the Templ component system.

## Features

- Session-based authentication using secure HTTP-only cookies
- Login, registration, and logout endpoints
- Middleware for protected routes
- SQLite storage for user data
- Password hashing with bcrypt
- Separated API and HTML views

## Routes
### HTML Endpoints

* `GET /auth/login` – Login page
* `GET /auth/register` – Registration page
* `GET /auth/logout` – Logout (requires session)

### API Endpoints

* `POST /api/auth/login` – Authenticate user
* `POST /api/auth/register` – Register new user
* `GET /api/auth/logout` – Invalidate session
* `GET /api/auth/state` – Session status (protected)

## Security

* Session data stored using `gin-contrib/sessions` with cookie backend
* Passwords are hashed with bcrypt
* CSRF protection and cookie flags should be handled at the reverse proxy level

## License

MIT License
