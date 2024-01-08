package cpopartnersdk

// PartnerTCECRequest defines model for PartnerTCECRequest.
type PartnerTCECRequest struct {
	Data       string `json:"Data"`
	OperatorID string `json:"OperatorID"`
	Seq        string `json:"Seq"`
	Sig        string `json:"Sig"`
	TimeStamp  string `json:"TimeStamp"`
}
