FROM node:20-alpine AS build

WORKDIR /app

# 只复制依赖配置，最小化缓存
COPY frontend/package.json frontend/tsconfig.json frontend/vite.config.ts ./
COPY frontend/src ./src
COPY frontend/index.html ./

# 【关键】加速命令：跳过审计、不弹广告、极速安装
RUN npm install --registry=https://registry.npmmirror.com && npm run build

FROM nginx:alpine

COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80