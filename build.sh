#!/bin/bash

cd view
npm install && npm run build
cd ..
go build -v
echo 'build rp done'
