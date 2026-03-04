---
applyTo: "ui/**/*.{ts,tsx}"
---

# React/TypeScript Coding Instructions

## Component Patterns

Define props interfaces above the component, use functional components:

```tsx
interface ItemCardProps {
    id: string
    label: string
    onSelect: (id: string) => void
}

export function ItemCard({ id, label, onSelect }: ItemCardProps) {
    return <button onClick={() => onSelect(id)}>{label}</button>
}
```

## State Management

- Use `useState` for client-only UI state
- Docker Desktop extension UI reinitializes on tab switch -- all persistent state must come from the backend via `ddClient`
- **No `useEffect` for data syncing** -- use nullable local override instead:

```tsx
// GOOD: nullable override, no useEffect sync
const [localOverride, setLocalOverride] = useState<string | null>(null)
const displayValue = localOverride ?? serverValue ?? ''
```

## API Integration

Always use `ddClient` from `@docker/extension-api-client`; never call `fetch` directly:

```tsx
import { createDockerDesktopClient } from '@docker/extension-api-client'
const ddClient = createDockerDesktopClient()

// GET from backend
const result = await ddClient.extension.vm.service.get('/api/checks')

// POST to backend
await ddClient.extension.vm.service.post('/api/check-all', {})
```

## Union Return Types

When a function returns a union type, add a type guard and use it at EVERY call site:

```tsx
// BAD: TS2339 -- access_token not on union
const result = await loginApi(user, pass)
setToken(result.access_token)

// GOOD: narrow first
const result = await loginApi(user, pass)
if (isMFAChallenge(result)) throw new Error('unexpected MFA')
setToken(result.access_token)
```

## TypeScript

- Strict mode -- no `any`; use `unknown` with type guards
- JSX short-circuit: `{expanded && item.details != null && <div/>}` (use `!= null`, not bare `&&` on `unknown`)
- Unused imports: ESLint catches these even when `tsc` does not -- verify every named import is used

## React Compiler Lint

Do not mutate `ref.current` during render -- wrap in `useEffect`:

```tsx
// BAD: ref mutation during render
onMessageRef.current = onMessage

// GOOD: wrap in effect
useEffect(() => { onMessageRef.current = onMessage }, [onMessage])
```

For Popper/Popover anchor elements, use callback ref with `useState` instead of `useRef`:

```tsx
// GOOD: callback ref avoids reading ref.current during render
const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null)
<Button ref={setAnchorEl}>Menu</Button>
<Popper anchorEl={anchorEl} open={open}>...</Popper>
```

## MUI v5 (Docker Desktop Extensions)

This project uses MUI v5 pinned via `@docker/docker-mui-theme`. Always reference MUI v5 patterns:

```tsx
// GOOD: MUI v5 TextField adornment
<TextField InputProps={{ startAdornment: <InputAdornment position="start">...</InputAdornment> }} />

// BAD: MUI v6 pattern -- does NOT work with v5
<TextField slotProps={{ input: { startAdornment: ... } }} />
```

Use `@docker/docker-mui-theme` for theming consistency with Docker Desktop:

```tsx
import { DockerMuiThemeProvider } from '@docker/docker-mui-theme'

function App() {
    return (
        <DockerMuiThemeProvider>
            <CssBaseline />
            {/* your app */}
        </DockerMuiThemeProvider>
    )
}
```

## Testing

Vitest + Testing Library. Use `mergeConfig` to inherit Vite plugins and defines:

```ts
// vitest.config.ts
import { mergeConfig, defineConfig } from 'vitest/config'
import viteConfig from './vite.config'

export default mergeConfig(viteConfig, defineConfig({
    test: { environment: 'jsdom', globals: true, setupFiles: './src/test-setup.ts' }
}))
```

Mock `@docker/extension-api-client` via resolve alias in `vitest.config.ts` (CJS/ESM mismatch):

```ts
resolve: {
    alias: {
        '@docker/extension-api-client': path.resolve(__dirname, 'src/__mocks__/@docker/extension-api-client.ts')
    }
}
```
