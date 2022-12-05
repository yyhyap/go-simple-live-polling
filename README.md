# Simple Live Polling web app built with Go and React
 To simulate Malaysia's 15th General Election live voting using websocket and MongoDB
 
 ## Tech stacks
 - Go
 - Websocket
 - Zap
 - MongoDB
 - React

## Environment variables
Requires 2 environment variables for .env file in _server_ folder:

```sh
PORT
MONGODB_URL
```

## Docker
Kindly change the variable _projectDirName_ in file _dotEnvUtil.go_ and _mongoDbConnection.go_ to "" before docker build.