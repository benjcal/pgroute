# pgroute

Simple library to expose all postgres functions to an endpoint.

Uses `chi` and `gorm`

### usage:

```
db, _ := gorm.Open(...)
r := chi.NewRouter()

r.Mount("/f", pgroute.MountFunctionRoute(db))
```