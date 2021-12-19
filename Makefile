build:
	go build \
		-ldflags "-X main.buildCommit=`git rev-parse --short HEAD` \
		-X main.buildTime=`date " +%Y-%m-%dT%H:%M:%S:%Z:00"`" \
		-o app

maria:
	docker run -p 127.0.0.1:3306:3306 --name todo-mariadb \
	-e MARIADB_ROOT_PASSWORD=my-secret-pw -e MARIADB_DATABASE=myapp -d mariadb:latest

image:
	docker build -t todo:test -f Dockerfile .

container:
	docker run -p 8081:8081 --env-file local.env --link todo-mariadb:db \
	--name myapp todo:test