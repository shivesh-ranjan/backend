pipeline {
    agent any
    tools { go"go1.23.4" }
    stages {
        stage('Setup DB') {
	    steps {
		script {
		    img = 'postgres:17.2'
		    docker.image("${img}").run("-d -p 5432:5432 -e POSTGRES_DB=auth -e POSTGRES_PASSWORD=secret --name postgres")
	    	}
	    }
	}
	// Installed in the docker image running Jenkins and docker
	//       stage('Install golang-migrate') {
	//    steps {
	//        sh '''
	//	    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz
	//	    mv migrate.linux-amd64 /usr/bin/migrate
	//            which migrate
	//	'''
	//    }
	//}
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
    }
}
