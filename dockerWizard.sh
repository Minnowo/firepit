#!/bin/bash

#* Small Wizard that "magically" handles deployment actions with Docker

# Function to delete deployment and network
delete_deployment() {
    echo "Stopping and removing containers..."
    docker stop firepit-mariadb
    docker rm firepit-mariadb

    docker stop firepit-go
    docker rm firepit-go

    echo "Removing network..."
    docker network rm firepit-network

    echo "Deletion complete."
}

# Function to deploy MariaDB and Firepit Backend
deploy() {
    echo "Deploying MariaDB..."
    docker run --detach --name "firepit-mariadb" --env MARIADB_ROOT_PASSWORD="root" --env MARIADB_DATABASE="firepit-mariadb" -p 127.0.0.1:3306:3306 mariadb:latest

    echo "Sleeping, to avoid MariaDB Issue"
    sleep 10

    echo "Creating network and deploying Firepit Backend..."
    docker network create firepit-network
    docker network connect firepit-network firepit-mariadb
    docker run -d --network firepit-network -e JWT_SECRET=51ef7b24b93de21487d852651ac30300 -p 3000:3000 --name firepit-go firepit-go-img

    echo "Deployment complete."
}

# Function to update Firepit Backend image
update_image() {
    git pull
    sleep 1
    echo "Updating Firepit Backend image..."
    docker build -t firepit-go-img .
    echo "Update complete."
}

# Main script logic based on command line argument
case "$1" in
    d|delete|rm)
        delete_deployment
        ;;
    deploy|dp)
        deploy
        ;;
    rebuild|rb)
        delete_deployment
        update_image
        deploy
        ;;
    *)
        echo "Usage: $0 {d|delete|rm|run|build|u|rebuild}"
        ;;
esac
