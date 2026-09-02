FROM node:22-alpine
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install --prefer-offline --no-audit --no-fund --loglevel=error
COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]