pipeline { 
    agent any  

     tools {
            maven 'maven3'
            jdk 'OpenJDK17'
    }
    
    stages { 
        
        stage('Build War File') { 
            steps { 
                
                sh "mvn --version"
                sh "mvn clean package"
                
                echo 'War file built' 
            }
        }
        
        stage('Remote Copy War File') {
            steps {
                dir("target"){
                    sshPublisher(publishers: [
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


                                        make_and_copy() {
                                            num="$1"

                                            auto_deploy="$HOME/volumes/glassfish/autodeploy_${num}"

                                            mkdir -p "${auto_deploy}"

                                            chmod 770 "${auto_deploy}"

                                            cp -f $HOME/warbuilds/firepit.war "$auto_deploy"
                                        }

                                        make_and_copy "a"
                                        make_and_copy "b"


                                    ''', 
                                    execTimeout: 120000, 
                                    flatten: false, 
                                    makeEmptyDirs: false, 
                            noDefaultExcludes: false, patternSeparator: '[, ]+', remoteDirectory: './warbuilds', remoteDirectorySDF: false, 
                            removePrefix: '', sourceFiles: 'firepit.war')
                        ], 
                        usePromotionTimestamp: false, useWorkspaceInPromotion: false, verbose: false)
                    ])
                }
                
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
