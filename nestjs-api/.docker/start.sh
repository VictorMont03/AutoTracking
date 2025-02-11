#!/bin/bash

npm install

# Runs migrations
npx prisma generate

# Run kafka consumer
npm run start:dev -- --entryFile=cmd/kafka.cmd &

# Run the app
npm run start:dev 

tail -f /dev/null 