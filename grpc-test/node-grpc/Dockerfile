FROM node:14

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY src ./src
COPY protobuf ./protobuf

# 시작 스크립트 추가
COPY start.sh ./
RUN chmod +x start.sh

CMD ["./start.sh"]