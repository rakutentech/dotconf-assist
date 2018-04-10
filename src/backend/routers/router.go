package routers

import (
	"github.com/gorilla/mux"
	"github.com/rakutentech/dotconf-assist/src/backend/controllers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	{
		"Login", "POST", "/v1/login", controllers.LoginHandler,
	},
	{
		"Logout", "POST", "/v1/logout", controllers.LogoutHandler,
	},

	//For Preflight
	{
		"Preflight Login", "OPTIONS", "/v1/login", controllers.PreflightHandler,
	},
	{
		"Preflight Logout", "OPTIONS", "/v1/logout", controllers.PreflightHandler,
	},
	{
		"Update Session", "POST", "/v1/update_session", controllers.UpdateSessionHandler,
	},
	{
		"Get Token", "GET", "/v1/tokens", controllers.GetTokenHandler,
	},
	{
		"Preflight Get Token", "OPTIONS", "/v1/tokens", controllers.PreflightHandler,
	},
	{
		"Preflight Update Session", "OPTIONS", "/v1/update_session", controllers.PreflightHandler,
	},

	//users
	{
		"Add User", "POST", "/v1/users/reset_password", controllers.ResetUserPasswordHandler,
	},

	{
		"Add User", "POST", "/v1/users", controllers.AddUserHandler,
	},
	{
		"Get Users", "GET", "/v1/users", controllers.GetUsersHandler,
	},
	{
		"Get User", "GET", "/v1/users/{username}", controllers.GetUserHandler,
	},
	{
		"Update User", "PUT", "/v1/users/{username}", controllers.UpdateUserHandler,
	},

	{
		"Delete User", "DELETE", "/v1/users/{username}", controllers.DeleteUserHandler,
	},
	{
		"Preflight User", "OPTIONS", "/v1/user", controllers.PreflightHandler,
	},
	{
		"Preflight Users", "OPTIONS", "/v1/users", controllers.PreflightHandler,
	},
	{
		"Preflight Users Parameter", "OPTIONS", "/v1/users/{username}", controllers.PreflightHandler,
	},

	//splunk_users
	{
		"Add SplunkUser", "POST", "/v1/splunk_users", controllers.AddSplunkUserHandler,
	},
	{
		"Get SplunkUsers", "GET", "/v1/splunk_users", controllers.GetSplunkUsersHandler,
	},
	{
		"Get SplunkUser", "GET", "/v1/splunk_users/{username}", controllers.GetSplunkUserHandler,
	},
	{
		"Update SplunkUser", "PUT", "/v1/splunk_users/{username}", controllers.UpdateSplunkUserHandler,
	},
	{
		"Delete SplunkUser", "DELETE", "/v1/splunk_users/{username}", controllers.DeleteSplunkUserHandler,
	},
	{
		"Preflight SplunkUser", "OPTIONS", "/v1/splunk_user", controllers.PreflightHandler,
	},
	{
		"Preflight SplunkUsers", "OPTIONS", "/v1/splunk_users", controllers.PreflightHandler,
	},
	{
		"Preflight SplunkUsers Parameter", "OPTIONS", "/v1/splunk_users/{username}", controllers.PreflightHandler,
	},

	//announcements
	{
		"Add Announcement", "POST", "/v1/announcements", controllers.AddAnnouncementHandler,
	},
	{
		"Get Announcements", "GET", "/v1/announcements", controllers.GetAnnouncementsHandler,
	},
	{
		"Get Announcement", "GET", "/v1/announcements/{id}", controllers.GetAnnouncementHandler,
	},
	{
		"Update Announcement", "PUT", "/v1/announcements/{id}", controllers.UpdateAnnouncementHandler,
	},
	{
		"Delete Announcement", "DELETE", "/v1/announcements/{id}", controllers.DeleteAnnouncementHandler,
	},
	{
		"Preflight Announcements", "OPTIONS", "/v1/announcements", controllers.PreflightHandler,
	},
	{
		"Preflight Announcements Parameter", "OPTIONS", "/v1/announcements/{id}", controllers.PreflightHandler,
	},

	//splunk hosts
	{
		"Add Splunk Host", "POST", "/v1/splunk_hosts", controllers.AddSplunkHostHandler,
	},
	{
		"Get Splunk Hosts", "GET", "/v1/splunk_hosts", controllers.GetSplunkHostsHandler,
	},
	{
		"Get Splunk Host", "GET", "/v1/splunk_hosts/{name}", controllers.GetSplunkHostHandler,
	},
	{
		"Update Splunk Host", "PUT", "/v1/splunk_hosts/{name}", controllers.UpdateSplunkHostHandler,
	},
	{
		"Delete Splunk Host", "DELETE", "/v1/splunk_hosts/{name}", controllers.DeleteSplunkHostHandler,
	},
	{
		"Preflight Splunk Hosts", "OPTIONS", "/v1/splunk_hosts", controllers.PreflightHandler,
	},
	{
		"Preflight Splunk Host Parameter", "OPTIONS", "/v1/splunk_hosts/{name}", controllers.PreflightHandler,
	},

	//Forwarders
	{
		"Add Forwarder", "POST", "/v1/forwarders", controllers.AddForwarderHandler,
	},
	{
		"Get Forwarders", "GET", "/v1/forwarders", controllers.GetForwardersHandler,
	},
	{
		"Get Forwarder", "GET", "/v1/forwarders/{name}", controllers.GetForwarderHandler,
	},
	{
		"Update Forwarder", "PUT", "/v1/forwarders/{name}", controllers.UpdateForwarderHandler,
	},
	{
		"Delete Forwarder", "DELETE", "/v1/forwarders/{name}", controllers.DeleteForwarderHandler,
	},
	{
		"Preflight Forwarders", "OPTIONS", "/v1/forwarders", controllers.PreflightHandler,
	},
	{
		"Preflight Forwarder Parameter", "OPTIONS", "/v1/forwarders/{name}", controllers.PreflightHandler,
	},

	//ServerClasses
	{
		"Add ServerClass", "POST", "/v1/server_classes", controllers.AddServerClassHandler,
	},
	{
		"Get ServerClasses", "GET", "/v1/server_classes", controllers.GetServerClassesHandler,
	},
	{
		"Get ServerClass", "GET", "/v1/server_classes/{name}", controllers.GetServerClassHandler,
	},
	{
		"Update ServerClass", "PUT", "/v1/server_classes/{name}", controllers.UpdateServerClassHandler,
	},
	{
		"Delete ServerClass", "DELETE", "/v1/server_classes/{name}", controllers.DeleteServerClassHandler,
	},
	{
		"Preflight ServerClasses", "OPTIONS", "/v1/server_classes", controllers.PreflightHandler,
	},
	{
		"Preflight ServerClass Parameter", "OPTIONS", "/v1/server_classes/{name}", controllers.PreflightHandler,
	},

	//Inputs
	{
		"Add Input", "POST", "/v1/inputs", controllers.AddInputHandler,
	},
	{
		"Get Inputs", "GET", "/v1/inputs", controllers.GetInputsHandler,
	},
	{
		"Get Input", "GET", "/v1/inputs/{id}", controllers.GetInputHandler,
	},
	{
		"Update Input", "PUT", "/v1/inputs/{id}", controllers.UpdateInputHandler,
	},
	{
		"Delete Input", "DELETE", "/v1/inputs/{id}", controllers.DeleteInputHandler,
	},
	{
		"Preflight Inputs", "OPTIONS", "/v1/inputs", controllers.PreflightHandler,
	},
	{
		"Preflight Input Parameter", "OPTIONS", "/v1/inputs/{id}", controllers.PreflightHandler,
	},

	//Apps
	{
		"Add App", "POST", "/v1/apps", controllers.AddAppHandler,
	},
	{
		"Get Apps", "GET", "/v1/apps", controllers.GetAppsHandler,
	},
	{
		"Get App", "GET", "/v1/apps/{id}", controllers.GetAppHandler,
	},
	{
		"Update App", "PUT", "/v1/apps/{id}", controllers.UpdateAppHandler,
	},
	{
		"Delete App", "DELETE", "/v1/apps/{id}", controllers.DeleteAppHandler,
	},
	{
		"Preflight Apps", "OPTIONS", "/v1/apps", controllers.PreflightHandler,
	},
	{
		"Preflight App Parameter", "OPTIONS", "/v1/apps/{id}", controllers.PreflightHandler,
	},

	//Deployment
	{
		"Add Deployment", "POST", "/v1/deployment", controllers.AddDeploymentHandler,
	},

	{
		"Create Deployment App", "POST", "/v1/deployment/create_app", controllers.CreateDeploymentAppHandler,
	},

	{
		"Create Deployment Server Class", "POST", "/v1/deployment/create_server_class", controllers.CreateDeploymentServerClassHandler,
	},

	{
		"Deinstall Deployment App", "POST", "/v1/deployment/deinstall_app", controllers.DeinstallDeploymentAppHandler,
	},

	// {
	// 	"Create Tag", "POST", "/v1/deployment/create_tag", controllers.CreateTagHandler,
	// },

	// {
	// 	"Change Tag permission", "POST", "/v1/deployment/change_tag_permission", controllers.ChangeTagPermissionHandler,
	// },

	{
		"Get Deployment List", "GET", "/v1/deployment", controllers.GetDeploymentListHandler,
	},
	{
		"Get Deployment", "GET", "/v1/deployment/{id}", controllers.GetDeploymentHandler,
	},
	{
		"Update Deployment", "PUT", "/v1/deployment", controllers.UpdateDeploymentHandler,
	},
	{
		"Delete Deployment", "DELETE", "/v1/deployment/{id}", controllers.DeleteDeploymentHandler,
	},
	{
		"Preflight Deployment", "OPTIONS", "/v1/deployment", controllers.PreflightHandler,
	},
	{
		"Preflight Deployment Parameter", "OPTIONS", "/v1/deployment/{id}", controllers.PreflightHandler,
	},

	//unit price
	{
		"Add Unit Price", "POST", "/v1/unit_price", controllers.AddUnitPriceHandler,
	},
	{
		"Get Unit Prices", "GET", "/v1/unit_price", controllers.GetUnitPricesHandler,
	},
	{
		"Get Unit Price", "GET", "/v1/unit_price/{id}", controllers.GetUnitPriceHandler,
	},
	{
		"Update Unit Price", "PUT", "/v1/unit_price/{id}", controllers.UpdateUnitPriceHandler,
	},
	{
		"Delete Unit Price", "DELETE", "/v1/unit_price/{id}", controllers.DeleteUnitPriceHandler,
	},
	{
		"Preflight Unit Prices", "OPTIONS", "/v1/unit_price", controllers.PreflightHandler,
	},
	{
		"Preflight Unit Price Parameter", "OPTIONS", "/v1/unit_price/{id}", controllers.PreflightHandler,
	},

	//usage
	{
		"Get Usages", "GET", "/v1/usage", controllers.GetUsagesHandler,
	},
	{
		"Preflight Usages", "OPTIONS", "/v1/usage", controllers.PreflightHandler,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		// handler = settings.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
