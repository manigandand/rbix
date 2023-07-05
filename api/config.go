package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	// EnvLoclDocker is the local docker environment
	EnvLoclDocker = "dev-docker"

	// EnvLocalK8s is the local k8s environment
	EnvLocalK8s = "dev-k8s"
)

var (
	// Env is the environment
	Env = EnvLoclDocker

	// Port is the port number to run the server
	Port = "8080"

	// SqrxWSLoadbalncerHost is the host name of the load balancer
	SqrxWSLoadbalncerHost = "localhost:8081"

	// SqrxRbiImage is the docker image name of sqrx-rbi
	SqrxRbiImage = "manigandanjeff/sqrx-rbi:latest"
)

// Initialize initializes all the env variables for this package.
// This function should be called only once!
//
// Example:
//
//	LoadFromJSON("config.json", "dash.json")
//	addNewEnvEntry("HOST", "localhost", &host)
//	addNewEnvEntry("PORT", "8080", &port)
//	...
//	load()
func Initialize(files ...string) {
	if initDone {
		panic("config initialization done already")
	}
	LoadFromJSON(files...)

	addNewEnvEntry("ENV", &Env, Env)
	addNewEnvEntry("PORT", &Port, Port)
	addNewEnvEntry("SQRX_WS_LOADBALANCER_HOST", &SqrxWSLoadbalncerHost, SqrxWSLoadbalncerHost)
	addNewEnvEntry("SQRX_RBI_IMAGE", &SqrxRbiImage, SqrxRbiImage)

	// load all the env variables. Must be called at the end.
	load()
	log.Println("Inited config...👍")
}

// cfg holds the configuration values for the application.
var cfg = make(map[string]string)
var initDone = false

// LoadFromJSON is used to load the config key value pairs from a json file
func LoadFromJSON(files ...string) {
	if initDone {
		panic("config initialization done already")
	}
	for _, file := range files {
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("[WARN] unable to load from config file: %s\n", err)
			return
		}

		err = json.Unmarshal(bs, &cfg)
		if err != nil {
			log.Printf("[WARN] unable to load from config file: %s\n", err)
			return
		}
	}
}

// runtimeValue contains data for a single env value which had to be fetched
// from the environment during runtime.
type runtimeValue struct {
	VarPtr  *string // Pointer to the variable which will holds the value
	Default string  // the default value if the env variable is not set
}

// Vars contains all the maps of env key to the pointer memory
var envVars = make(map[string]runtimeValue)

// addNewEnvEntry is used to add a new env entry into the Vars map
//
//	envKey     - What key this variable is defined as in env
//	varPtr     - Pointer to the variable which will holds the value
//	defaultVal - What should be the default value if it is not defined
func addNewEnvEntry(envkey string, varPtr *string, defaultVal string) {
	if initDone {
		panic("Initialization done already")
	}
	envVars[envkey] = runtimeValue{
		Default: defaultVal,
		VarPtr:  varPtr,
	}
}

// Load initializes all the variables for this package.
//
// This function reads/loads the environment variables and sets the values
// to the variables from the environment and/or configMaps.
//
// This function is should call after setting the source and other
// required variables to be read Initialize(addNewEnvEntry and LoadFromJSON).
func load() {
	if initDone {
		panic("config initialization done already")
	}

	for k, v := range envVars {
		mustEnv(k, v.VarPtr, v.Default)
	}
	initDone = true
}

// mustEnv reads the env variable with the name 'key' and assigns the value
// in 'varPtr'.
//
// this checks the env variable from either configMaps or env.
func mustEnv(key string, varPtr *string, defaultVal string) {
	// check the env variable in configMaps, use if not null
	val, ok := cfg[key]
	if ok && val != "" {
		*varPtr = val
		return
	}

	// else, check in the os.env
	if *varPtr = os.Getenv(key); *varPtr == "" {
		*varPtr = defaultVal
		log.Printf("%s env variable not set, using default value.\n", key)
	}
}
