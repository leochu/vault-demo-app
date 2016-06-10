package main

import (
	"fmt"
	"os"
  "encoding/json"
  "net/http"
  "time"

	"github.com/hashicorp/vault/api"
)

func main() {

	// vcapStr := `
// {
  // "vault": [
   // {
    // "credentials": {
     // "root": "secret/ffe43c22-f835-457b-a88e-a87f6ccb8aa6",
     // "token": "efdc7986-bc09-45bc-8353-8399eb915ad6",
     // "vault": "http://10.244.8.3:8200"
    // },
    // "label": "vault",
    // "name": "vault-instance",
    // "plan": "shared",
    // "provider": null,
    // "syslog_drain_url": null,
    // "tags": []
   // }
  // ]
 // }
// `

	// os.Setenv("VCAP_SERVICES", vcapStr)

  apiClient, path := getClient()

	go func() {
		for {
		  if apiClient == nil {
        fmt.Println("Cannot connect to Vault")
			} else {
				fmt.Println("Reading secret at ", path)
				secret, err := apiClient.Read(path)
				if err != nil {
					fmt.Println("Error in reading secret", err)
				}

				fmt.Println("secret is ", secret)
			}
			time.Sleep(5 * time.Second)
		}
  }()

	http.HandleFunc("/", hello)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "go, world")
}


func getClient() (*api.Logical, string){
	vcapJson := os.Getenv("VCAP_SERVICES")

	if vcapJson == "" {
		return nil, ""
	}

	var vcap map[string]interface{}
	err := json.Unmarshal([]byte(vcapJson), &vcap)
	if err != nil {
	  fmt.Println("error unmarshaling", err)
	}

	root := vcap["vault"]

  if root == nil {
		return nil, ""
	}

	element := root.([]interface{})[0]
	credentials := element.(map[string]interface{})["credentials"]
	path := credentials.(map[string]interface{})["root"]
	token := credentials.(map[string]interface{})["token"]
	url := credentials.(map[string]interface{})["vault"]

	os.Setenv("VAULT_ADDR", url.(string))
  // os.Setenv("VAULT_TOKEN", "efdc7986-bc09-45bc-8353-8399eb915ad6")

  config := api.DefaultConfig()

	err = config.ReadEnvironment()
	if err != nil {
	  fmt.Println("error reading environment variables", err)
	}

	client, err := api.NewClient(config)
	if err != nil {
	  fmt.Println("error creating client", err)
	}

	client.SetToken(token.(string))

	return client.Logical(), path.(string)
}
