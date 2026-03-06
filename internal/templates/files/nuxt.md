# CLAUDE.md — Nuxt Project

## Project Context
- **Name:** [project_name]
- **Nuxt:** 3.x (Nitro server engine)
- **Vue:** 3 (Composition API, `<script setup>`)
- **TypeScript:** strict mode
- **CSS:** Tailwind CSS 4 + PrimeVue 4 (passthrough mode)
- **State:** Pinia (via @pinia/nuxt)
- **Backend:** [API base URL or Nitro server routes]

---

## Nuxt Auto-Imports
```typescript
// Nuxt auto-importa automaticamente:
// - Vue APIs: ref, computed, watch, onMounted, etc.
// - Nuxt composables: useFetch, useAsyncData, useRoute, useRouter, etc.
// - Components de components/, layouts/, pages/
// - Composables de composables/
// - Utils de utils/
// - Pinia stores de stores/

// GOOD — sem import necessario
const route = useRoute()
const { data } = await useFetch('/api/users')

// BAD — import desnecessario com Nuxt
import { ref } from 'vue'
import { useRoute } from 'vue-router'
```

---

## Data Fetching
```typescript
// Preferir useFetch para chamadas simples
const { data: users, status, refresh } = await useFetch('/api/users', {
  query: { page: currentPage },
})

// useAsyncData para logica mais complexa
const { data: user } = await useAsyncData(
  `user-${id}`,
  () => $fetch(`/api/users/${id}`)
)

// IMPORTANTE: useFetch e useAsyncData executam no server (SSR) por padrao
// Usar lazy: true para client-only fetch
const { data } = await useFetch('/api/stats', { lazy: true })

// BAD — nunca usar onMounted + fetch para data fetching
onMounted(async () => {
  const data = await $fetch('/api/users') // WRONG
})
```

---

## Server Routes (Nitro)
```typescript
// server/api/users.get.ts
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  // ... business logic
  return { users }
})

// server/api/users.post.ts
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  // validar com zod
  const data = createUserSchema.parse(body)
  // ... create
  return { user }
})

// server/middleware/auth.ts — middleware de server
export default defineEventHandler((event) => {
  // auth check
})
```

---

## PrimeVue + Nuxt
```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  modules: ['@primevue/nuxt-module'],
  primevue: {
    autoImport: true, // auto-import de todos os componentes
    options: {
      theme: {
        preset: Aura, // ou Lara, Nora
      },
    },
  },
})

// Componentes PrimeVue disponiveis globalmente sem import
// <DataTable>, <Column>, <Dialog>, <Toast>, etc.
```

---

## Composables Padrao
```typescript
// composables/useApi.ts — wrapper para $fetch com auth
export function useApi() {
  const config = useRuntimeConfig()
  const auth = useAuthStore()

  return $fetch.create({
    baseURL: config.public.apiUrl,
    headers: {
      Authorization: `Bearer ${auth.token}`,
    },
    onResponseError({ response }) {
      if (response.status === 401) auth.logout()
    },
  })
}

// composables/useToast.ts — PrimeVue toast
// composables/useConfirm.ts — PrimeVue confirm
// composables/useForm.ts — vee-validate + zod
```

---

## Form Validation
- Library: `@vee-validate/nuxt` + `zod`
- VeeValidate auto-importado via module Nuxt
- Schema: sempre em `schemas/[feature].schema.ts`

```typescript
// schemas/user.schema.ts
export const createUserSchema = z.object({
  email: z.string().email(),
  name: z.string().min(2),
  password: z.string().min(8),
})

// Usar <Form> e <Field> do vee-validate (auto-importados)
```

---

## Pinia Store Pattern
```typescript
// stores/auth.ts
export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = useCookie('auth-token')

  const isAuthenticated = computed(() => !!token.value)

  async function login(credentials: LoginDTO) {
    const { data } = await useFetch('/api/auth/login', {
      method: 'POST',
      body: credentials,
    })
    user.value = data.value.user
    token.value = data.value.token
  }

  function logout() {
    user.value = null
    token.value = null
    navigateTo('/login')
  }

  return { user, token, isAuthenticated, login, logout }
})
```

---

## Middleware (Route)
```typescript
// middleware/auth.ts — proteger rotas
export default defineNuxtRouteMiddleware((to, from) => {
  const auth = useAuthStore()
  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }
})

// Usar no page:
// definePageMeta({ middleware: 'auth' })
```

---

## Project Structure
```
app.vue
nuxt.config.ts

assets/
  css/tailwind.css

components/
  ui/                  # design system primitives
  forms/               # form components
  layout/              # Header, Sidebar, Footer

composables/           # useApi, useForm, etc (auto-imported)

features/              # feature modules (nao auto-imported)
  users/
    components/
    composables/
    schemas/
    types/

layouts/
  default.vue
  auth.vue

middleware/             # route middleware
  auth.ts

pages/
  index.vue
  login.vue
  users/
    index.vue
    [id].vue

plugins/               # Nuxt plugins

schemas/               # zod schemas compartilhados

server/
  api/                 # Nitro API routes
  middleware/           # server middleware
  utils/               # server utilities

stores/                # Pinia stores (auto-imported)

types/                 # TypeScript types

utils/                 # utility functions (auto-imported)
```

---

## Runtime Config
```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  runtimeConfig: {
    // Server-only
    databaseUrl: '',
    authSecret: '',
    // Public (client + server)
    public: {
      apiUrl: 'http://localhost:8000/api/v1',
      appName: '',
    },
  },
})
```

```bash
# .env
NUXT_DATABASE_URL=
NUXT_AUTH_SECRET=
NUXT_PUBLIC_API_URL=http://localhost:8000/api/v1
NUXT_PUBLIC_APP_NAME=
```

---

## Build & Dev
```bash
npm run dev            # dev server (Nitro)
npm run build          # producao
npm run preview        # preview build
npm run generate       # static site generation (SSG)
npm run lint           # eslint
npm run type-check     # nuxi typecheck
```

---

## Known Gotchas
<!--
- useFetch vs $fetch: useFetch deduplica e funciona no SSR; $fetch e raw fetch
- Auto-imports: se nao resolver, rodar `nuxi prepare` para regenerar tipos
- Server routes: nao importar codigo client-side (Vue, Pinia) em server/
- Pinia + SSR: usar useCookie para persistir estado entre server/client
- PrimeVue + Nuxt: usar @primevue/nuxt-module, nao instalar manual
- Nitro: cada arquivo em server/api/ gera um endpoint — cuidado com naming
- Middleware: definePageMeta e macro de compilacao — nao usar em composables
-->
