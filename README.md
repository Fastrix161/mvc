# MVC assignment- "Urban Tadka"
The MVC assignment is a food ordering system that is built in Go using MySQL as the database, in the MVC pattern. For the app, two different ways of setup are provided(local and docker):- 
## Installation Guide
### **Clone repository:-**

```bash
git clone git@github.com:Fastrix161/mvc.git 
cd mvc
```

# I. Local Setup
1. **Configure .env.local file:-**<br>
Create a ```.env.local``` file similar to the ```.env.local.sample``` to set up Environmental variables.

2. **Create database:-**<br>
Create a database in MySQL, named *restaurant*

    ```mysql
    Create Database restaurant;
    ```

3. **Run Migrations:-**<br>
- *Create migrations* <br>

    Run:- 
    ```bash 
    make migrate-up
    ```
    in terminal to set up the migrations.

- *Delete migrations*<br>
    Run:- 
    ```bash
    make migrate-down
    ```
    in terminal to delete the migrations.
<br>

4. **Run the server:-** <br>
    Finally run the server/app using
    ```bash
    make run 
    ``` 
    in the Terminal.

Congarts!! You are good to test the server on local.

# II. Docker Setup
1. **Configure .env.docker file:-**<br>
Create a ```.env.docker``` file similar to the ```.env.docker.sample``` to set up docker Environmental variables.

2. **Run Docker:-**<br>
- *Create docker*<br>
    Run docker using the using 
    ```bash
    make docker-up
    ```
    in your Terminal.

- *Stop docker*<br>
    Stop/delete docker using the using 
    ```bash
    make docker-down
    ```
    in your Terminal.

And there you go!! You are good to test the server on docker.

