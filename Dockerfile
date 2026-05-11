FROM node:25-alpine3.22 AS frontend-build
WORKDIR /app
COPY package.json .
COPY package-lock.json .
RUN npm ci
COPY . .
RUN npx @tailwindcss/cli -i ./static/css/main.css -o ./static/css/output.css

FROM golang:1.26.2-alpine3.23 AS backend-build
WORKDIR /app
RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY --from=frontend-build /app .
RUN go mod download
ENV CGO_ENABLED=1
RUN go build -o seesaw .

FROM alpine:3.23
WORKDIR /app
RUN apk add --no-cache sqlite-libs
COPY --from=backend-build /app/seesaw .
CMD [ "./seesaw" ]