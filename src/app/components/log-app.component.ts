import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { AppService, InputService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { App, Message, Forwarder, FileInput, ScriptInput, UnixAppInput } from '../model';
import * as com from "../common";
import * as _ from 'lodash';
import * as conf from "../configuration";
declare var $;

@Component({
  selector: 'log-app',
  templateUrl: '../../assets/templates/log-app.component.html',
  providers: [AppService, MessageService, InputService]
})

export class LogAppComponent extends CommonComponent{
  admin: string = ""; // used for select ele
  fileInputs: FileInput[] = []; //all selectable inputs
  scriptInputs: ScriptInput[] = []; //all selectable inputs
  unixAppInputs: UnixAppInput[] = []; //all selectable inputs

  selectedFileInputIDs: number[] = [];
  selectedScriptInputIDs: number[] = [];
  selectedUnixAppInputIDs: number[] = [];

  fileInputChecked: boolean[] = [];
  scriptInputChecked: boolean[] = [];
  unixAppInputChecked: boolean[] = [];

  apps: App[];
  currentApp: App = new App();
  oldApp: App = new App();
  env: string = localStorage.getItem("env");
  user: string = localStorage.getItem("user");
  appType: string = "app"; //or unixapp
  addModalID: string = "addAppModal";
  editModalID: string = "editAppModal";
  inputsLoaded : boolean[] = [false, false, false];

  unixappDashboardLink : string = conf.UNIXAPP_DASHBOARD_LINK;

  constructor(public messageService: MessageService, 
    public fb: FormBuilder, 
    private appService: AppService, 
    private inputService: InputService){
    super(messageService, fb, appService, "app");
    this.form = fb.group({
      "name": this.nameCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.getRecords();
    // this.currentApp.file_inputs = [];
    // this.currentApp.script_inputs = [];
    // this.currentApp.unix_app_inputs = [];
    // this.currentApp.file_input_ids = [];
    // this.currentApp.script_input_ids = [];
    // this.currentApp.unix_app_input_ids = [];
  }

  getRecords(){
    this.inputsLoaded = [false, false, false];
    var callback = (result: any[]) : void => {
      this.apps = result;
      this.getFileInputRecords();//get selectable inputs
      this.getScriptInputRecords(); //get selectable inputs
      this.getUnixAppInputRecords(); //get selectable inputs
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.apps.length; i++){
        this.apps[i].created_at = com.formatTime(this.apps[i].created_at);
        this.apps[i].file_input_ids = [];
        this.apps[i].script_input_ids = [];
        this.apps[i].unix_app_input_ids = [];
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 2; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.appService, callback, [this.env, this.user, ""]);
  }

  getFileInputRecords(){
    var callback = (result: any[]) : void => {
      this.fileInputs = result;
      this.inputsLoaded[0] = true;
      this.fileInputChecked = [];
      for(var i = 0; i < this.fileInputs.length; i++){
        this.fileInputChecked.push(false);
      }
      $(window).trigger('resize');
    }
    super._getRecords(this.inputService, callback, [this.env, this.user, "file", -1]);
  }

  getScriptInputRecords(){
    var callback = (result: any[]) : void => {
      this.scriptInputs = result;
      this.inputsLoaded[1] = true;
      this.scriptInputChecked = [];
      for(var i = 0; i < this.scriptInputs.length; i++){
        this.scriptInputChecked.push(false);
      }
      $(window).trigger('resize');
    }
    super._getRecords(this.inputService, callback, [this.env, this.user, "script", -1]);
  }

  getUnixAppInputRecords(){
    var callback = (result: any[]) : void => {
      this.unixAppInputs = result;
      this.inputsLoaded[2] = true;
      this.unixAppInputChecked = [];
      for(var i = 0; i < this.unixAppInputs.length; i++){
        this.unixAppInputChecked.push(false);
      }
      $(window).trigger('resize');
    }
    super._getRecords(this.inputService, callback, [this.env, this.user, "unixapp", -1]);
  }

  addRecord(){
    this.currentApp.file_input_ids = this.selectedFileInputIDs;
    this.currentApp.script_input_ids = this.selectedScriptInputIDs;
    this.currentApp.unix_app_input_ids = this.selectedUnixAppInputIDs;
    this.currentApp.unix_app = this.appType=="app"? false: true;
    this.currentApp.user_name = this.user;
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._addRecord(this.appService, this.currentApp.name, this.currentApp, callback, true);
  }

  editRecord(r){
    super.submitEdit();
    this.currentApp.file_input_ids = this.selectedFileInputIDs;
    this.currentApp.script_input_ids = this.selectedScriptInputIDs;
    this.currentApp.unix_app_input_ids = this.selectedUnixAppInputIDs;
    var callback = (result: any) : void => {
      if(result == "ok") {
        // for( var i = 0; i < this.apps.length; i++){
        //   if (this.apps[i].id == this.oldApp.id)
        //   {
        //     this.apps[i] = _.cloneDeep(this.currentApp); //if don't clone, when click add, the record editted last time will be set as empty
        //   }
        // }
        this.getRecords();
      }
    }
    super._editRecord(this.appService, this.currentApp.id, this.oldApp, this.currentApp, callback, true);
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._deleteRecord(this.appService, this.currentApp.id, callback, true);
  }

  openAddModal(){
    this.appType = "app";
    super.openAddModal();
    this.currentApp = new App();
    this.currentApp.env = this.env;
    this.selectedFileInputIDs = [];
    this.selectedScriptInputIDs = [];
    this.selectedUnixAppInputIDs = [];
    this.fileInputChecked = [];
    this.scriptInputChecked = [];
    this.unixAppInputChecked = [];
    this.currentApp.file_input_ids = [];
    this.currentApp.script_input_ids = [];
    this.currentApp.unix_app_input_ids = [];
    $('#' + this.addModalID).modal('show');
    this.modalResizeTimeOut();
  }

  openEditModal(app : App){
    this.currentApp = _.cloneDeep(app);
    this.oldApp = _.cloneDeep(app);
    this.selectedFileInputIDs = [];
    this.selectedScriptInputIDs = [];
    this.selectedUnixAppInputIDs = [];
    this.fileInputChecked = [];
    this.scriptInputChecked = [];
    this.unixAppInputChecked = [];
    if(!app.unix_app){
      for (var i = 0; i < app.file_inputs.length; i++){
        this.fileInputChecked.push(true);
        this.selectedFileInputIDs.push(app.file_inputs[i].id);
      }
      for (var i = 0; i < this.fileInputs.length; i++){
        this.fileInputChecked.push(false);
      }
      for (var i = 0; i < app.script_inputs.length; i++){
        this.scriptInputChecked.push(true);
        this.selectedScriptInputIDs.push(app.script_inputs[i].id);
      }
      for (var i = 0; i < this.scriptInputs.length; i++){
        this.scriptInputChecked.push(false);
      }
    }else{
      for (var i = 0; i < app.unix_app_inputs.length; i++){
        this.unixAppInputChecked.push(true);
        this.selectedUnixAppInputIDs.push(app.unix_app_inputs[i].id);
      }
      for (var i = 0; i < this.unixAppInputs.length; i++){
        this.unixAppInputChecked.push(false);
      }
    }

    $('#' + this.editModalID).modal('show');
    this.modalResizeTimeOut();
  }

  openDelModal(obj){
    this.currentApp = _.cloneDeep(obj);
  }
  openShowModal(obj){
    this.currentApp = _.cloneDeep(obj);
    this.modalResizeTimeOut();
  }

  hideModal(id: string){
    $("#" + id).modal('hide');
  }

  onChangeAppType(type: string){
    this.appType = type;
    this.modalResizeTimeOut();
  }

  modalResizeTimeOut(){
    setTimeout(() => {
      $(window).trigger('resize');
    }, 800);
  }

  onChangeInput(type:string, inputID:number, isChecked: boolean){
    if(type == "file"){
      if(isChecked){
        this.selectedFileInputIDs.push(inputID);
      }else{ //remove from list
        var index = this.selectedFileInputIDs.indexOf(inputID, 0);
        if (index > -1)
          this.selectedFileInputIDs.splice(index, 1);
      }
    }else if(type == "script"){
      if(isChecked){
        this.selectedScriptInputIDs.push(inputID);
      }else{ //remove from list
        var index = this.selectedScriptInputIDs.indexOf(inputID, 0);
        if (index > -1)
          this.selectedScriptInputIDs.splice(index, 1);
      }
    }else if(type == "unixapp"){
      if(isChecked){
        this.selectedUnixAppInputIDs.push(inputID);
      }else{ //remove from list
        var index = this.selectedUnixAppInputIDs.indexOf(inputID, 0);
        if (index > -1)
          this.selectedUnixAppInputIDs.splice(index, 1);
      }
    }
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.currentApp = _.cloneDeep(obj);
    this.oldApp = _.cloneDeep(obj);
  }
}
