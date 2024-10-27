
# Firepit

Stateful Rooms!


## Deploy with Docker

Clone the repo and create a deploy.sh file like so:
```sh
#!/bin/sh

VERSION="0.0.1"
IMAGE="firepit:${VERSION}"
OUTPUT="firepit.tar"
SSH_SERVER="your-ssh-server-config-name"

(cd src && make build-js)

docker build -t "${IMAGE}" .
docker save -o "${OUTPUT}" "${IMAGE}"

scp "$OUTPUT" "$SSH_SERVER:/path/on/server/firepit.tar"

TERM=xterm-256color ssh "$SSH_SERVER"

```

