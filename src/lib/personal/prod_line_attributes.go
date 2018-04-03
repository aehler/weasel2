package personal

type ProdLineAttributes struct {
	FacilityME float64 `weaselform:"numeric" formLabel:"Facility PE bonus, %"`
	FacilityTE float64 `weaselform:"numeric" formLabel:"Facility TE bonus, %"`
	FacilityMTax float64 `weaselform:"numeric" formLabel:"Facility manufacturing tax, %"`
	FacilityRTax float64 `weaselform:"numeric" formLabel:"Facility lab tax, %"`
	ManufacturingSlots uint `weaselform:"numeric" formLabel:"My manufacturing slot count"`
	LabSlots uint `weaselform:"numeric" formLabel:"My lab slot count"`
	Investments float64 `weaselform:"numeric" formLabel:"I can invest"`
	TimeBasis string `weaselform:"select" formLabel:"Time basis"`
	ManufSystem uint `weaselform:"autocomplete" formLabel:"I manufacture in"`
	ResearchSystem uint `weaselform:"autocomplete" formLabel:"I invent in"`
}