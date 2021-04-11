package api


type APIResponse struct {
    Status string   `json:"status"`
    Error string    `json:"error"`
}


type CostEstimateResponse struct {
    Status string       `json:"status"`
    GasPrice float64    `json:"gasPrice"`
    EthCost float64     `json:"ethCost"`
    Error string        `json:"error"`
}
