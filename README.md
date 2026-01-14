# Product Discovery Web App (Go + React)

A full-stack product search app:
- **Backend:** Go + Postgres
- **Frontend:** React (Vite) + shadcn/ui + React Query
- **Auth:** email/password + Google OAuth, jwt

---

## Prerequisites

- Docker + Docker Compose (recommended)
- Node.js (LTS) + npm
- Go (only needed if running backend without Docker)


## Quick start (recommended): Docker Compose

### 1) Create `.env`
```bash
cp .env.example .env
```

### 2) Start Backend
```bash
docker compose up --build
```

### 3) Seed the database
```bash
docker compose exec backend ./seed
```

### 4) Start Frontend
```bash
cd frontend
npm install
npm run dev
```

---
## Testing
### Backend
```bash
go test ./... -v
```

### Frontend
```bash
npm run test
# or
npm run test:watch
```

Have fun exploring! 