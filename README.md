# go-proxy

Service written in Go language that allows splitting query destinations to a MySQL databases using query rules.

## QDS - Query destination split

**QDS** has for now two ways of splitting:

- Regex Rule Split
- Hash Rule Split

Both of them allow you to split the query direction but use different ways to find the destiny. Both are compatible with the caching system.

### RRS - Regex Rule Split

#### What you should know

- If many RRS's matches the query then the first in configuration will be used
- RRS **for now** is case-sensitive
- Queries that are checked against the regex are first normalized to make things simpler

#### How to use it

You can write your own regex rule that if matches the query, then it will be used by `go-proxy`. So, let's imagine that you want to redirect all `SELECT` queries made to table `large_table` to **replica** (ID of replica is `R1`).

You can create the RRS with this regex: `^SELECT.*FROM.*large_table.*` (you should spend more time creating this regex, this is a very basic example and as you can see it's far from being ideal):

Now you can add new regex rule to the `config.yml` file:

```yml
name: "REDIRECT ALL SELECT QUERIES MADE TO large_table TO R1"
regex_rule: "SELECT.*FROM.*large_table.*"
target_id: "R1"
```

### HRS - Hash Rule Split

## Configuration

Configuration is currently located in the `config.yml` file, and the structure looks as follows:

```yml
proxy:
  basics: # go-proxy basic configuration
    host: "127.0.0.1" # host on which go-proxy operates
    port: 1234 # port of communication with go-proxy
  access: # MySQL protocol access user to go-proxy  
    user: "user"
    password: ""
  servers: # primary/replica_1 db definition
    - name: "PRIMARY" # name of the db 
      id: "P1" # id of the db (it has to be unique)
      host: "192.168.250.230" # host on which MySQL db operates
      port: 3306 # port of communication with MySQL db
      required: true # if is required then go-proxy won't start up if MySQL db is down
      test_db: "test" # db used for test of communication
      default: true # if set to true then every util that doesn't hit util rule will be redirected to this db
    - name: "REPLICA" 
      id: "R1"
      host: "192.168.250.230"
      port: 3306
  rules:
    - name: "REDIRECT SELECT FOR UPDATE QUERIES TO PRIMARY" # name of the util
      regex_rule: "^SELECT FOR UPDATE.*" # regex rule - regexp definition of rule  
      target_id: "P1" # to which db this rule should direct   
    - name: "REDIRECT SELECT QUERIES TO REPLICA" 
      regex_rule: "^SELECT.*"
      target_id: "R1"
    - name: "SELECT * FROM versions WHERE major=?"
      hash_rule: "3c343df0eb5b1832b1c8443e63340718dae9c8dbaaa43193e3db435d40dffe94" # hash rule - SHA-256 representation of normalized util
      target_id: "R1"
  db_users:
    - target: "P1" # which db has this user 
      user: "root"
      password: "passwd"
    - target: "R1"
      user: "root"
      password: "passwd"
```

## Installation

Compile the project:

```shell
go build
chmod +x go-proxy
mv go-proxy /usr/local/bin
```

Add go-proxy user:

```shell
adduser \
   --system \
   --shell /bin/bash \
   --gecos 'Go Proxy user' \
   --group \
   --disabled-password \
   --home /home/go-proxy \
   go-proxy
```

Add required directory structure:

```shell
mkdir /etc/go-proxy
chown root:go-proxy /etc/go-proxy
chmod 770 /etc/go-proxy
```

Create your own config file and name it `config.yml` and then move to the `/etc/go-proxy` directory.

Set directory to read only:

```shell
chmod 750 /etc/go-proxy
chmod 640 /etc/go-proxy/config.yml
```

Add service using the `contrib\systemd\go-proxy.service`.
