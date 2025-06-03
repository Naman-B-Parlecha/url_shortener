pipeline {
    agent any
    environment {
        AWS_ID = 'add_id'
        AWS_REGION = 'add_region'
        BROKER_REPO = 'url-short/broker-service'
        SHORTENER_REPO = 'url-short/shortener-service'
        ANALYTICS_REPO = 'url-short/analytics-service'
        REDIRECT_REPO = 'url-short/redirect-service'
        IMAGE_TAG = "latest"
        AWS_CREDENTIALS = credentials('MinusX-Creds')
        ECS_SERVICE = 'nig-url-short-task-service-pvd7qeb2'
        ECS_CLUSTER = 'url-short-service'
        TASK_FAMILY = 'url-short-task'
    }
    stages {
        stage('Checkout Repository') {
            steps {
                echo "Checking out code"
                git url: 'https://github.com/Naman-B-Parlecha/url_shortener.git', branch: 'main', credentialsId: 'minusx-id'
            }
        }
        stage('Build Services') {
            parallel {
                stage('Build Broker') {
                    steps {
                        dir('broker-service') {
                            script {
                                dockerBrokerImage = docker.build("${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${BROKER_REPO}:${IMAGE_TAG}")
                            }
                        }
                    }
                }
                stage('Build Shortener') {
                    steps {
                        dir('url-shortener-service') {
                            script {
                                dockerShortenerImage = docker.build("${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${SHORTENER_REPO}:${IMAGE_TAG}")
                            }
                        }
                    }
                }
                stage('Build Redirect') {
                    steps {
                        dir('redirect-service') {
                            script {
                                dockerRedirectImage = docker.build("${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${REDIRECT_REPO}:${IMAGE_TAG}")
                            }
                        }
                    }
                }
                stage('Build Analytics') {
                    steps {
                        dir('analytics-service') {
                            script {
                                dockerAnalyticsImage = docker.build("${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ANALYTICS_REPO}:${IMAGE_TAG}")
                            }
                        }
                    }
                }
            }
        }
        stage('Push To ECR') {
            steps {
                script {
                    sh "aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
                    dockerBrokerImage.push()
                    dockerRedirectImage.push()
                    dockerAnalyticsImage.push()
                    dockerShortenerImage.push()
                }
            }
        }
        stage('Force ECS Deployment') {
            steps {
                script {
                    sh "aws ecs update-service --cluster ${ECS_CLUSTER} --service ${ECS_SERVICE} --force-new-deployment --region ${AWS_REGION}"
                }
            }
        }
    }
    post {
        always {
            sh """
            docker rmi ${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${BROKER_REPO}:${IMAGE_TAG} || true
            docker rmi ${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${SHORTENER_REPO}:${IMAGE_TAG} || true
            docker rmi ${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${REDIRECT_REPO}:${IMAGE_TAG} || true
            docker rmi ${AWS_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ANALYTICS_REPO}:${IMAGE_TAG} || true
            """
        }
    }
}
