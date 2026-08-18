# Multi-stage Dockerfile for DeepPhotos SvelteKit Frontend
# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Stage 2: Production - only install production dependencies (no devDeps)
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=5173
ENV HOST=0.0.0.0

# Copy only what's needed to run the built app
COPY --from=builder /app/package.json ./
COPY --from=builder /app/package-lock.json ./

# Install ONLY production runtime deps (strips ~350MB of dev packages)
RUN npm ci --omit=dev --ignore-scripts

COPY --from=builder /app/build ./build
COPY --from=builder /app/static ./static

EXPOSE 5173

CMD ["node", "build"]
