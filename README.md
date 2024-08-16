##  create SSL/TLS self signed certificate
```bash
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout $path_to_Key/private.key -out $path_to_cert/cert.crt
```
## create client cert
```bash
openssl req -new -key  $path_to_Key/private.key -out client.csr
```

