#!/bin/bash

cd view
npm install && npm run build
cd ..
go build -i -v
echo 'build rp done'
