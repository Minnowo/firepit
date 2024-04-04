pipeline { 
    agent any  

    stages { 

        stage('Build Env File') {

            steps {

                withCredentials([usernamePassword(credentialsId: 'FIREPIT_DB_CREDS', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                    withCredentials([string(credentialsId: 'FIREPIT_JWT', variable: 'JWT_VALUE')]) {

                                writeFile file: './env.sh', text: """#!/bin/sh
JWT_SECRET="${JWT_VALUE}" 
DB_USERNAME="${USERNAME}" 
DB_PASSWORD="${PASSWORD}" 
DB_HOSTNAME="firepit-mariadb" 
DB_NAME="firepit-mariadb" 
DB_PORT="3306" 
"""

                }
                }

            }
        }
        
        stage('Copy and Build Code on Remote') {

            steps {

                sshPublisher(
                    publishers: [
                        sshPublisherDesc(
                            configName: "${SSH_SERVER}", 
                            transfers: [
                                sshTransfer(
                                    cleanRemote: false, 
                                    excludes: '', 
                                    execCommand: '''
                                    
                                        log_dir="$HOME/log/jenkins-ssh"
                                        mkdir -p "$log_dir"
                                
                                        log_file="$log_dir/pipeline_build_firepit.out"
                                        touch "$log_file"
                                
                                        exec 3>&1 4>&2
                                        trap 'exec 2>&4 1>&3' 0 1 2 3
                                        exec 1>"$log_file" 2>&1

                                        cd ~/pipeline_firepit

                                        chmod +x env.sh
                                        . ./env.sh
                                        rm -rf env.sh

                                        docker build -t firepit-go-img .

                                        docker stop firepit_prod || true
                                        docker rm firepit_prod || true

                                        docker run \
                                            -d \
                                            -e DB_USERNAME="$DB_USERNAME" \
                                            -e DB_PASSWORD="$DB_PASSWORD" \
                                            -e DB_NAME=firepit-mariadb \
                                            -e DB_HOSTNAME=firepit-mariadb \
                                            -e JWT_SECRET="$JWT_SECRET" \
                                            -p 127.0.0.1:3003:3000 \
                                            --network firepit \
                                            --name firepit_prod \
                                            firepit-go-img:latest

                                    ''', 
                                    execTimeout: 120000, 
                                    flatten: false, 
                                    makeEmptyDirs: false, 
                            noDefaultExcludes: false, 
                            patternSeparator: '[, ]+', 
                            remoteDirectory: './pipeline_firepit', 
                            remoteDirectorySDF: false, 
                            removePrefix: '', 
                            sourceFiles: 'env.sh, Dockerfile, src/**'
                            )
                    ], 
                    usePromotionTimestamp: false, useWorkspaceInPromotion: false, verbose: false)
                ])
                
                echo "war file copied"
            }
        }
    }
    
        
    post {
        always {
            discordSend(
                        description: currentBuild.result, 
                        enableArtifactsList: false, 
                        footer: '', 
                        image: '', 
                        link: '', 
                        result: currentBuild.result, 
                        scmWebUrl: '', 
                        thumbnail: '', 
                        title: env.JOB_BASE_NAME, 
                        webhookURL: "${DISCORD_WEBHOOK_1}"
                    )
        }
    }
}
