pipeline { 
    agent any  

    stages { 
        
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

                                        docker build -t firepit .

                                        docker stop firepit_prod || true
                                        docker rm firepit_prod || true

                                        docker run \
                                            --name firepit_prod \
                                            -d \
                                            -p 127.0.0.1:3003:3000 \
                                            firepit:latest

                                    ''', 
                                    execTimeout: 120000, 
                                    flatten: false, 
                                    makeEmptyDirs: false, 
                            noDefaultExcludes: false, 
                            patternSeparator: '[, ]+', 
                            remoteDirectory: './pipeline_firepit', 
                            remoteDirectorySDF: false, 
                            removePrefix: '', 
                            sourceFiles: 'Dockerfile, src/**'
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
