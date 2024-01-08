# 使用 nginx 镜像作为基础镜像
FROM nginx:1.21-alpine

# 复制构建的前端文件到 nginx 容器中
COPY dist/ /usr/share/nginx/html
COPY nginx/default.conf /etc/nginx/conf.d/default.conf
COPY nginx/nginx.conf /etc/nginx/nginx.conf
# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
