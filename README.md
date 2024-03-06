# Firepit
Websocket-based API being re-written in Go, to allow creation of Stateful Rooms!


## Resources

Since I am learning go and have like 10 hours experience

- https://programmingpercy.tech/blog/mastering-websockets-with-go/

- https://www.whichdev.com/building-a-multi-room-chat-application-with-websockets-in-go-and-vue-js-part-2/


## Use our Shell Wizard for Docker Actions 🧙‍♂️

Run DB & Backend with Docker Network: `./dockerWizard.sh deploy` *(`dp` works aswell)*

Delete Deploys & Network: `./dockerWizard.sh delete` *(or `d` works aswell)*

**Recommended** Delete, Rebuild Image & Re-Deploy: `./dockerWizard.sh rebuild` *(`rb` works aswell)*

*Once running the recommended cm, check: **http://localhost:3000/quote***

---

## Manual with Docker

Build Image: `docker build -t firepit-go-img .`

Run as Container: `docker run -p 3000:3000 firepit-go-img`

Delete Entire Deployment & Network

```shell
docker stop firepit-mariadb
docker rm firepit-mariadb

docker stop firepit-go
docker rm firepit-go

docker network rm firepit-network
```

1. Deploy MariaDB:

`docker run --detach --name "firepit-mariadb" --env MARIADB_ROOT_PASSWORD="root" --env MARIADB_DATABASE="firepit-mariadb" -p 127.0.0.1:3306:3306 mariadb:latest`

2. Deploy Firepit Backend & Network Setup

```shell
docker network create firepit-network
docker network connect firepit-network firepit-mariadb
docker run -d --network firepit-network -e JWT_SECRET=51ef7b24b93de21487d852651ac30300 -p 3000:3000 --name firepit-go firepit-go-img
```

---

### Common Issue / Fix:

```go
func getDBConf() *database.DBConfig {
	return &database.DBConfig{
		Username:     "root",
		Password:     "root",
		Hostname:     "firepit-mariadb",
		Port:         3306,
		DatabaseName: "firepit-mariadb",
	}
}
```
In *src/cmd/backend/main.go* , make sure in the **getDBConf** method's return, Hostname is the Name of the Docker Container, using the Docker Setup Wizard, it should just be `firepit-mariadb` !