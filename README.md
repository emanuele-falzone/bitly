# bitly

## Usage
```
docker-compose up --detach
```

## GRPC

The GRPC server is available at `localhost:6060`.

The GRPC client is available at [http://localhost:8080](http://localhost:8080).

### Documentation
The GRPC client also fucntions as documentation explointing GRPC server reflection. However, in the Actions section you will find the documentation as `proto.html`.

## HTTP
The HTTP server is availble at `http://localhost:7070`.

The Swagger UI is available at [http://localhost:7070/swagger](http://localhost:7070/swagger).

### Documentation
The Swagger UI functions both as a client and documentation. Be aware the it will automacilly follow redirections and you will probably get an error due to CORS.
Copy and paste the redirection link in the browser and enjoy.
However, in the Actions section you will find the documentation as `swagger.yaml` and `swagger.json`.