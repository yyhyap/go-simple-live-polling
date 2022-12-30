# E-Voting web app built with Go and React
<p>
    A simple e-voting system comprising internet voting, 
    to simulate Malaysia's 15th General Election live voting.
    The difficulty for citizens and students to participate a short time frame during election day
    will promote a higher electoral participation and trust in the system overall.
    This e-voting web application will store all of the electronic ballots in cloud database, and reflect the live result immediately to the voter.
</p>
 
 ## Tech stacks
 - Go
 - Gin
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
<p>
    Kindly change the variable _projectDirName_ in file _dotEnvUtil.go_ and _mongoDbConnection.go_ to "" before docker build.
</p>