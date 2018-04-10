package models

import (
	"time"
)

type IDStruct struct {
	ID int `json:"id"`
}

type NameStruct struct {
	Name string `json:"name"`
}

type User struct {
	ID                int       `json:"id"`
	Admin             bool      `json:"admin"`
	UserName          string    `json:"user_name" sql:"unique_index"`
	Email             string    `json:"email"`
	EmailForEmergency string    `json:"email_for_emergency"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastLoginAt       time.Time `json:"last_login_at"`
	Password          string    `json:"password"`
	GroupName         string    `json:"group_name"`
	AppTeamName       string    `json:"app_team_name"`
	// ServiceID         int       `json:"service_id"`
	ServiceID string `json:"service_id"`
	Status    string `json:"status"`
}

type SplunkUser struct {
	ID            int       `json:"id"`
	UserName      string    `json:"user_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	RpaasUserName string    `json:"rpaas_user_name"`
	Env           string    `json:"env"`
	Memo          string    `json:"memo"`
	SearchHead    string    `json:"search_head"`
	Email         string    `json:"email"`
	PositionIDs   string    `json:"position_ids"`
	Password      string    `gorm:"-"`
	// Status            string    `json:"status"`
}

type Announcement struct {
	ID        int       `json:"id"`
	Content   string    `json:"content" gorm:"size:1000"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SplunkHost struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" sql:"unique_index"`
	Role      string    `json:"role"`
	Env       string    `json:"env"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Forwarder struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserName  string    `json:"user_name"`
	Env       string    `json:"env"`
	Share     string    `json:"share"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ServerClass struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	UserName     string    `json:"user_name"`
	Env          string    `json:"env"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ForwarderIDs []int     `json:"forwarder_ids" gorm:"-"` // for add, update, delete
	Forwarders   []string  `json:"forwarders" gorm:"-"`    // for display
}

type ServerClassForwarder struct {
	ID            int       `json:"id"`
	ServerClassID int       `json:"server_class_id"`
	ForwarderID   int       `json:"forwarder_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FileInput struct {
	ID                  int    `json:"id"`
	LogFilePath         string `json:"log_file_path"`
	Sourcetype          string `json:"sourcetype"`
	LogFileSize         string `json:"log_file_size"`
	DataRetentionPeriod string `json:"data_retention_period"`
	Memo                string `json:"memo"`
	// Status              int       `json:"status"`         // to remove
	Env      string `json:"env"`       //added later
	UserName string `json:"user_name"` //added later
	// SplunkUserID int       `json:"splunk_user_id"` // to remove
	AppID     int       `json:"app_id"`
	AppName   string    `json:"app_name" gorm:"-"` //for display
	Blacklist string    `json:"blacklist"`
	Crcsalt   string    `json:"crcsalt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ScriptInput struct {
	ID                  int    `json:"id"`
	Sourcetype          string `json:"sourcetype"`
	LogFileSize         string `json:"log_file_size"`
	DataRetentionPeriod string `json:"data_retention_period"`
	Env                 string `json:"env"` //added later
	// Memo                string    `json:"memo"`
	// Status       int       `json:"status"`         // to remove
	UserName string `json:"user_name"` //added later
	// SplunkUserID int       `json:"splunk_user_id"` // to remove
	AppID      int       `json:"app_id"`
	AppName    string    `json:"app_name" gorm:"-"` //for display
	OS         string    `json:"os"`
	Interval   string    `json:"interval"`
	ScriptName string    `json:"script_name"`
	Script     []byte    `json:"script" gorm:"size:10240"`
	ScriptCode string    `json:"script_code" gorm:"-"`
	Option     string    `json:"option"`
	Exefile    bool      `json:"exefile"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UnixAppInput struct {
	ID int `json:"id"`
	// Sourcetype          string `json:"sourcetype"`
	DataRetentionPeriod string `json:"data_retention_period"`
	// Status              int       `json:"status"`         // to remove
	Env      string `json:"env"`       // added later
	UserName string `json:"user_name"` //added later
	// SplunkUserID int       `json:"splunk_user_id"` //to remove
	AppID      int       `json:"app_id"`
	AppName    string    `json:"app_name" gorm:"-"` //for display
	Interval   string    `json:"interval"`
	ScriptName string    `json:"script_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type App struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Env          string `json:"env"`
	UserName     string `json:"user_name"` //added later
	DeployStatus int    `json:"deploy_status"`
	// UserID          int            `json:"user_id"`        //to remove
	// SplunkHostID    int            `json:"splunk_host_id"` //to remove
	UnixApp         bool           `json:"unix_app"`
	FileInputIDs    []IDStruct     `json:"file_input_ids" gorm:"-"`     //add, update
	ScriptInputIDs  []IDStruct     `json:"script_input_ids" gorm:"-"`   //add, update
	UnixAppInputIDs []IDStruct     `json:"unix_app_input_ids" gorm:"-"` //add, update
	ServerClass     []ServerClass  `json:"server_classes" gorm:"-"`     //for display in deployment
	ServerClassIDs  []IDStruct     `json:"server_class_ids" gorm:"-"`   //for update in deployment
	FileInputs      []FileInput    `json:"file_inputs" gorm:"-"`        //display
	ScriptInputs    []ScriptInput  `json:"script_inputs" gorm:"-"`      //display
	UnixAppInputs   []UnixAppInput `json:"unix_app_inputs" gorm:"-"`    //display
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type Deployment struct {
	ID            int       `json:"id"`
	AppID         int       `json:"app_id"`
	ServerClassID int       `json:"server_class_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DeploymentApp struct { // no tables
	AppID      int        `json:"app_id" gorm:"-"`
	FolderName string     `json:"folder_name" gorm:"-"`
	InputsConf string     `json:"inputs_conf" gorm:"-"`
	ScriptIDs  []IDStruct `json:"script_ids" gorm:"-"`
	AppType    string     `json:"app_type" gorm:"-"`
	Env        string     `json:"env" gorm:"-"`
}

type DeploymentServerClass struct { // no tables
	ServerClassName string       `json:"server_class_name" gorm:"-"`
	AppName         string       `json:"app_name" gorm:"-"`
	AppID           int          `json:"app_id" gorm:"-"`
	ForwarderNames  []NameStruct `json:"forwarder_names" gorm:"-"`
	User            string       `json:"user" gorm:"-"`
	Env             string       `json:"env" gorm:"-"`
}

type Tag struct {
	Name        string   `json:"name"`
	Hosts       []string `json:"hosts"`
	SourceTypes []string `json:"source_types"`
	Sources     []string `json:"sources"`
	GUIDs       []string `json:"guids"`
}

type SamlGroup struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

type UnitPrice struct {
	ID           int       `json:"id"`
	ServicePrice int       `json:"service_price"`
	StoragePrice int       `json:"storage_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LogSize struct {
	ServiceID string `json:"service_id"`
	Host      string `json:"host"`
	SizeMB    int    `json:"size_mb"`
}

type StorageSize struct {
	ServiceID string `json:"service_id"`
	IndexName string `json:"index_name"`
	SizeMB    int    `json:"size_mb"`
}

// type DeploymentInfo struct { // for emailing
// 	UserName      string        `json:"user_name"`
// 	Env           string        `json:"env"`
// 	AppName       string        `json:"app_name"`
// 	Inputs        []string      `json:"inputs"`
// 	ServerClasses []ServerClass `json:"server_classes"` //contains sc name and forwarder list
// }
