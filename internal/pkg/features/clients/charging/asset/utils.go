package asset

import (
	"time"
)

const (
	ASSET_SITES_PATH                         = "sites"
	ASSET_ANCESTORS_PATH                     = "hierarchies/ancestors"
	DEFAULT_ENERGY_AUTH_TOKEN_PATH           = "tokens"
	DEFAULT_ENERGY_AUTH_POWERGATE_TOKEN_PATH = "energy-auth/powergate_tokens"
	POWERGATE_ASSET_PROXY_SITES_PATH         = "powergate-asset-proxy/sites"
)

type TokenResponse struct {
	Data TokenContainer `json:"data"`
}

type TokenContainer struct {
	Token string `json:"token"`
}

type AssetQuery struct {
	SiteNumber     string `url:"site_number,omitempty"`
	Din            string `url:"din,omitempty"`
	ExternalSiteID string `url:"external_site_id,omitempty"`
	AssetSiteID    string `url:"target_id,omitempty"`
}

type SiteResponse struct {
	SiteInfo SiteInfoContainer `json:"data"`
}

type SiteResponseArray struct {
	SiteInfo []SiteInfoContainer `json:"data"`
}

type AddressContainer struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Zip          string `json:"zip,omitempty"`
	County       string `json:"county,omitempty"`
	Country      string `json:"country"`
	Timezone     string `json:"time_zone"`
}

type BatteryContainer struct {
	TotalNamePlateMaxChargePower    float64   `json:"total_nameplate_max_charge_power"`
	TotalNamePlateMaxDischargePower float64   `json:"total_nameplate_max_discharge_power"`
	TotalNamePlateEnergy            float64   `json:"total_nameplate_energy"`
	Batteries                       []Battery `json:"batteries"`
}

type Battery struct {
	DeviceId                   string  `json:"device_id"`
	DIN                        string  `json:"din"`
	SerialNumber               string  `json:"serial_number"`
	NamePlateMaxChargePower    float64 `json:"nameplate_max_charge_power"`
	NamePlateMaxDischargePower float64 `json:"nameplate_max_discharge_power"`
	NamePlateEnergy            float64 `json:"nameplate_energy"`
}

type GatewayContainer struct {
	TotalGateways int       `json:"total_gateways"`
	Gateways      []Gateway `json:"gateways"`
}

type Gateway struct {
	DeviceID       string    `json:"device_id"`
	LeaderDeviceID string    `json:"leader_device_id,omitempty"`
	DIN            string    `json:"din"`
	SerialNumber   string    `json:"serial_number"`
	IsActive       bool      `json:"is_active"`
	SiteId         string    `json:"site_id"`
	Version        string    `json:"firmware_version"`
	UpdatedAt      time.Time `json:"updated_datetime"`
}

type GeolocationContainer struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Source    string  `json:"source,omitempty"`
}

type Inverter struct {
	DeviceID string `json:"device_id"`
	DIN      string `json:"din"`
	IsActive bool   `json:"is_active"`
	SiteID   string `json:"site_id"`
	Version  string `json:"firmware_version,omitempty"`
}

type Meter struct {
	DeviceID           string `json:"device_id"`
	DeviceEnergySource string `json:"device_energy_source"`
}

type WallConnector struct {
	DeviceID string `json:"device_id"`
	Din      string `json:"din"`
	IsActive bool   `json:"is_active"`
}

type Customer struct {
	MyxxxUserid        int       `json:"myxxx_userid"`
	OwnershipStartdate time.Time `json:"ownership_startdate"`
	OwnershipEnddate   time.Time `json:"ownership_enddate"`
	CreatedAt          time.Time `json:"created_datetime"`
	UpdatedAt          time.Time `json:"updated_datetime"`
}
type Customers []Customer

type Inverters []Inverter

type Meters []Meter

type WallConnectors []WallConnector

// SiteInfoContainer is the direct struct response from asset/sites
type SiteInfoContainer struct {
	SiteId                string     `json:"site_id"`
	SiteNumber            string     `json:"site_number"`
	SiteName              string     `json:"site_name"`
	InstanceIds           []string   `json:"instance_ids"`
	ExternalSiteID        string     `json:"external_site_id,omitempty"`
	MarketType            string     `json:"market_type"` // @warning: may not be reliable for industrial vs. residential
	OperationDate         *time.Time `json:"operation_date,omitempty"`
	CreatedDate           *time.Time `json:"created_datetime,omitempty"`
	*AddressContainer     `json:"address"`
	*BatteryContainer     `json:"battery"`
	*GatewayContainer     `json:"gateway"`
	*GeolocationContainer `json:"geolocation"`
	*Customers            `json:"customers"`
	*Inverters            `json:"inverter"`
	*Meters               `json:"meter"`
	*WallConnectors       `json:"wall_connector"`
}
