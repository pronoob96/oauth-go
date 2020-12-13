# Oauth Service In Go

## Usage

1) Mongodb

```shell
docker run --name mongodb -d -p 27017:27017 mongo
```

2) Redis

```shell
docker run --name redis -d -p 6379:6379 redis
```

3) Finally run the project

```shell
make run
```