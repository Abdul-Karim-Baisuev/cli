## Mission

This repository powers the Pickup Haul / AutoBais LLC website for pickup-truck
hauling, delivery, junk removal, moving help, and assembly/install services in
Southern California.

Treat it as a production business system. Prioritize:

1. clear customer flows;
2. trustworthy public copy;
3. reliable booking and payment behavior;
4. privacy and security;
5. minimal disruption to existing working functionality.

## Operating Style

Work like a senior engineer who already understands the project context.

- Inspect the repository before making changes.
- Prefer completing clear tasks without unnecessary clarification.
- Ask one targeted question only when a missing decision could materially change
  the implementation, public behavior, security, pricing, or legal meaning.
- Make reasonable, conservative assumptions for minor implementation details.
- Preserve existing architecture, naming, style, and user work.
- Prefer focused, reviewable changes over broad rewrites.
- Explain important assumptions and report verification honestly.
- Never claim that a feature works unless it was verified at the relevant level.

## Standard Workflow

For every task:

1. Read the request and identify the affected user flow.
2. Search the codebase before creating components, hooks, utilities, routes,
   migrations, or API handlers.
3. Inspect related files and existing patterns.
4. Make the smallest coherent change that fully solves the task.
5. Check edge cases, loading states, error states, and sensitive-data handling.
6. Run the narrowest verification that proves the change.
7. Run the full verification set when the change affects shared or production
   behavior.
8. Summarize what changed, what was verified, and what could not be verified.

Do not stop at a partial implementation when the remaining work is clear and
safe to complete.

## Project Map

- `src/` contains the React 18 + Vite + TypeScript frontend.
- `src/components/ui/` contains shadcn/Radix primitives. Keep app-specific
  business UI outside that folder when practical.
- `src/pages/` contains route-level screens, including admin, payment, tracking,
  quote, auth, and service pages.
- `src/integrations/` contains Lovable/Supabase client code.
- `supabase/functions/` contains Edge Functions for booking, payments, maps,
  email, Telegram, and AI widget flows.
- `supabase/migrations/` contains database schema changes.
- `docs/security/` and `README.md` document credential rotation and secret
  scanning.

## Editing Principles

Prefer:

- extending existing components and utilities;
- composition over duplication;
- simple, strongly typed implementations;
- established repository patterns;
- readable code over unnecessary abstraction;
- backward-compatible changes unless a breaking change is explicitly requested.

Avoid:

- replacing the Lovable/Vite scaffold;
- unrelated refactors;
- duplicate components, utilities, routes, or data models;
- premature abstractions;
- introducing new dependencies when existing tools are sufficient;
- broad formatting changes unrelated to the task.

Do not add a production dependency unless it is clearly necessary. When a new
dependency is avoidable, use the existing stack.

## React and TypeScript

- Use functional React components.
- Keep components focused and separate business logic from presentation where
  practical.
- Extract reusable stateful logic into hooks when it improves clarity.
- Preserve strict TypeScript typing.
- Avoid `any`; use explicit types, generics, or `unknown` with validation.
- Reuse existing UI primitives before creating new ones.
- Preserve responsive behavior and mobile usability.
- Include appropriate loading, empty, error, and disabled states.
- Avoid unnecessary rerenders, duplicated state, and oversized client-side
  dependencies.

## Supabase and Edge Functions

- Prefer migrations over manual database changes.
- Never edit generated types manually.
- Keep Edge Functions small, explicit, and fail-safe.
- Validate external input at function boundaries.
- Preserve backward compatibility unless explicitly requested otherwise.
- Do not assume a local frontend proves a deployed function or remote integration
  works.
- Verify the real frontend-to-API-to-function path when credentials and access
  are available.
- Treat database policies, authentication, service-role usage, and public
  endpoints as security-sensitive.

## Security and Privacy

- Never print, commit, expose, or move secrets from `.env`, Supabase, Lovable,
  Stripe, Telegram, Mapbox, GitHub, or deployment platforms.
- Treat bookings, customers, driver tracking, payments, admin/CRM data, Telegram
  notifications, and AI chat logs as sensitive.
- Do not add PII to analytics, logs, test fixtures, screenshots, examples, or
  public files.
- Use synthetic data for tests and examples.
- Preserve authorization checks.
- Do not weaken RLS, authentication, input validation, rate limits, or secret
  handling to make a feature easier to implement.
- Report suspected secret exposure or unsafe configuration immediately without
  reproducing the secret value.

## Public Content and Business Trust

- Keep public business copy practical, credible, and easy to understand.
- Do not invent testimonials, certifications, insurance claims, pricing,
  guarantees, service coverage, or legal promises.
- Use `deposit` and `down payment` wording carefully so customer expectations are
  clear.
- Confirm with the user before publishing personal photos, work photos, private
  customer details, or any public-facing claim that could affect trust, pricing,
  insurance, or legal responsibility.
- Preserve consistent business naming, contact details, and service-area
  information across the site.

## Production-Critical Flows

Treat these areas as high risk:

- booking and quote submission;
- payments, deposits, refunds, and Stripe flows;
- authentication and authorization;
- driver tracking;
- admin and CRM features;
- customer notifications;
- pricing and public claims;
- database migrations and RLS;
- Supabase Edge Functions;
- Telegram, email, maps, and AI integrations.

For changes in these areas, verify the complete affected workflow rather than
only the edited file.

Do not perform destructive production actions, rotate credentials, deploy schema
changes, publish content, or modify live payment behavior unless the task
explicitly requires it and the necessary access is available.

## Accessibility, SEO, and UX

For public-facing changes:

- preserve semantic HTML and keyboard navigation;
- keep labels and controls accessible;
- maintain sensible heading structure;
- preserve metadata, canonical behavior, and route discoverability;
- avoid layout shifts and unnecessarily large assets;
- keep forms clear and error messages actionable;
- verify important flows on mobile-sized layouts.

## Working Tree Safety

- Preserve unrelated local changes.
- This repository may contain untracked installers, exports, generated files, or
  local artifacts.
- Do not delete, move, reset, clean, or overwrite unrelated files unless asked.
- Do not use destructive Git commands such as `git reset --hard`,
  `git clean -fd`, or forced checkout to discard user work.
- Review the current diff before and after editing when Git is available.

## Commands

- `npm run dev` starts the Vite dev server.
- `npm run build` creates the production build.
- `npm run lint` runs ESLint.
- `npm run test` runs Vitest once.
- `npm run test:watch` starts Vitest watch mode.
- `npm run preview` serves the production build locally.
- `npm run secrets:verify` runs post-rotation smoke checks when the required
  secrets and admin token are available.
- `gitleaks detect --config .gitleaks.toml --no-banner` checks for leaked
  secrets.

## Verification

Use `npm run lint`, `npm run test`, and `npm run build` as the normal local
verification set for code changes.

Use narrower checks only when:

- the change is clearly isolated;
- the full command is unavailable;
- dependencies are missing;
- required credentials or remote access are unavailable.

Always report skipped verification and the reason.

Suggested minimums:

- documentation-only changes: inspect formatting and links;
- isolated UI changes: lint plus relevant tests when present;
- shared frontend logic: lint, tests, and build;
- booking, payment, auth, admin, or tracking changes: lint, tests, build, and
  workflow-level verification where access permits;
- migrations and Edge Functions: validate syntax, affected contracts, security,
  and integration paths.

A successful local dev server is not proof that deployed services work.

## Tooling

- Prefer `rg` for repository search.
- Prefer `apply_patch` for focused manual edits.
- Read nearby code before introducing new patterns.
- Check existing dependencies before adding libraries.
- Use repository-provided scripts instead of inventing parallel tooling.