package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Node represents the API to a single FReD Node
type Node struct {
	URL string
}

// NewNode creates a new Node with the specified url (shuld have format: http://%s:%d/v%d/)
func NewNode(url string) (node *Node) {
	node = &Node{URL: url}
	return
}

// CreateKeygroup creates a new keygroup with the node. The Response should be empty if everything is correct
func (n *Node) CreateKeygroup(kgname string, expectedStatusCode int, expectEmptyResponse bool) (responseBody map[string]string) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Create Keygroup for group %s; expecting %d", kgname, expectedStatusCode)
	responseBody = n.sendPost("keygroup/"+kgname, nil, expectedStatusCode)
	if expectEmptyResponse && (responseBody != nil) {
		log.Warn().Str("node", n.URL).Msgf("Create Keygroup expected an empty response but got %#v", responseBody)
	}
	return
}

// DeleteKeygroup deletes the specified keygroup
func (n *Node) DeleteKeygroup(kgname string, expectedStatusCode int, expectEmptyResponse bool) (responseBody map[string]string) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete Keygroup for group %s; expecting %d", kgname, expectedStatusCode)
	responseBody = n.sendDelete("keygroup/"+kgname, nil, expectedStatusCode)
	if expectEmptyResponse && (responseBody != nil) {
		log.Warn().Str("node", n.URL).Msgf("Delete Keygroup expected an empty response but got %#v", responseBody)
	}
	return
}

// PutItem puts a key-value pair into a (already created) keygroup
func (n *Node) PutItem(kgname, item string, data string, expectedStatusCode int, expectEmptyResponse bool) (responseBody map[string]string) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Put for Item %s in KG %s; expecting %d", item, kgname, expectedStatusCode)
	reqBody := map[string]string{
		"data": data,
	}
	responseBody = n.sendPut(fmt.Sprintf("keygroup/%s/data/%s", kgname, item), reqBody, expectedStatusCode)
	if expectEmptyResponse && (responseBody != nil) {
		log.Warn().Str("node", n.URL).Msgf("PutItem expected an empty response but got %#v", responseBody)
	}
	return
}

// GetItem returns the stored item
func (n *Node) GetItem(kgname, item string, expectedStatusCode int, expectEmptyResponse bool) (responseBody map[string]string){
	log.Debug().Str("node", n.URL).Msgf("Sending a Get for Item %s in KG %s; expecting %d", item, kgname, expectedStatusCode)
	responseBody = n.sendGet(fmt.Sprintf("keygroup/%s/data/%s", kgname, item), expectedStatusCode)
	if expectEmptyResponse && (responseBody != nil) {
		log.Warn().Str("node", n.URL).Msgf("GetItem expected an empty response but got %#v", responseBody)
	} else if !expectEmptyResponse && responseBody == nil {
		log.Warn().Str("node", n.URL).Msg("GetItem expected a response but got nothing")
	}
	return
}

// DeleteItem deletes the item from the keygroup
func (n *Node) DeleteItem(kgname, item string, expectedStatusCode int, expectEmptyResponse bool) (responseBody map[string]string){
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete for Item %s in KG %s; expecting %d", item, kgname, expectedStatusCode)
	responseBody = n.sendDelete(fmt.Sprintf("keygroup/%s/data/%s", kgname, item), nil, expectedStatusCode)
	if expectEmptyResponse && (responseBody != nil) {
		log.Warn().Str("node", n.URL).Msgf("DeleteItem expected an empty response but got %#v", responseBody)
	} else if !expectEmptyResponse && responseBody == nil {
		log.Warn().Str("node", n.URL).Msg("DeleteItem expected a response but got nothing")
	}
	return
}

// RegisterReplica registers a new replica with this node
func (n *Node) RegisterReplica(nodeID, nodeIP string, nodePort int, expectedStatusCode int, expectEmptyResponse bool) (responseBody map[string]string){
	log.Debug().Str("node", n.URL).Msgf("Registering Replica %s ; expecting %d", nodeID, expectedStatusCode)
	// type Message struct {
	// 	Nodes []struct{
	// 		ID string
	// 		IP string
	// 		Port int
	// 	}
	// }
	// message := Message{Nodes: make([]struct {
	// 	ID   string
	// 	IP   string
	// 	Port int
	// }, 1)}
	// message.Nodes[0] = struct {
	// 	ID   string
	// 	IP   string
	// 	Port int
	// }{ID: nodeID, IP: nodeIP, Port: nodePort}
	//
	// json, _ := json.Marshal(message)
	json := []byte(fmt.Sprintf(`{"nodes":[{"id":"%s","addr":"%s","port":%d}]}`, nodeID, nodeIP, nodePort))
	responseBody = n.sendPost("replica", json, expectedStatusCode)
	if expectEmptyResponse && (responseBody != nil && len(responseBody) != 0) {
		log.Warn().Str("node", n.URL).Msgf("RegisterReplica expected an empty response but got %#v with len %d", responseBody, len(responseBody))
	} else if !expectEmptyResponse && (responseBody == nil || len(responseBody) == 0) {
		log.Warn().Str("node", n.URL).Msg("RegisterReplica expected a response but got nothing")
	}
	return
}

// GetAllReplica returns a list of all replica that this node has stored
func (n *Node) GetAllReplica(expectedStatusCode int, expectEmptyResponse bool) (responseBody []string){
	log.Debug().Str("node", n.URL).Msgf("Sending a Get for all Replicas; expecting %d", expectedStatusCode)
	responseBody = n.sendGetResponseArray("replica", expectedStatusCode)

	if expectEmptyResponse && (responseBody != nil && len(responseBody) != 0) {
		log.Warn().Str("node", n.URL).Msgf("GetAllReplica expected an empty response but got %#v with len %d", responseBody, len(responseBody))
	} else if !expectEmptyResponse && (responseBody == nil || len(responseBody) == 0) {
		log.Warn().Str("node", n.URL).Msg("GetAllReplica expected a response but got nothing")
	}
	return
}

func (n *Node) sendGet(path string, expectedStatusCode int) (responseBody map[string]string) {
	resp, err := http.Get(n.URL + path)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal().Str("node", n.URL).Err(err).Msg("SendGet got HTTP error")
		return nil
	}
	if resp.StatusCode != expectedStatusCode {
		log.Error().Str("node", n.URL).Msgf("SendGet got wrong HTTP Status Code Response. Expected: %d, Got: %d", expectedStatusCode, resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	// return if the response is empty
	if buf.Len() == 0 {
		return
	}
	// Load buf into responseBody
	err = json.Unmarshal(buf.Bytes(), &responseBody)

	if err != nil {
		log.Fatal().Str("node", n.URL).Err(err).Str("json", buf.String()).Msg("sendGet got a response with invalid json")
	}

	return
}

func (n *Node) sendGetResponseArray(path string, expectedStatusCode int) (responseBody []string) {
	resp, err := http.Get(n.URL + path)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal().Str("node", n.URL).Err(err).Msg("SendGet got HTTP error")
		return nil
	}
	if resp.StatusCode != expectedStatusCode {
		log.Error().Str("node", n.URL).Msgf("SendGet got wrong HTTP Status Code Response. Expected: %d, Got: %d", expectedStatusCode, resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	// return if the response is empty
	if buf.Len() == 0 {
		return
	}
	// Load buf into responseBody
	err = json.Unmarshal(buf.Bytes(), &responseBody)

	if err != nil {
		log.Fatal().Str("node", n.URL).Err(err).Str("json", buf.String()).Msg("sendGet got a response with invalid json")
	}

	return
}

func (n *Node) sendPut(path string, data map[string]string, expectedStatusCode int) (responseBody map[string]string) {
	client := &http.Client{}
	jsonBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPut, n.URL+path, bytes.NewBuffer(jsonBytes))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal().Str("node", n.URL).Err(err).Msg("sendPut got HTTP error")
		return nil
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if resp.StatusCode != expectedStatusCode {
		log.Error().Str("node", n.URL).Msgf("SendPut got wrong HTTP Status Code Response. Expected: %d, Got: %d. Response Body: %s", expectedStatusCode, resp.StatusCode, buf.String())
	} else {
		err = json.Unmarshal(buf.Bytes(), &responseBody)
		if err != nil {
			log.Error().Str("node", n.URL).Err(err).Str("json", buf.String()).Msg("sendPut got a response with invalid json")
		}
	}
	return
}

func (n *Node) sendPostDataString(path, data string, expectedStatusCode int) (responseBody map[string]string) {
	return n.sendPost(path, []byte(data), expectedStatusCode)
}

func (n *Node) sendPost(path string, data []byte, expectedStatusCode int) (responseBody map[string]string) {
	var jsonBytes []byte
	if data != nil {
		var err error
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			log.Fatal().Str("node", n.URL).Msgf("Cannot marshal JSON: %v", data)
		}
	}
	var resp *http.Response
	var err error
	if jsonBytes != nil {
		resp, err = http.Post(n.URL+path, "application/json", bytes.NewBuffer(jsonBytes))
	} else {
		resp, err = http.Post(n.URL+path, "", nil)
	}
	if err != nil {
		log.Fatal().Str("node", n.URL).Err(err).Msg("sendPost got HTTP error")
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedStatusCode {
		log.Error().Str("node", n.URL).Msgf("SendPost got wrong HTTP Status Code Response. Expected: %d, Got: %d.", expectedStatusCode, resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// Load buf into responseBody
	err = json.Unmarshal(buf.Bytes(), &responseBody)
	if err != nil && err.Error() != "unexpected end of JSON input" {
		log.Fatal().Str("node", n.URL).Err(err).Str("json", buf.String()).Msg("sendPost got a response with invalid json")
	}
	return
}

func (n *Node) sendDelete(path string, data map[string]string, expectedStatusCode int) (responseBody map[string]string) {
	var jsonBytes []byte
	if data != nil {
		var err error
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			log.Fatal().Str("node", n.URL).Msgf("Cannot marshal JSON: %v", data)
		}
	}
	var resp *http.Response
	var err error
	client := &http.Client{}
	if jsonBytes != nil {
		jsonBytes, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodDelete, n.URL+path, bytes.NewBuffer(jsonBytes))
		resp, err = client.Do(req)
	} else {
		req, _ := http.NewRequest(http.MethodDelete, n.URL+path, nil)
		resp, err = client.Do(req)
	}
	if err != nil {
		log.Fatal().Str("node", n.URL).Err(err).Msg("sendDelete got HTTP error")
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedStatusCode {
		log.Error().Str("node", n.URL).Msgf("SendDelete got wrong HTTP Status Code Response. Expected: %d, Got: %d.", expectedStatusCode, resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// Load buf into responseBody
	err = json.Unmarshal(buf.Bytes(), &responseBody)
	if err != nil && err.Error() != "unexpected end of JSON input" {
		log.Fatal().Str("node", n.URL).Err(err).Str("json", buf.String()).Msg("sendDelete got a response with invalid json")
	}
	return
}
