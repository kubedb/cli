Mount a file to /srv/pgadmin/secrets/.env with secrets

```
email=<email>
password=<password>
```
to configure the first user.

An alternative option is to pass email and password as command line arguments.

```sh
docker run -d -P -it --name=pgadmin-M2MxYWZh appscode/pgadmin:4-1.0_beta1 email password
```

Example:
```sh
docker run -d -P -it --name=pgadmin-M2MxYWZh appscode/pgadmin:4-1.0_beta1 demo@demo.com demo
```
