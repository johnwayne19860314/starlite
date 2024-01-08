package cpogovsdk

// GovChargingResponseEncryption defines model for GovChargingResponseEncryption.
type GovChargingResponseEncryption struct {
	Data *string `json:"Data,omitempty"`
	Msg  *string `json:"Msg,omitempty"`
	Ret  *int32  `json:"Ret,omitempty"`
	Sig  *string `json:"Sig,omitempty"`
}
