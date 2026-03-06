# CLAUDE.md — React Project

## Project Context
- **Name:** [project_name]
- **React:** 19+ (functional components only)
- **TypeScript:** strict mode
- **Build:** Vite
- **CSS:** Tailwind CSS 4
- **State:** Zustand (global) + React state (local)
- **Router:** TanStack Router
- **Backend:** [API base URL]

---

## API Layer
```typescript
// Todas as chamadas passam por services/api.ts
// Base URL: import.meta.env.VITE_API_URL
// HTTP client: ky ou axios com interceptors
// Auth: Bearer token via interceptor
// Server state: TanStack Query (useQuery, useMutation)
// Nunca fetch direto em componentes — sempre via hooks customizados
```

---

## Core Libraries

### TanStack Query (React Query)
- Cache: staleTime 5min para listas, Infinity para dados estaticos
- Mutations: sempre invalidar queries relacionadas no onSuccess
- Keys: usar queryKeyFactory pattern

```typescript
// GOOD — query key factory
export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (filters: UserFilters) => [...userKeys.lists(), filters] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
}

// GOOD — custom hook
export function useUsers(filters: UserFilters) {
  return useQuery({
    queryKey: userKeys.list(filters),
    queryFn: () => api.users.list(filters),
  })
}

// BAD — query inline no componente
useQuery({ queryKey: ['users'], queryFn: fetch('/api/users') })
```

### TanStack Router
- File-based routing via `@tanstack/router-plugin`
- Type-safe routes: sempre usar `Link` tipado, nunca `<a href>`
- Loaders: prefetch data no route loader quando possivel
- Search params: validar com zod schema

```typescript
// GOOD
<Link to="/users/$userId" params={{ userId: user.id }}>
  {user.name}
</Link>

// BAD
<a href={`/users/${user.id}`}>{user.name}</a>
```

### TanStack Table
- Usar para todas as tabelas com sorting/filtering/pagination
- Column definitions tipadas com `createColumnHelper`
- Server-side pagination: controlar via query params

### Framer Motion
- Usar para animacoes de entrada/saida, transicoes de pagina, micro-interacoes
- Preferir `motion.div` com variants para animacoes reutilizaveis
- Sempre definir `initial`, `animate`, `exit` para AnimatePresence
- Nunca animar layout properties (width, height) — usar `layoutId` ou transforms

```typescript
// GOOD — variants reutilizaveis
const fadeIn = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -20 },
}

<motion.div {...fadeIn} transition={{ duration: 0.3 }}>
  {children}
</motion.div>

// BAD — animacao inline sem padrao
<motion.div animate={{ opacity: 1 }}>{children}</motion.div>
```

---

## State Management

### Zustand (global state)
```typescript
// stores/auth.store.ts
interface AuthState {
  user: User | null
  token: string | null
  login: (credentials: LoginDTO) => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  devtools(
    persist(
      (set) => ({
        user: null,
        token: null,
        login: async (credentials) => { ... },
        logout: () => set({ user: null, token: null }),
      }),
      { name: 'auth-storage' }
    )
  )
)
```

### Regras de estado
- Server state: TanStack Query (nunca duplicar em Zustand)
- UI state local: useState/useReducer
- UI state global (theme, sidebar): Zustand
- Form state: React Hook Form
- URL state: TanStack Router search params

---

## Form Validation
- Library: `react-hook-form` + `@hookform/resolvers` + `zod`
- Schema: sempre definir em `schemas/[feature].schema.ts`
- Nunca validar manualmente

```typescript
// GOOD
const schema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
})

const { register, handleSubmit, formState: { errors } } = useForm({
  resolver: zodResolver(schema),
})

// BAD
if (!email.includes('@')) setError('Invalid email')
```

---

## Component Patterns
```typescript
// Componentes: sempre function + arrow, nunca class
// Props: sempre interface tipada, nunca inline
// Exportar como named export, nunca default (exceto pages)

// GOOD
interface UserCardProps {
  user: User
  onEdit: (id: string) => void
}

export function UserCard({ user, onEdit }: UserCardProps) {
  return (...)
}

// BAD — default export, props inline
export default function({ user }: { user: any }) { ... }
```

---

## Project Structure
```
src/
  assets/
  components/          # shared/reusable (Button, Modal, DataTable)
    ui/                # design system primitives
  features/            # feature modules
    users/
      components/      # feature-specific components
      hooks/           # feature-specific hooks
      schemas/         # zod schemas
      services/        # API calls
      types/           # feature types
  hooks/               # shared hooks (useDebounce, useMediaQuery)
  layouts/
  lib/                 # utilities, helpers, constants
  routes/              # TanStack Router route files
  stores/              # Zustand stores
  styles/              # global CSS, Tailwind config
  types/               # shared TypeScript types
```

---

## Environment Variables
```bash
VITE_API_URL=http://localhost:8000/api/v1
VITE_APP_NAME=
VITE_SENTRY_DSN=        # opcional
```

---

## Build & Dev
```bash
npm run dev            # dev server (Vite)
npm run build          # producao
npm run preview        # preview build
npm run lint           # eslint
npm run type-check     # tsc --noEmit
```

---

## Known Gotchas
<!--
- TanStack Router: file-based routing requer plugin Vite — nao esquecer de configurar
- Framer Motion: AnimatePresence precisa de `key` unica nos children
- Zustand persist: cuidado com migrations quando mudar a shape do store
- React 19: use() hook e server components mudam patterns — verificar compatibilidade
- TanStack Query v5: queryKey deve ser readonly array
-->
