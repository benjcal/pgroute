# pgroute

Simple library to expose all postgres pubic functions to HTTP.

Uses `chi` and `gorm`

### usage:


```golang
db, _ := gorm.Open(...)
r := chi.NewRouter()

r.Mount("/f", pgroute.MountFunctionRoute(db))
```


```sql
CREATE FUNCTION add_user(first_name TEXT, last_name TEXT, age INT) ...
```


```sh
curl -X
```
