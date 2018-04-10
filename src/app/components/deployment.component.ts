import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { DeploymentService, ForwarderService, AppService, ServerClassService, InputService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Deployment, Message, Forwarder, App, ServerClass, FileInput, ScriptInput, UnixAppInput, DeploymentApp, DeploymentServerClass } from '../model';
import * as com from "../common";
import * as _ from 'lodash';
declare var $;

@Component({
  selector: 'deployment',
  templateUrl: '../../assets/templates/deployment.component.html',
  providers: [DeploymentService, MessageService, ForwarderService, AppService, ServerClassService, InputService]
})

export class DeploymentComponent extends CommonComponent{
  serverClasses: ServerClass[] = [];
  forwarders: Forwarder[] = [];
  selectedSCIDs: number[] = [];
  scChecked: boolean[] = [];
  deploymentList: Deployment[];
  apps: App[];
  currentApp: App = new App();
  currentDeployment: Deployment = new Deployment();
  oldDeployment: Deployment = new Deployment();
  env: string = localStorage.getItem("env");
  user: string = localStorage.getItem("user");
  addModalID: string = "addDeploymentModal";
  editModalID: string = "editDeploymentModal";
  resourceLoaded : boolean[] = [false, false];

  constructor(public messageService: MessageService, 
    public fb: FormBuilder, 
    private deploymentService: DeploymentService, 
    private forwarderService: ForwarderService,
    private appService: AppService,
    private serverClassService: ServerClassService,
    private inputService: InputService){
    super(messageService, fb, deploymentService, "deployment");
    this.form = fb.group({
      "name": this.nameCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.resourceLoaded = [false, false];
    this.getApps();
    this.getServerClasses();
  }

  getApps(){
    this.resourceLoaded[0] = false;
    var callback = (result: any[]) : void => {
      this.apps = result;
      this.resourceLoaded[0] = true;
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.apps.length; i++){
        this.apps[i].created_at = com.formatTime(this.apps[i].created_at);
        this.apps[i].file_input_ids = [];
        this.apps[i].script_input_ids = [];
        this.apps[i].unix_app_input_ids = [];
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 1; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.appService, callback, [this.env, this.user, ""]);
  }

  // getRecords(){
  //   var callback = (result: any[]) : void => {
  //     this.deploymentList = result;
  //     this.getForwarders();
  //     for (var i = 0; i < this.deploymentList.length; i++){
  //       this.deploymentList[i].created_at = com.formatTime(this.deploymentList[i].created_at);
  //     }
  //   }
  //   super._getRecords(this.deploymentService, callback, [this.env, this.user, ""]);
  // }

  getForwarders(){
    this.resourceLoaded[1] = false;
    var callback = (result: any[]) : void => {
      this.forwarders = result;
      this.resourceLoaded[1] = true;
      $(window).trigger('resize');
    }
    super._getRecords(this.forwarderService, callback, [this.env, this.user]);
  }

  getServerClasses(){
    var callback = (result: any[]) : void => {
      this.serverClasses = result;
      this.getForwarders();
    }
    super._getRecords(this.serverClassService, callback, [this.env, this.user, ""]);
  }

  // getSCFwdrs(sc: ServerClass){
  //   if (sc.forwarders != null) {
  //     return
  //   }
  //   sc.forwarders = [];
  //   for (var i = 0; i < sc.forwarder_ids.length; i++){
  //     for (var j = 0; j < this.forwarders.length; j++){
  //       if (sc.forwarder_ids[i] == this.forwarders[j].id){
  //         sc.forwarders.push(this.forwarders[j].name);
  //         break;
  //       }
  //     }
  //   } 
  // }

  hideModal(id: string){
    $('#' + id).modal('hide');
  }
  
  editRecord(r){
    super.submitEdit();
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getApps();
      }
    }
    super._editRecord(this.deploymentService, this.currentApp.name, this.currentApp.id, this.selectedSCIDs, callback);
  }

  // deleteDeployment(){
  //   var callback = (result: any) : void => {
  //     if(result == "ok") {
  //       this.getApps();
  //     }
  //   }
  //   super._deleteRecord(this.deploymentService, this.currentApp.id, callback, true, this.currentApp.id);
  // }

  openApplyModal(app: App){
    this.currentApp = _.cloneDeep(app);
  }

  applyDeployment(){
    // console.log(this.currentApp);
    this.createDeploymentApp();
  }

  getIndexName(retention: string){
    var indexName = "idx_common_1mon";
    switch (retention) {
      case "1 Day":
        indexName = "idx_common_1d";
        break;
      case "2 Days":
        indexName = "idx_common_2d";
        break;
      case "3 Days":
        indexName = "idx_common_3d";
        break;
      case "1 Week":
        indexName = "idx_common_1w";
        break;
      case "2 Weeks":
        indexName = "idx_common_2w";
        break;
      case "3 Weeks":
        indexName = "idx_common_3w";
        break;
      case "1 Month":
        indexName = "idx_common_1mon";
        break;
      case "3 Months":
        indexName = "idx_common_3mon";
        break;
      case "6 Months":
        indexName = "idx_common_6mon";
        break;
      case "1 Year":
        indexName = "idx_common_1y";
        break;
      case "3 Years":
        indexName = "idx_common_3y";
        break;
      case "6 Years":
        indexName = "idx_common_6y";
        break;
      default:
        indexName = "idx_common_1mon";
        break;
    }
    return indexName;
  }

  createDeploymentApp(){
    // console.log(this.currentApp);
    var deploymentApp: DeploymentApp = new DeploymentApp();
    deploymentApp.folder_name = "dsapp_" + this.user + "_" + this.currentApp.name;
    deploymentApp.env = this.env;
    deploymentApp.inputs_conf = "";
    deploymentApp.app_id = this.currentApp.id;
    deploymentApp.script_ids = [];
    // console.log(this.currentApp.unix_app);
    if (!this.currentApp.unix_app) {
      deploymentApp.app_type = "app";
      for (var i = 0; i < this.currentApp.file_inputs.length; i++) {
        var re = /\\/gi;
        var blacklist = this.currentApp.file_inputs[i].blacklist;
        var logFilePath = this.currentApp.file_inputs[i].log_file_path;
        blacklist = blacklist.replace(re, "\\\\");
        logFilePath = logFilePath.replace(re, "\\\\");
        deploymentApp.inputs_conf += "[monitor://" + logFilePath + "]\\n";
        deploymentApp.inputs_conf += "disabled=false" + "\\n";
        if (this.currentApp.file_inputs[i].sourcetype != "auto") {
          deploymentApp.inputs_conf += "sourcetype=" + this.currentApp.file_inputs[i].sourcetype + "\\n";
        }
        deploymentApp.inputs_conf += "index=" + this.getIndexName(this.currentApp.file_inputs[i].data_retention_period) + "\\n";
        if(this.currentApp.file_inputs[i].crcsalt == "<SOURCE>")
          deploymentApp.inputs_conf += "crcSalt=<SOURCE>" + "\\n";
        if (blacklist != "") {
          deploymentApp.inputs_conf += "blacklist=" + blacklist + "\\n";
        }
        deploymentApp.inputs_conf += "ignoreOlderThan=10d" + "\\n\\n";
      }

      for (var i = 0; i < this.currentApp.script_inputs.length; i++) {
        if (this.currentApp.script_inputs[i].exefile) {
          var option = this.currentApp.script_inputs[i].option;
          option = option==""? option : " " + option;
          if(this.currentApp.script_inputs[i].os == "*nix")
            deploymentApp.inputs_conf += "[script://./bin/" + this.currentApp.script_inputs[i].script_name + option + "]\\n";
          else
            deploymentApp.inputs_conf += "[script://.\\\\bin\\\\" + this.currentApp.script_inputs[i].script_name + option + "]\\n";
          deploymentApp.inputs_conf += "disabled=false" + "\\n";
          deploymentApp.inputs_conf += "sourcetype=" + this.currentApp.script_inputs[i].sourcetype + "\\n";
          deploymentApp.inputs_conf += "index=" + this.getIndexName(this.currentApp.script_inputs[i].data_retention_period) + "\\n";
          deploymentApp.inputs_conf += "interval=" + this.currentApp.script_inputs[i].interval + "\\n\\n";
        }
        deploymentApp.script_ids.push(this.currentApp.script_inputs[i].id);
      }
    }else{
      deploymentApp.script_ids = [];
      deploymentApp.app_type = "unixapp";
      for (var i = 0; i < this.currentApp.unix_app_inputs.length; i++) {
          deploymentApp.inputs_conf += "[script://./bin/" + this.currentApp.unix_app_inputs[i].script_name + "]\\n";
          deploymentApp.inputs_conf += "disabled=false" + "\\n";
          deploymentApp.inputs_conf += "index=" + this.getIndexName(this.currentApp.unix_app_inputs[i].data_retention_period) + "\\n";
          deploymentApp.inputs_conf += "interval=" + this.currentApp.unix_app_inputs[i].interval + "\\n\\n";
        deploymentApp.script_ids.push(this.currentApp.unix_app_inputs[i].id);
      }
    }
    
    var callback = (result: any) : void => {
      if(result == "ok") {
        // console.log("created deployment app");
        this.createDeploymentServerClass();
      }else{
        this.getApps();
      }
    }
    super._createDeploymentApp(this.deploymentService, this.currentApp.name, deploymentApp, callback, true, "deployment app");
  }

  deinstallDeploymentApp(){ // Remove mapping of app from all server classes and delete it from client target repositories, 
  //server class remains, need to delete the folder under deployment-apps from deployment server with linux command
    // console.log(this.currentApp);
    var deploymentApp: DeploymentApp = new DeploymentApp();
    deploymentApp.folder_name = "dsapp_" + this.user + "_" + this.currentApp.name;
    deploymentApp.app_id = this.currentApp.id;
    deploymentApp.env = this.env;
    var callback = (result: any) : void => {
      if(result == "ok") {
        // console.log("deinstalled deployment app");
        this.getApps();
      }
    }
    super._deinstallDeploymentApp(this.deploymentService, this.currentApp.name, deploymentApp, callback, true, "deployment app");
  }

  createDeploymentServerClass(){
    // console.log(this.currentApp.server_classes);
    var deploymentServerClasses: DeploymentServerClass[] = [];
    for (var i = 0; i < this.currentApp.server_classes.length; i ++){
      if (this.currentApp.server_classes[i].forwarders != null) {
        var dplySC : DeploymentServerClass = new DeploymentServerClass();
        dplySC.forwarder_names = [];
        dplySC.env = this.env;
        dplySC.user = this.user;
        dplySC.app_name = "dsapp_" + this.user + "_" + this.currentApp.name;
        dplySC.app_id = this.currentApp.id;
        dplySC.server_class_name = "sc_" + this.user + "_" + this.currentApp.server_classes[i].name;
        for (var j = 0; j < this.currentApp.server_classes[i].forwarders.length; j++){
          dplySC.forwarder_names.push(this.currentApp.server_classes[i].forwarders[j]);
        }
        deploymentServerClasses.push(dplySC);
      }
    }
    // console.log(deploymentServerClasses);

    var callback = (result: any) : void => {
      if(result == "ok") {
        // console.log("created deployment server class");
      }
      this.getApps();
    }
    super._createDeploymentServerClass(this.deploymentService, this.currentApp.name, deploymentServerClasses, callback, true, "deployment app");
  }

  createTag(){

  }

  changeTagPermission(){

  }

  openEditModal(app: App){
    this.currentApp = new App();
    this.currentApp = _.cloneDeep(app);
    this.scChecked = [];
    this.selectedSCIDs = [];
    for( var i = 0; i < this.serverClasses.length; i++){
      var found = false;
      for( var j = 0; j < this.currentApp.server_classes.length; j++){
        if(this.serverClasses[i].id == this.currentApp.server_classes[j].id)
        {
          this.selectedSCIDs.push(this.serverClasses[i].id);
          found = true;
          break;
        }
      }
      if(found)
        this.scChecked.push(true);
      else
        this.scChecked.push(false);
    }
    $('#' + this.editModalID).modal('show');
  }

  openDelModal(app: App){
    // this.currentDeployment = _.cloneDeep(obj);
    this.currentApp = new App();
    this.currentApp = _.cloneDeep(app);
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.currentDeployment = _.cloneDeep(obj);
    this.oldDeployment = _.cloneDeep(obj);
  }

  onChangeSC(scID:number, isChecked: boolean){
    if(isChecked){
      this.selectedSCIDs.push(scID);
    }else{ //remove from list
      var index = this.selectedSCIDs.indexOf(scID, 0);
      if (index > -1) {
         this.selectedSCIDs.splice(index, 1);
      }
    }
  }
}
