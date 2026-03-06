# CLAUDE.md — Vue.js / PrimeVue Project

## Project Context
- **Name:** [project_name]
- **Vue:** 3 (Composition API, `<script setup>`)
- **TypeScript:** [Yes/No]
- **Build:** Vite
- **CSS:** Tailwind CSS 3 + PrimeVue 4 (passthrough mode)
- **State:** Pinia
- **Router:** Vue Router 4
- **Backend:** [API base URL]

---

## API Layer
```typescript
// Todas as chamadas passam por services/api.ts
// Base URL: import.meta.env.VITE_API_URL
// Auth: Bearer token via interceptor Axios/ofetch
// Retry: 1x em 401 → refresh token → retry
```

---

## PrimeVue Usage
- **Theme:** [Aura | Lara | Nora] — unstyled: [Yes/No]
- **Passthrough:** customizações em `presets/` 
- **Regra:** Sempre usar componente PrimeVue existente antes de criar custom
- **Imports:** Auto-import via `unplugin-vue-components`

Componentes mais usados neste projeto:
```
DataTable, Column → listagens com paginação server-side
Dialog, DynamicDialog → modais
Toast → notificações (via useToast)
Form, InputText, Select → formulários
Breadcrumb → navegação
```

---

## Composables Padrão
```typescript
// Sempre disponíveis via auto-import (src/composables/)
useApi()         // wraps ofetch com auth + error handling
useToast()       // PrimeVue toast notifications  
useConfirm()     // PrimeVue confirm dialog
useAuth()        // user, isAuthenticated, login, logout
usePagination()  // server-side pagination helper
useForm(schema)  // vee-validate + zod integration
```

---

## Form Validation
- Library: `vee-validate` + `zod`
- Schema: sempre definir em `schemas/[feature].schema.ts`
- Nunca validar manualmente — sempre via schema

```typescript
// GOOD
const schema = z.object({
  email: z.string().email(),
  password: z.string().min(8)
})

// BAD
if (!email.includes('@')) { error = 'invalid' }
```

---

## Pinia Store Pattern
```typescript
// stores/user.store.ts
export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isAdmin = computed(() => user.value?.role === 'admin')
  
  async function fetchProfile() { ... }
  
  return { user, isAdmin, fetchProfile }
})
```

---

## i18n
- Library: `vue-i18n`
- Arquivos: `locales/pt-BR.json`, `locales/en.json`
- Nunca string hardcoded no template — sempre `$t('key')`

---

## Environment Variables
```bash
VITE_API_URL=http://localhost:8000/api/v1
VITE_APP_NAME=
VITE_SENTRY_DSN=        # opcional
```

---

## Build & Deploy
```bash
yarn dev           # dev server
yarn build         # produção
yarn preview       # preview build
yarn lint          # eslint
yarn type-check    # tsc --noEmit
```

---

## Nuxt Compatibility
<!--
Se migrar para Nuxt 3:
- Remover vue-router manual → Nuxt usa file-based routing (pages/)
- Remover imports de ref, computed, etc → Nuxt auto-importa
- Mover composables/ para raiz → Nuxt auto-importa
- Pinia: usar @pinia/nuxt module
- PrimeVue: usar @primevue/nuxt-module (auto-import de componentes)
- API calls: substituir useApi() por useFetch() para SSR
- Env vars: VITE_* → NUXT_PUBLIC_* (runtimeConfig)
- Considere usar nuxt.md template em vez deste para projetos Nuxt
-->

---

## Known Gotchas
<!--
- PrimeVue 4: `pt` (passthrough) substituiu `unstyled` props — usar pt:root, pt:label etc
- Tailwind + PrimeVue: adicionar PrimeVue ao content do tailwind.config.js
- Vue Router 4: lazy routes com defineAsyncComponent para chunks menores
-->
