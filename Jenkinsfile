pipeline {
    agent any
    tools { 
	go"go1.23.4"
    }
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
	stage('SonarQube Analysis') {
            steps {
                script {
                    def scannerHome = tool 'sonarqube'
                    withSonarQubeEnv() {
                        sh "${scannerHome}/bin/sonar-scanner"
                    }
                }
            }
        }
	stage('Building Docker Image'){
	    steps {
		sh '''
		    docker build -t derekshaw/gatewaymicro:${GIT_COMMIT} ./auth/
		'''
	    }
	}
	stage('Publish Image to Dockerhub') {
	    steps {
		script {
		    withDockerRegistry(credentialsId: 'docker-hub-credentials', toolName: 'docker') {
			sh 'docker push derekshaw/gatewaymicro:$GIT_COMMIT'
		    }
		}
	    }
	}
	stage('Update and Commit Image Tag for ArgoCD') {
	    steps {
		sh 'git clone -b main https://github.com/shivesh-ranjan/backend-ops.git'
		sh '''
		    git checkout main
		    sed -i "s#derekshaw/gatewaymicro.*#derekshaw/gatewaymicro:$GIT_COMMIT#g" app-services.yml
		    git commit -am "updated docker image"
		    git push -u origin main
		'''
	    }
	}
    }
}
