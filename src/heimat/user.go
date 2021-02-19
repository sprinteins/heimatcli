package heimat

// User _
type User struct {
	ID                 int                `json:"id"`
	FirstName          string             `json:"firstName"`
	LastName           string             `json:"lastName"`
	BirthDate          string             `json:"birthdate"`
	Email              string             `json:"email"`
	BusinessPhone      string             `json:"businessPhone"`
	PrivatePhone       string             `json:"privatePhone"`
	AboutMe            string             `json:"aboutMe"`
	JoiningDate        string             `json:"joiningDate"`
	Location           Location           `json:"location"`
	DomainServices     []DomainService    `json:"domainServices"`
	Domains            []Domain           `json:"domains"`
	PSL                PSL                `json:"peopleSuccessLead"`
	IsPSL              bool               `json:"isPeopleSuccessLead"`
	Image              string             `json:"image"`
	imageExtension     string             `json:"imageExtension"`
	LeadRoles          []LeadRole         `json:"leadRoles"`
	CSLProps           []CSLProp          `json:"customerSuccessLeadProperties"`
	EmploymentDuration EmploymentDuration `json:"employmentDuration"`
}

// Location _
type Location struct {
	ID   int    `json:"id"`
	City string `json:"city"`
}

// PSL _
type PSL struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// LeadRole _
type LeadRole struct {
	ID string `json:"id"`
}

// CSLProp _
type CSLProp struct {
	ID   int    `json:"id"`
	Name string `json:"string"`
}

// EmploymentDuration _
type EmploymentDuration struct {
	DurationType string `json:"durationType"`
	DurationUnit int    `json:"durationUnit"`
}
