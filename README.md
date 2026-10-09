# Pizza Shop Order

A single-restaurant pizza shop: customers browse the menu, check out, and pay with Paystack. Staff manage pizzas, categories, images, and orders from an admin dashboard.

One Go process serves the website and the API. The HTML, CSS, and JavaScript are embedded in the binary. Postgres holds the data. Uploaded pizza images are written to `./uploads` and served at `/uploads`.

## Pages

| URL | Who uses it |
|---|---|
| `/` | Shop homepage and menu |
| `/pizza` | One pizza |
| `/checkout` | Cart and payment |
| `/payment-callback` | Paystack return page |
| `/userlogin` | Customer signup and login |
| `/my-orders` | A customer's orders |
| `/login` | Admin login |
| `/admin` | Menu, categories, and images |
| `/activity` | Orders and dashboard |
| `/health` | Process and database check |

Each page is also available with a `.html` suffix.

## Requirements

- Go 1.25 or newer
- Postgres (local, or Supabase)
- A Paystack secret key, if you take card payments

## Configure

```bash
cp .env.example .env
```

| Variable | Purpose |
|---|---|
| `DATABASE_URL` | Full Postgres URL. Used instead of the `DB_*` variables when set. |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Postgres connection. |
| `DB_SSLMODE` | `require` for Supabase. Defaults to `require` when unset. |
| `JWT_SECRET` | Signs admin and customer sessions. Required when `GIN_MODE=release`. |
| `PAYSTACK_SECRET_KEY` | Paystack **secret** key (`sk_live_...` or `sk_test_...`). |
| `PAYSTACK_CALLBACK_URL` | Public URL of `/payment-callback.html`. |
| `ALLOWED_ORIGINS` | Extra browser origins, comma-separated. |
| `TRUSTED_PROXIES` | Proxy CIDRs when the host sits behind a load balancer. |
| `PORT` | Listen port. Defaults to `8080`. |
| `GIN_MODE` | `debug` locally, `release` in production. |

On Supabase, use the IPv4 pooler host (`*.pooler.supabase.com`, port `5432`, user `postgres.<project-ref>`). The direct `db.<project-ref>.supabase.co` host is IPv6-only and fails on many networks.

## Database

On an empty database, run the SQL files in `migrations/` in this order:

1. `000_full_schema_supabase.sql`
2. `002_add_payment_columns.sql`
3. `005_add_categories_and_customers.sql`
4. `006_add_reviews.sql`
5. `007_add_category_images.sql`

`000_full_schema_supabase.sql` already creates the core tables. The later files add columns and tables that script does not include. Skip a file if that change is already present.

## Run locally

```bash
go run .
```

Open `http://127.0.0.1:8080`. `GET /health` returns `{"status":"ok"}` when Postgres answers.

## Admin account

There is no admin signup page. Create an account from the project root (this reads `.env` and stores a bcrypt hash):

```bash
go run ./cmd/seed -username=you@example.com -password='at-least-8-chars'
```

Log in at `/login` with that username and password. The command refuses to replace an existing username. Delete that row in `admins` first if you need to reset the password.

## API

Public:

- `GET /menu`, `GET /api/pizzas/:id`, `GET /api/pizzas/:id/images`, `GET /api/search`
- `GET /api/categories`, `GET /api/reviews`
- `POST /api/orders`
- `POST /api/auth/login`
- `POST /api/customers/signup`, `POST /api/customers/login`
- `POST /api/payments/initialize`, `GET /api/payments/verify/:reference`, `POST /api/payments/webhook`

Customer session (Bearer token, 30 days):

- `GET /api/customers/orders`
- `POST /api/reviews` (only for a delivered order belonging to that customer)

Admin session (Bearer token, 12 hours):

- `GET /api/auth/verify`
- `POST /api/pizzas`, `PUT /api/pizzas/:id`, `DELETE /api/pizzas/:id`
- `POST /api/pizzas/:id/images`, `DELETE /api/images/:id`
- `GET /api/orders`, `PATCH /api/orders/:id/status`
- `GET /api/admin/stats`, `GET /api/admin/best-seller`, `GET /api/admin/search`
- `POST /api/categories`, `PUT /api/categories/:id`, `DELETE /api/categories/:id`

Payment amounts are read from the database. The Paystack webhook checks `X-Paystack-Signature`.

## Docker

```bash
docker build -t pizza-shop .
docker run --env-file .env -p 8080:8080 pizza-shop
```

The image builds a static Linux binary, copies the frontend, listens on `8080`, and sets `GIN_MODE=release`. It runs as a non-root user.

## Deploy

Use the `Dockerfile` on Render, Fly.io, Railway, or any host that runs a container. Set the same environment variables as `.env`. Do not commit `.env`.

Set `PAYSTACK_CALLBACK_URL` to `https://<your-host>/payment-callback.html`.

`GET /health` is the readiness check. A free Render instance sleeps when idle, so the first request after that can take a minute.

`./uploads` lives on the container disk. A new deploy replaces that disk unless the host has a persistent volume, so uploaded images disappear on redeploy.

## Project layout

```
main.go            server, embedded pages, /health
routes/            HTTP routes
handlers/          request handlers
repositories/      SQL
middleware/        admin and customer sessions
migrations/        Postgres schema
cmd/seed/          create an admin user
frontend/          HTML, CSS, and JavaScript
uploads/           pizza images written at runtime
Dockerfile
```
