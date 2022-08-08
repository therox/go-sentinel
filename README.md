# go-sentinel
Sentinel client

# Using
Instantiate new client with copernicus credentials and timeout on http requests (0 for no timeout):
```Go
 client := sentinel.NewClient(user, password, 60*time.Minute)
 ```

Define query
```Go
searchParameters := sentinel.SearchParameters{
    Platforms: []sentinel.Platform{sentinel.PlanformSentinel2},
    ProductTypes:  []string{"S2MSI2A", "S2MS2Ap"},
    BeginDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
}
```

And do the query
```Go 
res, err := client.Query(searchParameters)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Total entries: ", res.Feed.TotalResults) += res.Feed.TotalResults
```

Check, if product is online
```Go
isOnline, err := client.IsOnline(entry.ID)
if err != nil {
    fmt.Println(err)
}
```

Dounload product
```Go
    err := client.Download(entry.GetID(), "/tmp")
    if err != nil {
        fmt.Println(err)
    }
```