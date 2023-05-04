# Auth API Server for the Test

## Demo Server 

A running demo is on [onesgo.simonxu.net](https://onesgo.simonxu.net/), which is on a single node kubernetes with arm64 arch.

## API Examples

1. To register an new user: 
  `curl --request POST --url https://onesgo.simonxu.net/register
    --header 'Content-Type: application/json'
    --data '{
			"email":      "test@example.com",
			"password":   "testpassword",
			"first_name": "Lingkai",
			"last_name":  "Xu",
    }'`
2. To login: 
  `curl --request POST --url https://onesgo.simonxu.net/login
    --header 'Content-Type: application/json'
    --data '{
      "password": "testpassword",
      "email": "test@example.com"
    }'`
3. To get an user's information:
   `curl --request GET --url https://onesgo.simonxu.net/users/[USEE_ID] 
    --header 'Authorization: bearer [YOUR_JWT_TOKEN]'`
4. To update an user's information:
  `curl --request POST --url https://onesgo.simonxu.net/users/[USER_ID] 
    --header 'Authorization: bearer [YOUR_JWT_TOKEN]' 
    --header 'Content-Type: application/json' 
    --data '{
      "phone": "12345678",
      "first_name": "Simon"
    }'`
5. To delete (or deactive) current user:
  `curl --request DELETE --url http://localhost:8080/users/[USER_ID] 
  --header 'Authorization: bearer [YOUR_JWT_TOKEN]'`

## Comment Convention

I put commments with annotation `// ?` in the source code, including config files, to record my thoughts related to that part.

## Other Thoughts

- Existing tests only prove that the service is available under normal conditions. More tests are needed to cover edge cases and internal util functions. 
- Every user input should be more carefully verified to prevent injection attacks and wrong data from entering the system.
- The output to users should be checked more carefully to prevent leakage of sensitive information and  exposures of attack points, such as internal state of the service.
- Although the data is stored in the database, the container is ephemeral and the data will be destroyed along with the container. I didn't allocate persistent volume to the container, because I think storage should be treated as an independent layer. If it was a formal deployment, I will use a separately deployed database (or managed database service) and object storage service for high availability and performance.
- The limit of current implementation is that it uses a fixed database configuration. The data layer ( or say model layer) should be abstracted, to provide a more generalised interface to service layer ( or say controller layer). Ideally, We can flexibly replace (or inject) the underlying storage.
- I spent too long for setting up my server and trying to use Github Actions to generate images in different architectures.
- I put other thoughts in code comments.