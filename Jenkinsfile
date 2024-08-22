@Library("shared-libraries") _

def s3FileName = "go-auth"
def bucketName = "go-auth-bucket"

def appName =   "go-auth"

def deployConfig = [
    testing: [
        revisionTag: appName,
        revisionLocation: 'go-auth-assets',
        assetsPath: 'app/',
        codeDeployAppName: 'go-auth-backend',
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
        stage('Build Application Image') {
            when {
                branch 'main'
            }
            steps {
                script {
                    buildDockerImage(imageTag: imageTag, buildContext: '.')
                }
            }
        }

        stage('Push to Registry') {
            when {
                branch 'main'
            }
            steps {
                script {
                    pushDockerImage(image: imageTag, registry: imageRegistry, awsCreds: awsCreds)
                }
            }
        }

        stage('Prepare To Deploy') {
            when {
                branch 'main'
            }
            steps {
                prepareToDeploy(s3FileName: s3FileName, appName: appName, bucketName: bucketName)
            }
        }

        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                script {
                    makeDeploymentECR(environment: currentBranch, deploymentConfig: deployConfig, awsCreds: awsCreds)
                }
            }
        }

        stage('Clean Up Build') {
            when{
                branch 'main'
            }
            steps {
                script {
                    sh "docker rmi ${imageTag}"
                    sh 'docker system prune -f'
                }
            }
        }

        stage('CleanUp WS') {
            steps {
                script {
                    cleanWs()
                }
            }
        }
    }
}
