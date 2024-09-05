@Library("shared-libraries") _

// def s3FileName = "go-auth"
// def bucketName = "go-auth-bucket"

def appName = "go-auth"

def deployConfig = [
    main: [
        revisionTag: appName,
        revisionLocation: 'go-auth-assets',
        assetsPath: 'app/',
        codeDeployAppName: 'internal-projects',
        codeDeployGroup: appName
    ]
]

def awsCreds = [
    region: 'eu-west-1',
    iamCredId: 'aws-cred-training-center'
]

pipeline {
    agent any
    
    environment {
        currentBranch = "${env.BRANCH_NAME}"
        gitUser = sh(script: 'git log -1 --pretty=format:%ae', returnStdout: true).trim()
        gitSha = sh(script: 'git log -n 1 --pretty=format:"%H"', returnStdout: true).trim()
        imageRegistry = '909544387219.dkr.ecr.eu-west-1.amazonaws.com'
        imageName = "go-auth"
        imageTag = "${imageRegistry}/${imageName}:${gitSha}"
    }
    
    stages{
        stage('Build Docker Image') {
            steps {
                script {
                    buildDockerImage(imageTag: imageTag, buildContext: '.')
                }
            }
        }

        stage('Push to Registry') {
            steps {
                script {
                    pushDockerImage(image: imageTag, registry: imageRegistry, awsCreds: awsCreds)
                }
            }
        }

        stage('Prepare Deploy') {
            steps {
                when {
                    branch 'main'
                }
                script {
                    sh 'mkdir -p app/'
                    sh 'cp docker-compose.yml app/'
                    sh 'cp deployment-scripts/ -r app/'
                    sh 'cp appspec.yml app/'
                    sh 'sed -i "s|image: go-auth-backend:latest|image: $imageTag|g" app/docker-compose.yml'
                    prepareToDeployECR(environment: currentBranch, deploymentConfig: deployConfig, awsCreds: awsCreds)
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    makeDeploymentECR(environment: currentBranch, deploymentConfig: deployConfig, awsCreds: awsCreds)
                }
            }
        }

        stage('Clean Up Build') {
            steps {
                script {
                    sh "docker rmi ${imageTag}"
                    sh 'docker system prune -f'
                }
            }
        }
    }
}

post{
    always{
        script {
            cleanWs()
        }
    }
}
