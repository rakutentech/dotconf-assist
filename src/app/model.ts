//the attributes in each class have the same name in models of golang backend
export class User {
  constructor(
    public id:number = null,
    public user_name: string = null,
    public group_name: string = null,
    public app_team_name: string = null,
    // public service_id: number = null,
    public service_id: string = null,
    public email: string = null,
    public email_for_emergency: string = null,
    public password: string = null,
    public admin: boolean = null,
    public created_at: string = null,
    public updated_at: string = null,
    public last_login_at: string = null,
    public status: string = null
    ){
  }
}

export class SplunkUser {
  constructor(
    public id:number = null,
    public user_name: string = null,
    public rpaas_user_name: string = null,
    public env: string = null,
    public memo: number = null,
    public search_head: string = null,
    public email: string = null,
    public password: string = null,
    public position_ids: string = null,
    public created_at: string = null,
    public updated_at: string = null,
    // public status: string = null
    ){
  }
}

export class Announcement {
  constructor(
    public id:number = null,
    public content: string = null,
    public created_at: string = null
    ){
  }
}

export class SplunkHost {
  constructor(
    public id:number = null,
    public name: string = null,
    public role: string = null,
    public env: string = null,
    public created_at: string = null,
    public updated_at: string = null
    ){
  }
}

export class Forwarder {
  constructor(
    public id:number = null,
    public name: string = null,
    public env: string = null,
    public user_name: string = null,
    public server_classes: string[] = null, //for display
    public share: string = null,
    public created_at: string = null,
    public updated_at: string = null
    ){
  }
}

export class ServerClass {
  constructor(
    public id:number = null,
    public name: string = null,
    public env: string = null,
    public user_name: string = null,
    public forwarders: string[] = null, //for display
    public forwarder_ids: number[] = null, //for add, update
    public created_at: string = null,
    public updated_at: string = null
    ){
  }
}

export class ServerClassForwarder {
  constructor(
    public id:number = null,
    public server_class_id: number = null,
    public forwarder_id: number = null,
    public created_at: string = null,
    public updated_at: string = null
    ){
  }
}

export class CostAllocation {
  constructor(
    public id:number = null,
    public service_id: string = null,
    public rate: number = null,
    public mb: number = null
    ){
  }
}

export class FileInput {
  constructor(
    public id:number = null,
    public log_file_path:string = null,
    public sourcetype:string = null,
    public log_file_size:string = null,
    public data_retention_period:string = null,
    public memo:string = null,
    // public status:number = null, // to delete
    public env:string = null, //added later
    public user_name: string = null, //added later
    public splunk_user_id:number = null, //to remove
    public app_id:number = null, //for add, update
    public app_name:string = null, //for display
    public blacklist:string = null,
    public crcsalt:string = null,
    public created_at:string = null,
    public updated_at:string = null
    ){
  }
}

export class ScriptInput {
  constructor(
    public id:number = null,
    public sourcetype:string = null,
    public log_file_size:string = null,
    public data_retention_period:string = null,
    // public memo:string = null,
    // public status:number = null, // to delete
    public env:string = null, // added later
    public user_name: string = null, //added later
    // public splunk_user_id:number = null, //to remove
    public app_id:number = null, //for add, update
    public app_name:string = null, //for display
    public os:string = null,
    public interval:string = null,
    public script_name:string = null,
    public script: any = null,  //for add, update
    public script_code: string = null, // for display
    public option:string = null,
    public exefile:boolean = null,
    public created_at:string = null,
    public updated_at:string = null
    ){
  }
}

export class UnixAppInput {
  constructor(
    public id:number = null,
    public data_retention_period:string = null,
    // public status:number = null, // to delete
    public env:string = null, // added later
    public user_name: string = null, //added later
    // public splunk_user_id:number = null, //to remove
    public app_id:number = null, //for add, update
    public app_name:string = null, //for display
    public interval:string = null,
    public script_name:string = null,
    public created_at:string = null,
    public updated_at:string = null
    ){
  }
}

export class App {
  constructor(
    public id:number = null,
    public name:string = null,
    public env:string = null,
    public user_name: string = null, //added later
    public user_id:number = null, // to remove
    public deploy_status:number = null,
    public splunk_host_id:number = null, //deployment server
    public unix_app:boolean = null,
    public file_input_ids: number[] = null, //for add, update
    public file_inputs: FileInput[] = null, //for display
    public script_input_ids: number[] = null, //for add, update
    public script_inputs: ScriptInput[] = null, //for display
    public unix_app_input_ids: number[] = null, //for add, update
    public unix_app_inputs: UnixAppInput[] = null, //for display
    public server_class_ids: number[] = null, //for add, update deployment
    public server_classes: ServerClass[] = null, //for display in deployment
    public created_at:string = null,
    public updated_at:string = null
    ){
  }
}

export class Deployment {
  constructor(
    public id:number = null,
    public app_id:number = null,
    public server_class_id:number = null,
    public created_at:string = null,
    public updated_at:string = null
    ){
  }
}

export class DeploymentApp { //for creating deployment app
  constructor(
    public app_id:number = null,
    public folder_name:string = null,
    public inputs_conf:string = null,
    public script_ids:number[] = null,
    public app_type:string = null,
    public env: string = null
    ){
  }
}

export class DeploymentServerClass { //for creating server class
  constructor(
    public server_class_name:string = null, // sc_<user_name>_<sc_name>
    public app_name:string = null, //app foldername dsapp_<user_name>_<app_name>
    public app_id:number = null,
    public forwarder_names:string[] = null,
    public user: string = null,
    public env: string = null
    ){
  }
}

export class UnitPrice {
  constructor(
    public id:number = null,
    public service_price:number = null,
    public storage_price:number = null,
    public created_at:string = null,
    public updated_at:string = null
    ){
  }
}

export class LogSize {
  constructor(
    public service_id:string = null,
    public host:string = null,
    public size_mb:number = null
    ){
  }
}

export class StorageSize {
  constructor(
    public service_id:string = null,
    public index_name:string = null,
    public size_mb:number = null
    ){
  }
}
export class RouterOutlet {
    constructor(
    public name: string = null,
    public link: string = null,
    public icon: string = null
    ){
  }
}

export class Message {
    constructor(
    public type: string = null,
    public status: string = null,
    public text: string = null,
    public icon: string = null
    ){
  }
}
