package cpogovsdk

// GovChargingRequestEncryption defines model for GovChargingRequestEncryption.
type GovChargingRequestEncryption struct {
	Data       *string `json:"Data,omitempty"`
	OperatorID *string `json:"OperatorID,omitempty"`
	Seq        *string `json:"Seq,omitempty"`
	Sig        *string `json:"Sig,omitempty"`
	TimeStamp  *string `json:"TimeStamp,omitempty"`
}
