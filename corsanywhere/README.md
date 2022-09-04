# CorsAnywhere
for your all cors pain problems

## Overview
This allow all, no config cors middleware that
will solve all your cors problem. Compatible 
with standard `net/http` routers
like `gorilla/mux` or `httprouter`.

## How to use
This is an example on how to use it with
`gorilla/mux` router, but any other router
that compatible with standard `net/http`
will be the same basically:

```go
router := mux.NewRouter()

// any other router setup
...

srv := &http.Server{
    // it's important to use it like this
    // so the middleware will be always executed before the router
    Handler: corsanywhere.CorsAnywhere(router), 
    Addr:    fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.Port),

    // Good practice: enforce timeouts for servers you create!
    WriteTimeout: 10 * time.Second,
    ReadTimeout:  10 * time.Second,
}

log.Fatal(srv.ListenAndServe())
```