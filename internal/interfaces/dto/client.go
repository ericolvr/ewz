package dto

type ClientRequest struct {
	CustomerName  string  `json:"cliente_nome"      binding:"required"`
	CustomerEmail string  `json:"cliente_email"     binding:"required,email"`
	RequestType   string  `json:"tipo_solicitacao"  binding:"required"`
	AssetValue    float64 `json:"valor_patrimonio"  binding:"required,gt=0"`
}

type ClientResponse struct {
	ID            int64   `json:"id"`
	CustomerName  string  `json:"cliente_nome"`
	CustomerEmail string  `json:"cliente_email"`
	RequestType   string  `json:"tipo_solicitacao"`
	AssetValue    float64 `json:"valor_patrimonio"`
	Status        string  `json:"status"`
}
