FROM node:latest as build

WORKDIR /app

COPY package*.json ./

RUN npm install && npm install -g @angular/cli

COPY . .

RUN npm run build

FROM nginx:latest

COPY --from=build /app/dist/frontend/browser /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]