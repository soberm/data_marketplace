FROM node:13-alpine as builder

LABEL maintainer="michael.sober@ymail.com"
LABEL version="1.0"
LABEL description="Dockerfile for the angular frontend of the data marketplace"

COPY package.json yarn.lock patch-webpack.js ./

RUN yarn install && mkdir /ng-app && mv ./node_modules ./ng-app

WORKDIR /ng-app

COPY . .

RUN npm run ng build -- --prod --output-path=dist

FROM nginx:1.14.1-alpine

COPY nginx/default.conf /etc/nginx/conf.d/

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /ng-app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
