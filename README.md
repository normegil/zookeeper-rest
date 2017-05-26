# Zookeeper-REST

REST Server used to manage a zookeeper node.

## Installation

`go get github.com/normegil/zookeeper-rest`

## Requirements

  * [MongoDB](https://www.mongodb.com/)

## Usage

Just run Zookeeper-rest. For command line options, use:

`zookeeper-rest --help`

## Configuration
You can configure the server using a config file (default is $HOME/.zookeeper-rest.toml)in [TOML](https://github.com/toml-lang/toml) format, environment variables and/or by directly specifying values to command line options. 

All the environment variables are prefixed with 'ZK\_REST', and are named from the file keys. '.' become '\_' and all letters are UPPER CASED. As an example, the server port is configured with 'ZK\_REST\_SERVER\_PORT'

| Environment Variables | File              | Default Value  | Description                                                |
|-----------------------|-------------------|----------------|------------------------------------------------------------|
| SERVER\_PORT          | server.port       | 8080           | Port on which the rest server will listen.                 |
| LOG\_ROTATION         | logging.rotation  | 7              | Number of days between log file rotation.                  |
| LOG\_DIRECTORY        | logging.directory | /tmp           | Directory where the files will be stored                   | 
| ZOOKEEPER\_ADDRESS    | zookeeper.address | 127.0.0.1      | Address of the Zookeeper server                            |
| VERBOSE               | logging.verbose   | false          | Verbose mode (Debug messages and location for all messages |
| MONGO\_ADDRESS        | mongo.address     | 127.0.0.1      | Address of the mongo DB                                    |
| MONGO\_PORT           | mongo.port        | 27017          | Port of the mongo DB                                       |
| MONGO\_DATABASE       | mongo.database    | zookeeper-rest | Database where logs & Users will be stored                 |
| MONGO\_USER           | mongo.user        |                | User to connect to MongoDB                                 |
| MONGO\_PASS           | mongo.pass        |                | Password to connect to MongoDB                             |

### Authentication
Only basic authentication is supported for now, more will be added later. I strongly advise running the server over HTTPS to ensure the encryption of User/Password.

As there is no API as of yet to get/create/modify/delete users, all operations on them need to be done directly on the MongoDB. On a database of your choices (Specified in the configuration of the REST server), create a collection "users". Add the user needed to connect to the server with the following format:

```
{
  name: 'test',
  password: 'test'
}
```

### Logging
Logging is performed using [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus). If `stdin` is a terminal, then you will have a colored output. 

Additionaly, there is a log file that will be stored according to the LOG\_DIRECTORY folder. The messages will be stored as JSON, and there is a rotation of the log files depending on LOG\_ROTATION.

Those messages are also stored in the 'log' collection in the MongoDB.

### API & Endpoints
The root for all rest services is '/rest'. You can access the following endpoints:
  * node

#### Node
##### Get node
Perform a GET query on

`http://zookeeper-rest:8080/rest/node/:nodeID`

where '/:nodeID' should be the the ID of the request node. To access the root, leave this parameter empty (`http://zookeeper-rest:8080/rest/node`).

It will return a message looking like:
```
{
  "ID": "58bd0797-706c-487e-9cf3-20d6506ccfec",
  "URL": "http://zookeeper-rest:8080/rest/node/58bd0797-706c-487e-9cf3-20d6506ccfec",
  "Path": "/",
  "Content": "{\"nom\": \"updated1\"}",
  "Childs": {
    "/test": "http://zookeeper-rest:8080/rest/node/543d1832-c3f3-48bd-9173-7ef6e0214bbd",
    "/zk-rest": "http://zookeeper-rest:8080/rest/node/1317130b-12b4-4a94-9782-0b55a78170ea",
    "/zookeeper": "http://zookeeper-rest:8080/rest/node/e913fd3e-d787-491d-b69c-9720096044dc"
  }
}
```
The query informations are:
  * *ID:* Id of the node
  * *URL:*  Address to query the node
  * *Path:* Path of the node
  * *Content:* Content of the node
  * *Childs:* A map of association between paths and IDs of the childs of the current node.

##### Create node
In this context, I considered that PUT and POST are one and the same method (even if they're actually not). Perform a PUT/POST on

`http://zookeeper-rest:8080/rest/node`

to create a node. The body should look like:

```
{
	"path": "/test/subtest4",
	"content": "Hello World !"
}
```

##### Update node
In this context, I considered that PUT and POST are one and the same method (even if they're actually not). Perform a PUT/POST on

`http://zookeeper-rest:8080/rest/node/:nodeID`

to update the node. The body should look like:

```
{
	"content": "Hello updated World !"
}
```

##### Delete node
Perform a DELETE on 

`http://zookeeper-rest:8080/rest/node/:nodeID`
