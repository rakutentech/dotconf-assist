# dotconf-assist

*Web GUI and REST API Designed for Splunk Configuration Management*

## Overview

**Admin**

- **Home** management
- **Forwarder** management
- **Server Class** management
- **Inputs** management
- **Apps** management
- **Peployment** management
- **Users** management
- **Announcements** management
- **Splunk Hosts** management
- **Unit Price** management

**User**

- **Home** management
- **Forwarder** management
- **Server Class** management
- **Inputs** management
- **Apps** management
- **Peployment** management
- **Usage** management


## Code Structure

```
|-- dotconf_assist_project
    |-- pkg
    |-- bin
    |-- src
        |-- github.com
            |-- rakutentech
                |-- dotconf-assist (Current Repository)
                    |--e2e
                    |--Godeps
                    |--node_modules (not tracked by git)
                    |--src
                        ...
        |-- golang.org
        |-- gopkg.in
```

## Dependencies

- go >1.9.1
- npm >3.10.10
- node >6.11.4
- Angular >4.4.5
- @angular/cli: >1.4.7
- typescript >2.3.4

## Preparation and Installation

**Mac**

- **node, npm** <https://nodejs.org/en/>
- **go** <https://golang.org/dl/>
- **@angular/cli** `npm install -g @angular/cli`

```
# tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz
$ vi ~/.bash_profile
PATH=$PATH:$HOME/.local/bin:$HOME/bin:/usr/local/go/bin
$ source ~/.bash_profile
```

**CentOS 6.x** 

- **node, npm** <https://nodejs.org/en/download/package-manager/#debian-and-ubuntu-based-linux-distributions>
- **go** <https://golang.org/dl/>
- **@angular/cli** `npm install -g @angular/cli`

```
# curl --silent --location https://rpm.nodesource.com/setup_6.x | bash -
# yum -y install nodejs

# wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz

# tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz
$ vi ~/.bash_profile
PATH=$PATH:$HOME/.local/bin:$HOME/bin:/usr/local/go/bin
$ source ~/.bash_profile
```

**CentOS 7.x** 

- **node, npm**
- **go** <https://golang.org/dl/>
- **@angular/cli** `npm install -g @angular/cli`

```
# yum install epel-release
# yum install nodejs npm

# wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz

# tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz
$ vi ~/.bash_profile
PATH=$PATH:$HOME/.local/bin:$HOME/bin:/usr/local/go/bin
$ source ~/.bash_profile
```

**Ubuntu**

- **node, npm** <http://qiita.com/seibe/items/36cef7df85fe2cefa3ea>
- **go** <https://golang.org/dl/>
- **@angular/cli** `npm install -g @angular/cli`

```
# tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz
$ vi ~/.bash_profile
PATH=$PATH:$HOME/.local/bin:$HOME/bin:/usr/local/go/bin
$ source ~/.bash_profile
```

**Version confirmation**

```
$ ng -v
@angular/cli: 1.4.7
@angular/animations: 4.4.5
...
typescript: 2.3.4

$ node -v
v6.11.4

$ npm -v
3.10.10

$ go version
go version go1.9.1
```

## Development & Deployment

### 1. (If required) Preparing environment

```
$ ssh <your server> or localhost
$ vi ~/.bash_profile
export GOROOT=/usr/local/go
export GOPATH=<dotconf_assist_project>
export PATH=$PATH:/usr/local/go/bin/:$GOPATH/bin
$ source ~/.bash_profile
```

### 2. Getting code ready

```
(first time)
$ mkdir -p <dotconf_assist_project>/src/github.com/rakutentech/
$ cd <dotconf_assist_project>/src/github.com/rakutentech/
$ git clone https://github.com/rakutentech/dotconf-assist.git

(except for first time)
$ cd <dotconf_assist_project>/src/github.com/rakutentech/dotconf-assist
$ git pull origin master //git pull origin develop
```

### 3. (If required) Editing `src/backend/settings/conf.json`
```
$ cp src/backend/settings/conf.example.json src/backend/settings/conf.json
$ vi src/backend/settings/conf.json
```

### 4. (If required) Editing `src/app/configuration.ts`
```
$ cp src/app/configuration.example.ts src/app/configuration.ts
$ vi src/app/configuration.ts
```

### 5. Installiing packages

**API**

```
(For the first time)
$ go get github.com/tools/godep (golang package manager)
$ $GOPATH/bin/godep restore (install packages, 'godep save' must have been issued before git push)
$ npm install -g gulp
$ npm install -g http-server

(If new go package imported)
$ godep restore (install packages, 'godep save' must have been issued before git push)

(If go file changed)
$ go install (dotconf-assist execuation file will be generated in $GOPATH/bin/)
$ ls -l <dotconf_assist_project>/bin/dotconf-assist
```

**GUI**

```
(If new npm package imported)
$ npm install

(If ts file changed)
$ ng build (build ts files to js files and scss file to css file into dist folder)
$ ls -l dist
```

### 6. Enabling API and GUI (Development)

**API**

```
(if table schema changed)
$ go run server.go --with-migration

(regular start)
$ go run server.go
```

**GUI**

```
$ ng serve

(If SSL enabled)
$ ng serve --ssl -true --ssl-key "../cert/key.pem" --ssl-cert "../cert/cert.pem"
```

### 7. Enableing API and GUI (Production)

```
$ sudo su -
# cd /etc/systemd/system
```
**API** 

```
# vi dotconf-assist-api.service
```

```
[Unit]
Description=dotconf-assist api server daemon
[Service]
 
WorkingDirectory=/opt/gitrepo/dotconf_assist_project/src/github.com/rakutentech/dotconf-assist
ExecStart=/opt/gitrepo/dotconf_assist_project/bin/dotconf-assist
Restart=on-failure
RestartSec=1s
[Install]
WantedBy=multi-user.target
```
**GUI**

```
# vi dotconf-assist-gui.service
```

```
[Unit]
Description=dotconf-assist gui server daemon
 
[Service]
WorkingDirectory=/opt/gitrepo/dotconf_assist_project/src/github.com/rakutentech/dotconf-assist/dist
ExecStart=/bin/http-server -S -C ../../cert/cert.pem -K ../../cert/key.pem -p 443 -c-1
Restart=on-failure
RestartSec=1s
[Install]
WantedBy=multi-user.target
```

```
# systemctl daemon-reload (If the service file above changed)

(api)
# systemctl restart dotconf-assist-api (if `go install` executed)

(gui)
# systemctl restart dotconf-assist-gui (if ts files changed)
```

### 8. License

MIT License

### 9. Author

Peng Yang