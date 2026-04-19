# Client

React + Vite single-page app using TanStack Router.

## Getting Started

```bash
npm install
npm run dev
```

## MSW (Optional)

Enable browser API mocking in development:

```powershell
$env:VITE_ENABLE_MSW='true'; npm run dev
```

Disable it in the current shell:

```powershell
Remove-Item Env:VITE_ENABLE_MSW
```

## Build

```bash
npm run build
```

## Tests

```bash
npm run test
```

## Linting & Formatting

```bash
npm run lint
npm run format
npm run check
```

## Routing

Routes are file-based under `src/routes`.
