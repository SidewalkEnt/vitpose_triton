FROM node:14

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY src ./src
COPY protobuf ./protobuf

CMD ["npx", "ts-node", "src/grpc_client.ts"]