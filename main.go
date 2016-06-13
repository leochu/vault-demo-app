package main

import (
	"fmt"
	"os"
  "encoding/json"
  "net/http"

	"github.com/hashicorp/vault/api"
	"github.com/nu7hatch/gouuid"
)

func main() {
  // apiClient, path := getClient()

	// go func() {
		// for {
		  // if apiClient == nil {
  //       fmt.Println("Cannot connect to Vault")
			// } else {
				// fmt.Println("Reading secret at ", path)
				// secret, err := apiClient.Read(path)
				// if err != nil {
					// fmt.Println("Error in reading secret", err)
				// }

				// fmt.Println("secret is ", secret)
			// }
			// time.Sleep(5 * time.Second)
		// }
  // }()

	http.HandleFunc("/secrets", secrets)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}


func secrets(res http.ResponseWriter, req *http.Request) {
  if req.Method == "GET" {
    read(res, req)
	} else if req.Method == "POST" {
    write(res, req)
	}
}

func read(res http.ResponseWriter, req *http.Request) {
  apiClient, path := getClient()

  var respStr string

	if apiClient == nil {
		respStr = "Not connected to Vault"
	} else {
		fmt.Println("Reading secret at ", path)
		secret, err := apiClient.Read(path)
		if err != nil {
			fmt.Println("Error in reading secret", err)
		}

    if secret == nil {
      respStr = "{}"
		} else {
			j, _ := json.Marshal(secret.Data)
			respStr = string(j)
		}
	}

	fmt.Fprintln(res, respStr)
}

func write(res http.ResponseWriter, req *http.Request) {
  apiClient, path := getClient()

	var respStr string

	if apiClient == nil {
		respStr = "Not connected to Vault"
	} else {
		secret, err := apiClient.Read(path)
		if err != nil {
			fmt.Println("Error in reading secret", err)
		}

		fmt.Println("writing secret at ", path)

    u, _ := uuid.NewV4()
    tenantId := "Tenant_" + u.String()
    u2, _ := uuid.NewV4()
    secretId := "Secret_" + u2.String()

		var data map[string]interface{}
		if secret == nil {
      data = make(map[string]interface{})
		} else {
			data = secret.Data
		}

    data[tenantId] = secretId

		_, err = apiClient.Write(path, data)
		if err != nil {
			fmt.Println("Error in writting secret", err)
		}

		respStr = fmt.Sprintf("secret written")
	}

	fmt.Fprintln(res, respStr)
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
