pipeline {
    agent any
    tools { 
	go"go1.23.4"
    }
	//   environment {
	//SONARQUBE_URL = 'http://45.198.13.210:9000'
	//SONARQUBE_PROJECT_KEY = 'backend'
	//   }
    stages {
        stage('Setup DB') {
	    steps {
	        sh '''
		    docker run -d -p 5432:5432 -e POSTGRES_DB=auth -e POSTGRES_PASSWORD=secret --name postgres postgres:17.2
		'''
	    }
	}
        stage('Install golang-migrate') {
	    steps {
	        sh '''
		    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz
		    sudo mv migrate.linux-amd64 /usr/bin/migrate
	            which migrate
		'''
	    }
	}
        stage('Run Migrations') {
	    steps {
		sh '''
		    cd auth
		    make migrateup
		'''
	    }
	}
        stage('Testing') {
	    steps {
		sh '''
		    cd auth
		    make test
		'''
	    }
	}
	stage('Cleaning Up') {
	    steps {
		sh '''
		    docker stop postgres
		    docker rm postgres
		'''
	    }
	}
	//stage('SonarQube Analysis') {
	//    steps {
	//	withCredentials([string(credentialsId: 'sonarqube', variable: 'SONARQUBE_TOKEN')]) {
	//	    sh '''
	//                       /opt/sonar-scanner/bin/sonar-scanner \
	//                         -Dsonar.projectKey=$SONARQUBE_PROJECT_KEY \
	//                         -Dsonar.sources=. \
	//                         -Dsonar.host.url=$SONARQUBE_URL \
	//                         -Dsonar.login=${SONARQUBE_TOKEN}
	//                   '''
	//	}
	//    }
	//}
	stage('SonarQube Analysis') {
            steps {
                script {
                    def scannerHome = tool 'SonarScanner'
                    withSonarQubeEnv() {
                        sh "${scannerHome}/bin/sonar-scanner"
                    }
                }
            }
        }
    }
}
