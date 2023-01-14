migrate -source file:migrations -database "mysql://root:1@tcp(localhost:3306)/test?parseTime=true&multiStatements=true" up 1
