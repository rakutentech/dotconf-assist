import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { ForwarderService, UserService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Forwarder, Message, User, ServerClass } from '../model';
import * as com from "../common";
import * as _ from 'lodash';
declare var $;

@Component({
  selector: 'forwarder',
  templateUrl: '../../assets/templates/forwarder.component.html',
  providers: [ForwarderService, MessageService, UserService]
})

export class ForwarderComponent extends CommonComponent{
  forwardersFromDeploymentServer: Forwarder[];
  selectedForwarders: Forwarder[] = [];
  selectedForwardersStr: string[] = [];
  forwarderChecked: boolean[] = [];
  forwarders: Forwarder[];
  serverClasses: ServerClass[];
  users: User[];
  currentForwarder: Forwarder = new Forwarder();
  oldForwarder: Forwarder = new Forwarder();
  env: string = localStorage.getItem("env");
  user: string = localStorage.getItem("user");

  constructor(public messageService: MessageService, public fb: FormBuilder, 
    private forwarderService: ForwarderService, 
    private userService: UserService){
    super(messageService, fb, forwarderService, "forwarder");
    this.form = fb.group({
      "name": this.nameCtl,
      "user": this.nameCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.getRecords();
    this.getUsers();
  }

  getRecords(){
    var callback = (result: any[]) : void => {
      this.forwarders = result;
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.forwarders.length; i++){
        this.forwarders[i].created_at = com.formatTime(this.forwarders[i].created_at);
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 1; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.forwarderService, callback, [this.env, this.user, "local"]);
  }

  getUsers(){
    var callback = (result: any[]) : void => {
      this.users = result;
    }
    super._getRecords(this.userService, callback);
  }

  getRecordsFromDeploymentServer(){
    var callback = (result: any[]) : void => {
      for (var i = 0; i < result.length; i++){
        var found = false;
        for (var j = 0; j < this.forwarders.length; j++){
          if (this.forwarders[j].name == result[i].name){
            found = true;
            break;
          }
        }
        if (!found){ //filter the forwarders already added
          this.forwardersFromDeploymentServer.push(result[i]);
          this.forwarderChecked.push(false);
        }
      }
      
      $(window).trigger('resize');
      // $('#addForwarderModal').modal('refresh'); //doesn't work
    }
    super._getRecords(this.forwarderService, callback, [this.env, this.user, "deployment_server"]);
  }

  addRecord(){
    for(var i = 0; i < this.selectedForwardersStr.length; i++){
      this.selectedForwarders.push(new Forwarder(0, this.selectedForwardersStr[i], this.env, this.user, [], "", ""))
    }
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    // super._addRecord(this.forwarderService, this.currentForwarder.name, this.currentForwarder, callback);
    super._addRecord(this.forwarderService, this.selectedForwarders.length + " forwarders", this.selectedForwarders, callback);
  }

  hideAddModal(){
    $('#addForwarderModal').modal('hide');
  }
  
  editRecord(r){
    super.submitEdit();
    var callback = (result: any) : void => {
      if(result == "ok") {
        for( var i = 0; i < this.forwarders.length; i++){
          if (this.forwarders[i].name == this.oldForwarder.name)
          {
            this.forwarders[i] = _.cloneDeep(this.currentForwarder); //if don't clone, when click add, the record editted last time will be set as empty
          }
        }
      }
    }
    super._editRecord(this.forwarderService, this.currentForwarder.name, this.oldForwarder, this.currentForwarder, callback);
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._deleteRecord(this.forwarderService, this.currentForwarder.name, callback, true, this.currentForwarder.user_name);
  }

  openAddModal(){
    super.openAddModal();
    this.selectedForwarders = [];
    this.selectedForwardersStr = [];
    this.forwardersFromDeploymentServer = [];
    this.getRecordsFromDeploymentServer();
    this.forwarderChecked = [];
    $('#addForwarderModal').modal('show');
  }

  openDelModal(obj){
    this.currentForwarder = _.cloneDeep(obj);
  }
  
  onChangeFwdr(fwdr:string, isChecked: boolean){
    if(isChecked){
      this.selectedForwardersStr.push(fwdr);
    }else{ //remove from list
      var index = this.selectedForwardersStr.indexOf(fwdr, 0);
      if (index > -1) {
         this.selectedForwardersStr.splice(index, 1);
      }
    }
  }
  trackByIndex(index) {
    return index;
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.currentForwarder = _.cloneDeep(obj);
    this.oldForwarder = _.cloneDeep(obj);
  }

  // submitEdit(row, col){
  //   super.submitEdit();
  // }
}
