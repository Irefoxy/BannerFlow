FROM postman/newman_alpine33
WORKDIR /app
COPY ./tests/postman_collection ./tests
ENTRYPOINT newman run tests/test.postman_collection.json --env-var "SERVER_ADDR=$SERVER_ADDR"
