FROM node:26-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

FROM node:26-alpine AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG PUBLIC_API_URL=http://localhost:8080
ENV PUBLIC_API_URL=$PUBLIC_API_URL

RUN npm run build
RUN npm prune --omit=dev

FROM node:26-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PUBLIC_API_URL=http://localhost:8080

RUN addgroup -S -g 1001 sveltekit
RUN adduser -S -u 1001 -G sveltekit sveltekit

COPY --from=builder --chown=sveltekit:sveltekit /app/build ./build
COPY --from=builder --chown=sveltekit:sveltekit /app/node_modules ./node_modules
COPY --chown=sveltekit:sveltekit package.json .

USER sveltekit

EXPOSE 3000
ENV PORT=3000
ENV HOST=0.0.0.0

CMD ["node", "build/index.js"]
