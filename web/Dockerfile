FROM node:19-alpine3.16 AS builder
WORKDIR /app
ENV NODE_ENV production
COPY package.json ./
RUN npm install
COPY . ./
RUN npm run build

FROM nginx:1.23.3-alpine
WORKDIR /usr/share/nginx/html
RUN rm -rf ./*
COPY --from=builder /app/build .
USER user:user
ENTRYPOINT ["nginx", "-g", "daemon off;"]