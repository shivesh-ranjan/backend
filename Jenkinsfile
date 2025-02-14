pipeline {
    agent any
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
		    mv migrate.linux-amd64 ~/bin/migrate
    		    export PATH="/var/lib/jenkins/bin:$PATH"
	            which migrate
		'''
	    }
	}
        stage('Run Migrations') {
	    steps {
		sh '''
		    cd auth
		    export DB_SOURCE=postgresql://postgres:secret@host.docker.internal:5432/auth?sslmode=disable
    		    export PATH="~/bin:$PATH"
		    make migrateup
		'''
	    }
	}
        stage('Testing') {
	    steps {
		sh '''
		    cd auth
		    export DB_SOURCE=postgresql://postgres:secret@host.docker.internal:5432/auth?sslmode=disable
    		    export PATH="~/bin:$PATH"
		    make test
		'''
	    }
	}
    }
}
