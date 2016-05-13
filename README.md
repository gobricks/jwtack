# jwtack
JWTAck – jwt acknowledgements as package, coroutine or microservice

##Usage:
1) as package:
```
import jwtackpckg "github.com/gobricks/jwtack/src"
jwtack := jwtackpckg.NewService()
```

2) as coroutine:
```
import jwtack "github.com/gobricks/jwtack/src"
jwtack.Run(jwtack.Config{Port:8001})
```

3) as microservice:
```
(export PORT=36701; ./jwtack)
```
CreateToken:
```
$ curl -H 'X-Jwtack-Key:12345' -d '{"payload":{"asdf":1,"qwer":"wwertwert"}}' http://localhost:36701/api/v1/token;
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhc2RmIjoxLCJxd2VyIjoid3dlcnR3ZXJ0In0.didDNTbAieHZd6QezQuJw46qSBjuMR39bZe6eJuMe1c"}
```

ParseToken:
```
$ curl -H 'X-Jwtack-Key:12345' http://localhost:36701/api/v1/token/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhc2RmIjoxLCJxd2VyIjoid3dlcnR3ZXJ0In0.didDNTbAieHZd6QezQuJw46qSBjuMR39bZe6eJuMe1c;
{"payload":{"asdf":1,"qwer":"wwertwert"}}
```
