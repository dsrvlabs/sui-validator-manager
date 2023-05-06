package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/dsrvlabs/sui-validator-manager/types"
)

type requestMessage struct {
	ID      string `json:"id,omitempty"`
	JsonRPC string `json:"jsonrpc,omitempty"`
	Method  string `json:"method,omitempty"`
	Params  []any  `json:"params,omitempty"`
}

func newRequestMessage(method string, params []any) requestMessage {
	msg := requestMessage{
		ID:      "1",
		JsonRPC: "2.0",
		Method:  method,
		Params:  []any{},
	}

	if params != nil {
		msg.Params = params
	}

	return msg
}

type SuiClient interface {
	LatestCheckpointSequenceNumber() (*big.Int, error)
	Checkpoint(no *big.Int) (*types.Checkpoint, error)
	LatestSuiSystemState() (*types.SuiSystemState, error)

	GetStakes(address string) (types.StakeInfoList, error)
}

type rpcClient struct {
	endpoints  []string
	httpClient http.Client
}

func (c *rpcClient) LatestCheckpointSequenceNumber() (*big.Int, error) {
	msg := newRequestMessage("sui_getLatestCheckpointSequenceNumber", nil)
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoints[0], bytes.NewReader(rawMsg))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := struct {
		ID      string `json:"id"`
		Jsonrpc string `json:"jsonrpc"`
		Result  string `json:"result"`
	}{}

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return nil, err
	}

	ret, ok := new(big.Int).SetString(respData.Result, 10)
	if !ok {
		return nil, errors.New("convert response failed")
	}

	return ret, nil
}

func (c *rpcClient) Checkpoint(no *big.Int) (*types.Checkpoint, error) {
	msg := newRequestMessage("sui_getCheckpoint", []any{no.String()})
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoints[0], bytes.NewReader(rawMsg))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respMsg := struct {
		Jsonrpc string           `json:"jsonrpc"`
		Result  types.Checkpoint `json:"result"`
	}{}

	err = json.Unmarshal(respBody, &respMsg)
	if err != nil {
		return nil, err
	}

	return &respMsg.Result, nil
}

func (c *rpcClient) LatestSuiSystemState() (*types.SuiSystemState, error) {
	msg := newRequestMessage("suix_getLatestSuiSystemState", []any{})
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoints[0], bytes.NewReader(rawMsg))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respMsg := struct {
		Jsonrpc string               `json:"jsonrpc"`
		Result  types.SuiSystemState `json:"result"`
	}{}

	err = json.Unmarshal(respBody, &respMsg)
	if err != nil {
		return nil, err
	}

	return &respMsg.Result, nil
}

func (c *rpcClient) GetStakes(address string) (types.StakeInfoList, error) {
	msg := newRequestMessage("suix_getStakes", []any{address})
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoints[0], bytes.NewReader(rawMsg))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respMsg := struct {
		Jsonrpc string            `json:"jsonrpc"`
		Id      string            `json:"id"`
		Result  []types.StakeInfo `json:"result"`
	}{}

	err = json.Unmarshal(respBody, &respMsg)

	return respMsg.Result, nil
}

func NewClient(endpoints []string) SuiClient {
	return &rpcClient{
		endpoints: endpoints,
		httpClient: http.Client{
			Timeout: time.Second * 15,
		},
	}
}
