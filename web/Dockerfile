FROM node:10 AS builder
WORKDIR /app
COPY package.json ./
RUN npm install 
COPY . ./
RUN npm run build

FROM nginx:1.19.0
WORKDIR /usr/share/nginx/html
RUN rm -rf ./*
COPY --from=builder /app/build .
USER user:user
ENTRYPOINT ["nginx", "-g", "daemon off;"]