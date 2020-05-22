#!/bin/bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"isdn":"xyz","title":"my book","author":"john doe","pages":15}' \
  http://localhost:8080/books