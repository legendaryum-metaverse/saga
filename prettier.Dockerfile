ARG BUN_VERSION=alpine
FROM oven/bun:${BUN_VERSION} AS base
WORKDIR /app

FROM base AS runner

RUN bun init -y
RUN bun add prettier

ENTRYPOINT ["bunx", "prettier", "--", "\"**/*.{md,yml,yaml,json,mdx}\"", "--single-quote=true", "--write"]
