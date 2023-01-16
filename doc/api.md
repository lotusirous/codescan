# API document

> We can utilize Swagger, but this document is more straightforward to write and maintain compared to the large JSON file required by the Swagger documentation. It should also be consistent with the `server/server.go` file for each API.

**NOTE**

The server handles error codes (400, 404, 500) with the following format

```json
{
    "message" : "server error message goes here"
}
```

# Repository

## List all repositories

`GET http://localhost:8080/api/repos`

Returned a list of repositories

```
[
    {
        "id": 3,
        "name": "vuln",
        "link": "https://github.com/lotusirous/vuln",
        "created": 1673836676,
        "updated": 1673836676
    }
]
```

## Create a repository

`POST http://localhost:8080/api/repos`

Example request body:

```
{
    "name": "vuln-example",
    "link": "https://github.com/lotusirous/vuln"
}
```

## Update a repository

`PUT http://localhost:8080/api/repos/{id}`

Example request body:

```json
{
    "name": "vuln-example"
}
```

## Delete a repository

`DELETE http://localhost:8080/api/repos/{id}`

Return a `NO CONTENT` http status code if success

# Scans

## List all scans

`GET http://localhost:8080/api/scans`

```
[
    {
        "id": 2,
        "repository": 3,
        "status": "Success",
        "enqueuedAt": 1673836959,
        "startedAt": 1673836959,
        "finishedAt": 0
    },
    {
        "id": 3,
        "repository": 3,
        "status": "Success",
        "enqueuedAt": 1673837449,
        "startedAt": 1673837449,
        "finishedAt": 1673837450
    }
]
```

## Scan a repository

`POST http://localhost:8080/api/scans
`

```
{
 "repoID" : 3
}
```

## Get the scan result

`GET http://localhost:8080/api/scans/2`

```
{
    "id": 2,
    "repoName": "vuln",
    "repoURL": "https://github.com/lotusirous/vuln",
    "status": "Success",
    "enqueuedAt": 1673836959,
    "startedAt": 1673836959,
    "finishedAt": 0,
    "findings": [
        {
            "type": "sast",
            "ruleId": "G402",
            "location": {
                "path": "/src/main.py",
                "positions": {
                    "begin": {
                        "line": 1
                    }
                }
            },
            "metadata": {
                "description": "Leak the cryptography keys",
                "severity": "HIGH"
            }
        },
        {
            "type": "sast",
            "ruleId": "G402",
            "location": {
                "path": "/src/main.py",
                "positions": {
                    "begin": {
                        "line": 2
                    }
                }
            },
            "metadata": {
                "description": "Leak the cryptography keys",
                "severity": "HIGH"
            }
        }
    ]
}
```
