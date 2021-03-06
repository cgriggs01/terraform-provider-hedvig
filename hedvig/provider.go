/*Copyright 2015 Container Solutions
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Modified by Katrina for Hedvig
*/
package hedvig

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type LoginResponse struct {
	Result struct {
		Datacenters []interface{} `json:"datacenters"`
		DisplayName string        `json:"displayName"`
		Roles       struct {
			Hedvig string `json:"Hedvig"`
		} `json:"roles"`
		Dualdc        bool   `json:"dualdc"`
		SessionID     string `json:"sessionId"`
		UserName      string `json:"userName"`
		PrimaryTenant string `json:"primaryTenant"`
	} `json:"result"`
	RequestID string `json:"requestId"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type HedvigClient struct {
	Username string
	Password string
	Node     string
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{ // Source https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/provider.go#L20-L43
		Schema:        providerSchema(),
		ResourcesMap:  providerResources(),
		ConfigureFunc: providerConfigure,
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"node": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := HedvigClient{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Node:     d.Get("node").(string),
	}

	return &client, nil
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"hedvig_vdisk":  resourceVdisk(),
		"hedvig_lun":    resourceLun(),
		"hedvig_mount":  resourceMount(),
		"hedvig_access": resourceAccess(),
	}
}

func GetSessionId(d *schema.ResourceData, p *HedvigClient) string {
	u := url.URL{}
	u.Host = p.Node
	u.Path = "/rest/"
	u.Scheme = "http"

	q := url.Values{}
	q.Set("request", fmt.Sprintf("{type:Login,category:UserManagement,params:{userName:'%s',password:'%s',cluster:'%s.hedviginc.com'}}",
		p.Username, p.Password, d.Get("cluster").(string)))
	log.Printf("URL: %+v\n", u.String())

	u.RawQuery = q.Encode()

	log.Printf("QUERY: %v\n", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	login := LoginResponse{}
	err = json.Unmarshal(body, &login)

	if err != nil {
		log.Fatalf("Error unmarshalling: %s", err)
	}
	log.Printf("login: %+v", login)

	log.Printf("session: %s", login.Result.SessionID)

	return login.Result.SessionID
}
