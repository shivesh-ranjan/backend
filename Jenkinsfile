pipeline {
    agent {
        docker {
	    image 'derekshaw/golang_with_docker:1.0'
	    args '--user root -v /var/run/docker.sock:/var/run/docker.sock'
	}
    }
    stages {
        stage('Setup DB') {
	    steps {
	        sh '''
		    docker run -d -p 5432:5432 -e POSTGRES_DB=auth -e POSTGRES_PASSWORD=secret postgres:17.2
		'''
	    }
	}
        stage('Install golang-migrate') {
	    steps {
	        sh '''
		    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz
		    mv migrate.linux-amd64 /usr/bin/migrate
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
    }
}
