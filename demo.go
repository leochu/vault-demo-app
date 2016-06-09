package main

import (
	"fmt"
	"os"
  "encoding/json"

	"github.com/hashicorp/vault/api"
)

func main() {

	vcapStr := `
{
  "vault": [
   {
    "credentials": {
     "root": "secret/ffe43c22-f835-457b-a88e-a87f6ccb8aa6",
     "token": "efdc7986-bc09-45bc-8353-8399eb915ad6",
     "vault": "http://10.244.8.3:8200"
    },
    "label": "vault",
    "name": "vault-instance",
    "plan": "shared",
    "provider": null,
    "syslog_drain_url": null,
    "tags": []
   }
  ]
 }
`

	// os.Setenv("VCAP_SERVICES", vcapStr)


	vcapJson := os.Getenv("VCAP_SERVICES")

	var vcap map[string]interface{}
	err := json.Unmarshal([]byte(vcapJson), &vcap)
	if err != nil {
	  fmt.Println("error unmarshaling", err)
	}

	root := vcap["vault"]

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

	apiClient := client.Logical()

	secret, err := apiClient.Read(path.(string))
	if err != nil {
	  fmt.Println("Error in reading secret", err)
	}

	fmt.Printf("secret is %#v", secret)

	// http.HandleFunc("/", hello)
	// fmt.Println("listening...")
	// err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	// if err != nil {
	// 	panic(err)
	// }



}

// func hello(res http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintln(res, "go, world")
// }
