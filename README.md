# Task  | Instagram Backend API

Name: Susruta Das
Registration Number: 19BEC1366
University: Vellore Institute of Technology, Chennai


#### Create a user
```http
  POST /users
  
  payload = {
    id
    name
    email
    password
  }
```

#### Get users By ID 

```http
  GET /users/{id}
```
#### Create a post

```http
  POST /posts
  
  payload = {
    id
    caption
    image
    timestamp
  }
```
#### Get post By ID 

```http
  GET /posts/{id}
```
#### Get all posts by a user

```http
  GET /posts/users/{id}
  ```

## Testing Instructions

## 1. Run Main file

```
go run main.go
```
## 2. Simulate requests using curl

* To create an User
```
curl -X POST -H 'content-type: application/json' --data '{"id":"2","name":"Appointy","email":"Appointy@gmail.com","password":"123xyz"}' http://localhost:8080/users
```

* To get User details from User ID
```
curl http://localhost:8080/users/1
```
* To create a new post using User ID
```
curl -X POST -H 'content-type: application/json' --data '{"id":"2","caption":"Beautiful walls","image":"walls.png","timestamp":"2021-10-09T13:49:51.141Z"}' http://localhost:8080/posts
```

* To Get a post using id
```
curl http://localhost:8080/posts/1
```
* To list all posts of user
```
curl http://localhost:8080/posts
```
