pipeline {
    agent any

    environment {
        REGISTRY = 'register.l3xx.cc'
        IMAGE_NAME = "${REGISTRY}/mint-server"
        GOPROXY = 'https://goproxy.cn,direct'
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '20'))
        disableConcurrentBuilds(abortPrevious: true)
        timeout(time: 30, unit: 'MINUTES')
        timestamps()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_SHORT_SHA = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                    env.GIT_BRANCH_NAME = env.BRANCH_NAME ?: env.GIT_BRANCH?.replaceFirst('origin/', '') ?: 'unknown'
                }
            }
        }

        stage('Verify') {
            steps {
                sh '''
                    set -e
                    echo "branch=${GIT_BRANCH_NAME} commit=${GIT_SHORT_SHA}"
                    go version
                    go build ./...
                    go vet ./...
                    go test ./...
                    if command -v buf >/dev/null 2>&1; then
                        buf lint
                    else
                        echo "buf not installed, skipping buf lint"
                    fi
                    if command -v golangci-lint >/dev/null 2>&1; then
                        golangci-lint run ./...
                    else
                        echo "golangci-lint not installed, skipping"
                    fi
                '''
            }
        }

        stage('Build & Push') {
            when {
                allOf {
                    branch 'master'
                    expression {
                        def commitMessage = sh(script: 'git log -1 --pretty=%B', returnStdout: true).trim()
                        return !commitMessage.contains('[skip ci]')
                    }
                }
            }
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'REGISTRY_CREDENTIALS',
                    usernameVariable: 'REGISTRY_USER',
                    passwordVariable: 'REGISTRY_PASSWORD'
                )]) {
                    sh '''
                        set -e
                        echo "${REGISTRY_PASSWORD}" | docker login "${REGISTRY}" -u "${REGISTRY_USER}" --password-stdin
                        docker build -t "${IMAGE_NAME}:${GIT_SHORT_SHA}" .
                        docker push "${IMAGE_NAME}:${GIT_SHORT_SHA}"
                        docker logout "${REGISTRY}" || true
                    '''
                }
            }
        }

        stage('Update Manifest') {
            when {
                allOf {
                    branch 'master'
                    expression {
                        def commitMessage = sh(script: 'git log -1 --pretty=%B', returnStdout: true).trim()
                        return !commitMessage.contains('[skip ci]')
                    }
                }
            }
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'GIT_PUSH_CREDENTIALS',
                    usernameVariable: 'GIT_USER',
                    passwordVariable: 'GIT_TOKEN'
                )]) {
                    sh '''
                        set -e
                        sed -i "s|image: register.l3xx.cc/mint-server:.*|image: register.l3xx.cc/mint-server:${GIT_SHORT_SHA}|" deploy/k8s/mint-server.yaml
                        git config user.name "Jenkins CI"
                        git config user.email "ci@l3xx.cc"
                        git add deploy/k8s/mint-server.yaml
                        git diff --cached --quiet && echo "manifest unchanged, skip commit" && exit 0
                        git commit -m "ci: update image to ${GIT_SHORT_SHA} [skip ci]"
                        git push "https://${GIT_USER}:${GIT_TOKEN}@github.com/ethereal3x/mint-server.git" HEAD:master
                    '''
                }
            }
        }
    }

    post {
        success {
            echo "mint-server ${GIT_SHORT_SHA} pipeline succeeded"
        }
        failure {
            echo "mint-server ${GIT_SHORT_SHA} pipeline failed"
        }
        always {
            sh 'docker logout register.l3xx.cc 2>/dev/null || true'
        }
    }
}
