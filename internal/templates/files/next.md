# CLAUDE.md — Next.js Project

## Project Context
- **Name:** [project_name]
- **Next.js:** 15+ (App Router)
- **React:** 19+
- **TypeScript:** strict mode
- **CSS:** Tailwind CSS 4
- **State:** Zustand (client) + TanStack Query (server state)
- **Auth:** [NextAuth.js / Auth.js | Supabase Auth | Custom]
- **Deploy:** [Vercel | Docker | AWS]

---

## Architecture: App Router

### Server vs Client Components
```typescript
// Por padrao, tudo e Server Component no App Router
// Adicionar 'use client' APENAS quando necessario:
// - useState, useEffect, useContext
// - Event handlers (onClick, onChange)
// - Browser APIs (window, document)
// - Hooks de terceiros (useQuery, useForm, motion)

// GOOD — Server Component (padrao)
// app/users/page.tsx
export default async function UsersPage() {
  const users = await getUsers() // fetch direto no servidor
  return <UserList users={users} />
}

// GOOD — Client Component (quando precisa interatividade)
// components/user-search.tsx
'use client'
export function UserSearch() {
  const [query, setQuery] = useState('')
  return <input value={query} onChange={e => setQuery(e.target.value)} />
}

// BAD — 'use client' desnecessario
'use client'
export function UserCard({ user }: { user: User }) {
  return <div>{user.name}</div>  // nao precisa de client!
}
```

### Server Actions
```typescript
// Preferir Server Actions para mutacoes simples
// app/users/actions.ts
'use server'

export async function createUser(formData: FormData) {
  const data = createUserSchema.parse(Object.fromEntries(formData))
  await db.users.create(data)
  revalidatePath('/users')
}

// Para mutacoes complexas com loading state, usar TanStack Query no client
```

### Data Fetching
```typescript
// Server Components: fetch direto (com cache control)
const users = await fetch('/api/users', {
  next: { revalidate: 60 }, // ISR: revalidar a cada 60s
})

// Client Components: TanStack Query
'use client'
const { data: users } = useUsers(filters)

// Nunca usar useEffect + fetch para data fetching
```

---

## API Routes
```typescript
// app/api/[resource]/route.ts
// Usar Route Handlers para API endpoints

import { NextRequest, NextResponse } from 'next/server'

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams
  // ... logic
  return NextResponse.json(data)
}

// Validar input com zod em TODA route handler
```

---

## Core Libraries

### TanStack Query (client-side)
- Usar apenas em Client Components para dados que precisam de cache reativo
- Server Components: fetch direto, sem TanStack Query
- Pattern: hooks em `hooks/` ou `features/[feature]/hooks/`

### TanStack Table
- Sempre Client Component (precisa de interatividade)
- Server-side pagination via searchParams

### Framer Motion
- Usar para animacoes de entrada/saida e transicoes de layout
- Componentes com motion sao sempre Client Components
- Page transitions: usar `AnimatePresence` no layout

```typescript
// layouts com animacao
'use client'
import { motion, AnimatePresence } from 'framer-motion'

export function PageTransition({ children }: { children: React.ReactNode }) {
  return (
    <AnimatePresence mode="wait">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -20 }}
      >
        {children}
      </motion.div>
    </AnimatePresence>
  )
}
```

### Zustand (client state)
- Apenas para estado global de UI (theme, sidebar, preferences)
- Nunca duplicar server state no Zustand — usar TanStack Query

---

## Form Validation
- Library: `react-hook-form` + `@hookform/resolvers` + `zod`
- Server Actions: validar com zod no server tambem
- Schemas compartilhados entre client e server em `schemas/`

---

## Project Structure (App Router)
```
app/
  (auth)/              # route group: login, register
    login/page.tsx
    register/page.tsx
  (dashboard)/         # route group: area autenticada
    layout.tsx         # sidebar + header
    users/
      page.tsx         # lista
      [id]/page.tsx    # detalhe
      actions.ts       # server actions
  api/                 # Route Handlers
    users/route.ts
  layout.tsx           # root layout
  page.tsx             # home
  loading.tsx          # loading UI global
  error.tsx            # error boundary global
  not-found.tsx

components/
  ui/                  # design system (Button, Input, Modal)
  forms/               # form components reutilizaveis
  layout/              # Header, Sidebar, Footer

features/
  users/
    components/
    hooks/
    schemas/
    services/
    types/

hooks/                 # shared hooks
lib/                   # utils, constants, config
stores/                # Zustand stores
types/                 # shared types
middleware.ts          # auth middleware, redirects
```

---

## Middleware
```typescript
// middleware.ts — proteger rotas
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  // Auth check, redirects, headers
}

export const config = {
  matcher: ['/dashboard/:path*'],
}
```

---

## Environment Variables
```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8000/api/v1
NEXT_PUBLIC_APP_NAME=

# Server-only (sem NEXT_PUBLIC_)
DATABASE_URL=
AUTH_SECRET=
```

---

## Build & Dev
```bash
npm run dev            # dev server (turbopack)
npm run build          # producao
npm run start          # start producao
npm run lint           # next lint
npm run type-check     # tsc --noEmit
```

---

## Known Gotchas
<!--
- App Router: cuidado com 'use client' propagando — um client component torna todos os filhos client
- Server Actions: nao retornar objetos grandes — serializa tudo
- Middleware: roda no Edge Runtime — sem Node.js APIs (fs, path)
- fetch no Server Component: Next.js faz cache automatico — usar { cache: 'no-store' } se precisa fresh
- Image: sempre usar next/image, nunca <img> — otimizacao automatica
- Font: usar next/font, nunca importar de CDN
- Metadata: exportar metadata em cada page.tsx para SEO
-->
