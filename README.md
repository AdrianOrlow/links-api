![License](https://img.shields.io/github/license/AdrianOrlow/links-api)
[![CodeFactor](https://www.codefactor.io/repository/github/adrianorlow/links-api/badge)](https://www.codefactor.io/repository/github/adrianorlow/links-api)
![Go](https://img.shields.io/github/go-mod/go-version/AdrianOrlow/links-api)
# Links API

My personal link shortener API. Made with Go, GORM, Gorilla Mux and MySQL.

[Links frontend](https://github.com/AdrianOrlow/links)

![thumbnail](https://user-images.githubusercontent.com/10941338/71741122-5076fc00-2e5e-11ea-9ee4-253e3ae56654.png)

## Getting started

Firstly, rename `.env.sample` to `.env` and fill all the fields with your data.
It should me mentioned that `ADMIN_GMAIL_ADDRESSES` is array of Google accounts email addresses (separated with commas) which
can login to the system and perform link creation operation.

Once you filled the config you can run the server via

```
go run main.go
```

If you want to build the package, run

```
go build
```

## Deployment (Dokku)

Create the app container

```
dokku apps:create app_name
```

create the mysql database container

```
dokku mysql:create app_name-db
```

link database to the container

```
dokku mysql:link app_name-db app_name
```

set all the environment variables
   
```
dokku config:set PORT=5000 HASH_ID_SALT= ...
```

add Dokku remote repository

```
git remote add dokku dokku@server_ip:app_name
```

and finally push code to the repo

```
git push dokku master
```

## License

[MIT](https://choosealicense.com/licenses/mit/)