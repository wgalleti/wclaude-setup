# CLAUDE.md — Supabase Project

## Project Context
- **Name:** [project_name]
- **Supabase:** self-hosted / cloud
- **Frontend:** [React | Next.js | Vue | Nuxt | Flutter]
- **Auth:** Supabase Auth (email, OAuth, magic link)
- **Database:** PostgreSQL via Supabase
- **Storage:** Supabase Storage
- **Realtime:** [Yes/No]
- **Edge Functions:** [Yes/No]

---

## Supabase Client Setup
```typescript
// lib/supabase.ts (React/Next)
import { createClient } from '@supabase/supabase-js'
import type { Database } from '@/types/database.types'

export const supabase = createClient<Database>(
  process.env.NEXT_PUBLIC_SUPABASE_URL!,
  process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
)

// Para SSR (Next.js App Router)
import { createServerClient } from '@supabase/ssr'
import { cookies } from 'next/headers'

export async function createServerSupabase() {
  const cookieStore = await cookies()
  return createServerClient<Database>(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        getAll: () => cookieStore.getAll(),
        setAll: (cookiesToSet) => {
          cookiesToSet.forEach(({ name, value, options }) =>
            cookieStore.set(name, value, options)
          )
        },
      },
    }
  )
}
```

---

## Database (PostgreSQL)

### Queries tipadas
```typescript
// SEMPRE usar tipos gerados
// Gerar tipos: npx supabase gen types typescript --local > types/database.types.ts

// GOOD — query tipada
const { data, error } = await supabase
  .from('users')
  .select('id, email, full_name, created_at')
  .eq('organization_id', orgId)
  .order('created_at', { ascending: false })

// BAD — select * sem tipagem
const { data } = await supabase.from('users').select('*')
```

### Row Level Security (RLS)
```sql
-- SEMPRE habilitar RLS em todas as tabelas
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Policies claras e bem nomeadas
CREATE POLICY "Users can view own profile"
  ON users FOR SELECT
  USING (auth.uid() = id);

CREATE POLICY "Users can update own profile"
  ON users FOR UPDATE
  USING (auth.uid() = id)
  WITH CHECK (auth.uid() = id);

-- Multi-tenant: scope por organization
CREATE POLICY "Members can view org data"
  ON projects FOR SELECT
  USING (
    organization_id IN (
      SELECT organization_id FROM memberships
      WHERE user_id = auth.uid()
    )
  );
```

### Migrations
```bash
# Criar migration
npx supabase migration new [nome_descritivo]

# Aplicar localmente
npx supabase db reset

# Push para remote
npx supabase db push

# Nunca editar migrations ja aplicadas
# Uma concern por migration
```

---

## Auth
```typescript
// Sign up
const { data, error } = await supabase.auth.signUp({
  email,
  password,
})

// Sign in
const { data, error } = await supabase.auth.signInWithPassword({
  email,
  password,
})

// OAuth
const { data, error } = await supabase.auth.signInWithOAuth({
  provider: 'google',
  options: { redirectTo: `${origin}/auth/callback` },
})

// Listener de auth state
supabase.auth.onAuthStateChange((event, session) => {
  // handle auth state changes
})

// IMPORTANTE: sempre tratar erros de auth
// Nunca confiar apenas no client — RLS protege no server
```

---

## Storage
```typescript
// Upload
const { data, error } = await supabase.storage
  .from('avatars')
  .upload(`${userId}/avatar.png`, file, {
    cacheControl: '3600',
    upsert: true,
  })

// URL publica
const { data } = supabase.storage
  .from('avatars')
  .getPublicUrl(`${userId}/avatar.png`)

// URL assinada (privada, expira)
const { data } = await supabase.storage
  .from('documents')
  .createSignedUrl(`${userId}/contract.pdf`, 3600)

// Storage policies: configurar no dashboard ou via SQL
```

---

## Realtime
```typescript
// Subscribe a changes
const channel = supabase
  .channel('room1')
  .on(
    'postgres_changes',
    { event: 'INSERT', schema: 'public', table: 'messages' },
    (payload) => {
      console.log('New message:', payload.new)
    }
  )
  .subscribe()

// Cleanup
channel.unsubscribe()

// Broadcast (sem DB)
const channel = supabase.channel('room1')
channel.send({
  type: 'broadcast',
  event: 'cursor-pos',
  payload: { x: 100, y: 200 },
})
```

---

## Edge Functions
```typescript
// supabase/functions/hello/index.ts
import { serve } from 'https://deno.land/std@0.177.0/http/server.ts'
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2'

serve(async (req) => {
  const supabase = createClient(
    Deno.env.get('SUPABASE_URL')!,
    Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!,
  )

  // ... logic

  return new Response(JSON.stringify({ data }), {
    headers: { 'Content-Type': 'application/json' },
  })
})
```

```bash
# Deploy
npx supabase functions deploy hello

# Testar local
npx supabase functions serve hello
```

---

## Type Generation
```bash
# Gerar tipos do schema do banco
npx supabase gen types typescript --local > types/database.types.ts

# IMPORTANTE: regenerar apos cada migration
# Adicionar ao script de pre-build
```

---

## Environment Variables
```bash
# .env.local
NEXT_PUBLIC_SUPABASE_URL=http://127.0.0.1:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=
SUPABASE_SERVICE_ROLE_KEY=       # server-only, NUNCA expor no client

# Para Vue/Nuxt
VITE_SUPABASE_URL=
VITE_SUPABASE_ANON_KEY=
# ou NUXT_PUBLIC_SUPABASE_URL= etc
```

---

## Local Development
```bash
npx supabase start            # inicia stack local (DB, Auth, Storage, etc)
npx supabase stop             # para stack local
npx supabase status           # mostra URLs e keys locais
npx supabase db reset         # reset DB + aplica migrations + seed
npx supabase migration new X  # cria nova migration
npx supabase gen types typescript --local > types/database.types.ts
```

---

## Known Gotchas
<!--
- RLS: tabelas sem policies habilitadas sao acessiveis por qualquer usuario autenticado
- Service Role Key: NUNCA usar no client-side — bypassa RLS
- Realtime: precisa habilitar replication na tabela (supabase dashboard > Database > Replication)
- Auth: tokens expiram — usar onAuthStateChange para refresh automatico
- Storage: policies sao separadas das policies de tabela — configurar ambas
- Edge Functions: runtime Deno, nao Node.js — imports diferentes
- gen types: sempre regenerar apos migrations — tipos desatualizados causam bugs silenciosos
- self-hosted: configurar SMTP para emails de confirmacao/recuperacao
-->
