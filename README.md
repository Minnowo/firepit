# Firepit
Websocket-based API written in Java, to allow creation of Stateful Rooms!

---

## System Requirements for building of .war 

- JDK 17 [oracle.com/jdk17](https://www.oracle.com/java/technologies/javase/jdk17-archive-downloads.html)
- Maven [maven.apache.org](https://maven.apache.org/download.cgi)

**Windows Installation of Maven**: [Random Tutorial, Setting ENVIRONMENT vars](https://phoenixnap.com/kb/install-maven-windows)

Verify your maven installation: `maven --version`

---
## Building of Project (.war):
`mvn clean package` 

Expected output *(Indicating Successful Build)*:
```
[INFO] Building war: ...\firepit\target\firepit.war
```
---

## GlassFish7 Containerization & Running

#### Build Glassfish7 Container from *Dockerfile*:
`docker build -t glassfish7 .`

1. Make sure you've got `password_1.txt` & `password_2.txt` in the same directory as the Dockerfile
2. Make sure you've got a copy of Glassfish (`glassfish.zip` in same dir. as Dockerfile), if you don't, download it here: [DOWNLOAD GLASSFISH ZIP](https://download.eclipse.org/ee4j/glassfish/glassfish-7.0.4.zip) and rename it to: `glassfish.zip`

#### Run API Application: 
`docker run --rm -p 4848:4848 -p 8080:8080 glassfish7:latest`

- **NOTE:** `firepit.war` file is copied into Container in *Dockerfile*, so the run command can be simpler!