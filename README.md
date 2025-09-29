## go-webapp

Personal playground project with a small Go backend and Next.js frontend. Not intended for production use.

### Stack
- **Backend**: Go (chi, oauth2)
- **Frontend**: Next.js 15, React 19, Tailwind CSS

### Local development
1. Backend
   - Create a `.env` in `backend/` (see variables below)
   - Run: `go run ./backend`
2. Frontend
   - `cd frontend`
   - Create a `.env.local` (see variables below)
   - Run: `npm install && npm run dev`

### Environment variables (examples)
Backend (`backend/.env`):
```
PORT=8080
FRONTEND_URL=http://localhost:3000
COOKIE_DOMAIN=localhost

# Optional OAuth
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT_URL=http://localhost:8080/api/auth/github/callback
```

Frontend (`frontend/.env.local`):
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Notes
- This repo is for personal use; no support is expected.
- See `LICENSE.md` for licensing.


