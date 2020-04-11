# dechoconf

[![Build Status](https://travis-ci.org/Wuvist/dechoconf.svg?branch=master)](https://travis-ci.org/Wuvist/dechoconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/Wuvist/dechoconf?v=1)](https://goreportcard.com/report/github.com/Wuvist/dechoconf)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)

## Dependency

dechoconf depends on :

* go **1.11** and above
* github.com/BurntSushi/toml v0.3.1
* gopkg.in/yaml.v2 v2.2.8

## Introduction

`dechoconf` propose a convention when defining struct for unmarshalling config data: ```prefix:"decode_prefix"```.

Assuming we have two module dependencies in our application: `db` and `api`, each requires it's own config:

```go
type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type APIConfig struct {
	URL string
}
```

Usually, we will need to define additional config struct for our application that contains above both structs:

```go
type AppConfig struct {
	DB  DBConfig
	API APIConfig
}

func Setup () error {
	var configs AppConfig
	if _, err = toml.DecodeFile("config file path", &configs); err != nil {
		return err
	}

	db.Setup(configs.DB)
	api.Setup(configs.API)
}
```

Assuming our config file is in toml format:

```conf
[DB]
host = "localhost"
port = 3306
username = "root"
password = ""

[api]
url = "https://localhost:8080"
```

It works, but didn't scale well when our dependencies grows. Each time we introduce a new dependency, we need to change both our `AppConfig` struct and `Setup` method.

`Dependency Injection` is a common technic to ease such situation, and `dechoconf` aims to make configuration easier here.

## prefix tag annotion

When using `dechoconf`, we could define dependency's config struct in the following convention:

```go
type DBConfig struct {
	_        string `prefix:"db"`
	Host     string
	Port     int
	Username string
	Password string
}

type APIConfig struct {
	_   string `prefix:"api"`
	URL string
}
```

go doesn't support any kind of annotation on struct level, so we put it to the special field `-`.

With the prefix hint from `-` field annotation, `dechoconf` will be able to perform smart unmarshalling:
```go
var dbConfig DBConfig
var apiConfig APIConfig
if err := dechoconf.DecodeTomlFile(tomlData, &dbConfig, &apiConfig); err != nil {
    return
}
```

When using together with dependency injection framework like [wire](https://github.com/google/wire), the code could be even simpler:

```go
func LoadConfigs() (appConfig app, apiConfig backendAPI, err error) {
	err = dechoconf.DecodeTomlFile("config file path", &appConfig, &apiConfig)
	return
}
```

When adding new dependency, we just need to change `LoadConfigs` method signature, and pass new parameter to `dechoconf.DecodeTomlFile`. New returned config struct will be automactically wired to dependency.
