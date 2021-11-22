FROM node:12.17.0-alpine3.9

ENV APP_PATH /opt/nodejs-app

RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH
COPY build/nodejs-app $APP_PATH

RUN npm install

ENTRYPOINT node app.js
