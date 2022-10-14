up:
	docker-compose up -d --build

down:
	docker-compose down

migrateup:
	docker run -v ${PWD}/migrations:/migrations --network user-auth_default migrate/migrate -path=/migrations/ -database "mysql://root:password@tcp(mysql_user:3306)/user" up $(numb)
migratedown:
	docker run -v ${PWD}/migrations:/migrations --network user-auth_default migrate/migrate -path=/migrations/ -database "mysql://root:password@tcp(mysql_user:3306)/user" down $(numb)