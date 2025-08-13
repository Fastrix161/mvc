# MVC assignment- "Urban Tadka"
The MVC assignment is a food ordering system that is built in Go using MySQL as the database, in the MVC pattern.
## Installation Guide
1. **Clone repository:-**

```bash
git clone git@github.com:Fastrix161/mvc.git 
cd mvc
```
2. **Download dependencies:-**<br>
```bash
go mod download
```
3. **Configure .env file:-**<br>
Create a ```.env``` file similar to the ```.env.sample``` to set up Environmental variables.
4. **Create database:-**<br>
Create a database in MySQL, named *restaurant*

```mysql
Create Database restaurant;
```

5. **Run Migrations:-**<br>
Run the SQL migrations from ```database/migration``` directory.
<br>

- Set up the database, by running
```bash
migrate -path database/migration -database "mysql://username:secretkey@tcp(localhost:5432)/database_name up
```
in the Terminal.

- To delete/reset the database, run 
```bash
migrate -path database/migration -database "mysql://username:password@tcp(localhost:port)/database_name down
```
<br>

>Remember to change username, secretkey, localhost, port to the values of your .env file for MySQL !!

6. **Run the server:-**
Finally run the servire using `Makefile`. Simply run 
```bash
make run 
``` 
in the Terminal.

Congarts!! You are good to test the server.
