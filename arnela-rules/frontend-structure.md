# Frontend Structure

```
frontend/
├── app/                          # App Router (Next.js 16)
│   ├── layout.tsx               # Root layout
│   ├── page.tsx                 # Landing page
│   ├── (auth)/                  # Auth routes group
│   ├── (client)/                # Client area group
│   └── (backoffice)/            # Admin/Employee group
├── components/
│   ├── ui/                      # Shadcn components
│   ├── common/                  # Shared components
│   ├── client/                  # Client-specific
│   └── backoffice/              # Admin/Employee
├── stores/                       # Zustand stores
│   ├── useAuthStore.ts
│   ├── useClientStore.ts
│   └── ...
├── hooks/                        # Custom hooks
│   ├── useDebounce.ts
│   ├── useAuth.ts
│   └── ...
├── lib/                          # Utilities
│   ├── api.ts                   # API client
│   ├── constants.ts
│   ├── validators.ts
│   └── formatters.ts
└── types/                        # TypeScript types
    ├── user.ts
    ├── client.ts
    └── ...
```
