# dechoconf

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

## Dependency

dechoconf depends on :

* github.com/BurntSushi/toml

## Todo

* [ ] Add version
* [ ] Setup CI
* [ ] Confirm go version
* [X] Support prefix chain
* [X] Add yaml support
* [X] Prefix validation
* [ ] Support prefix customization
