maria:
	docker run -p 127.0.0.1:3306:3306 --name todo-mariadb \
	-e MARIADB_ROOT_PASSWORD=my-secret-pw -e MARIADB_DATABASE=myapp -d mariadb:latest