`nginx.local.conf` is set up for a private, local nginx reverse-proxy that forwards 81 to 80:

```
nginx -p $(pwd) -c nginx.local.conf
```

Nginx will still complain about not being able to open `/var/log/nginx/error.log`, but this is expected and can be ignored.

