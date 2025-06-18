# ChatGPT Codex Chess App

This repository contains a minimal chess application with a Go backend and a React frontend.

## Backend

The Go server lives in the `server` directory. It exposes two endpoints:

- `GET /board` – returns the current board state as an array of strings.
- `POST /move` – accepts JSON `{"from": "e2", "to": "e4"}` to move a piece.

The server does not validate chess rules; it simply moves pieces on the board.

Run the backend:

```bash
cd server
go run .
```

This starts the server on port `8080`.

## Frontend

The frontend is a very small React app in the `client` directory. Open `index.html` in a browser after starting the backend. The board is displayed and you can select a piece then a target square to move it.

No build step is required because React is loaded from a CDN.

### Running the frontend with Bun

The `client` folder now contains a `package.json` with a small dev script that
uses `bunx serve` to host the static files. After installing [Bun](https://bun.sh),
start the React app with:

```bash
cd client
bun install
bun run dev
```

This serves `index.html` on <http://localhost:3000>.

