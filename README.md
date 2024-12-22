# Cocktails search

List from [IBA](https://iba-world.com/) of all cocktails. The goal of the API is to bring a fine grain search to find coktails

# API

## GET Cocktails

```
curl -X GET .../cocktails?term=booze&notIncluded=booze1
```

`term` are liquors, ingredientes or cocktail name to be *included* in your search.

`notIncluded` does not return cocktails that contain that string