import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { ServerClassService, ForwarderService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { ServerClass, Message, Forwarder } from '../model';
import * as com from "../common";
import * as _ from 'lodash';
declare var $;

@Component({
  selector: 'serverClass',
  templateUrl: '../../assets/templates/server-class.component.html',
  providers: [ServerClassService, MessageService, ForwarderService]
})

export class ServerClassComponent extends CommonComponent{
  forwarders: Forwarder[] = [];
  selectedForwardersID: number[] = [];
  forwarderChecked: boolean[] = [];
  serverClasses: ServerClass[];
  currentServerClass: ServerClass = new ServerClass();
  oldServerClass: ServerClass = new ServerClass();
  env: string = localStorage.getItem("env");
  user: string = localStorage.getItem("user");
  addModalID: string = "addServerClassModal";
  editModalID: string = "editServerClassModal";
  popupSCNameStr: string = "Server Class Name CANNOT contain \"[\" or \"]\"";

  constructor(public messageService: MessageService, public fb: FormBuilder, private serverClassService: ServerClassService, private forwarderService: ForwarderService){
    super(messageService, fb, serverClassService, "server class");
    this.form = fb.group({
      "name": this.nameCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.getRecords();
  }

  getRecords(){
    var callback = (result: any[]) : void => {
      this.serverClasses = result;
      this.getForwarders();
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.serverClasses.length; i++){
        this.serverClasses[i].created_at = com.formatTime(this.serverClasses[i].created_at);
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 2; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.serverClassService, callback, [this.env, this.user]);
  }

  getForwarders(){
    var callback = (result: any[]) : void => {
      this.forwarders = result;
      this.forwarderChecked = [];
      for(var i = 0; i < this.forwarders.length; i++){
        this.forwarderChecked.push(false);
      }
      $(window).trigger('resize');
    }
    super._getRecords(this.forwarderService, callback, [this.env, this.user, "local"]);
  }

  getSCFwdrs(sc: ServerClass){
    if (sc.forwarders != null) {
      return
    }
    sc.forwarders = [];
    for (var i = 0; i < sc.forwarder_ids.length; i++){
      for (var j = 0; j < this.forwarders.length; j++){
        if (sc.forwarder_ids[i] == this.forwarders[j].id){
          sc.forwarders.push(this.forwarders[j].name);
          break;
        }
      }
    } 
  }

  addRecord(){
    this.currentServerClass.forwarder_ids = this.selectedForwardersID;
    // console.log(this.currentServerClass.forwarder_ids);
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._addRecord(this.serverClassService, this.currentServerClass.name, this.currentServerClass, callback);
  }

  hideModal(id: string){
    $("#" + id).modal('hide');
  }
  
  editRecord(r){
    super.submitEdit();
    this.currentServerClass.forwarder_ids = this.selectedForwardersID;
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
        // for( var i = 0; i < this.serverClasses.length; i++){
        //   if (this.serverClasses[i].name == this.oldServerClass.name)
        //   {
        //     this.serverClasses[i] = _.cloneDeep(this.currentServerClass); //if don't clone, when click add, the record editted last time will be set as empty
        //     this.serverClasses[i].forwarders = null;
        //   }
        // }
      }
    }
    super._editRecord(this.serverClassService, this.currentServerClass.name, this.oldServerClass, this.currentServerClass, callback);
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._deleteRecord(this.serverClassService, this.currentServerClass.name, callback, true, this.currentServerClass.user_name);
  }

  openAddModal(){
    super.openAddModal();
    this.currentServerClass = new ServerClass();
    this.currentServerClass.user_name = this.user;
    this.currentServerClass.env = this.env;
    this.selectedForwardersID = [];
    $('#' + this.addModalID).modal('show');
  }

  openEditModal(sc : ServerClass){
    this.currentServerClass = _.cloneDeep(sc);
    this.oldServerClass = _.cloneDeep(sc);
    this.selectedForwardersID = [];
    for(var i = 0; i < this.forwarders.length; i++){
      var found = false;
      if(sc.forwarders != null){
        for (var j = 0; j < sc.forwarders.length; j++){
          if (this.forwarders[i].name == sc.forwarders[j]){
            this.selectedForwardersID.push(this.forwarders[i].id);
            this.forwarderChecked[i] = true;
            found = true;
            break;
          }
        }
      }
      
      if (!found){
        this.forwarderChecked[i] = false;
      }
    }
    $('#' + this.editModalID).modal('show');
  }

  openDelModal(obj){
    this.currentServerClass = _.cloneDeep(obj);
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.currentServerClass = _.cloneDeep(obj);
    this.selectedForwardersID = [];
    if(obj.forwarders != null){
      for(var i = 0; i < this.forwarders.length; i++){
        for (var j = 0; j < obj.forwarders.length; j++){
          if (this.forwarders[i].name == obj.forwarders[j]){
            this.selectedForwardersID.push(this.forwarders[i].id);
            break;
          }
        }
      }
    }
    this.oldServerClass = _.cloneDeep(obj);
  }

  onChangeFwdr(fwdrID:number, isChecked: boolean){
    if(isChecked){
      this.selectedForwardersID.push(fwdrID);
    }else{ //remove from list
      var index = this.selectedForwardersID.indexOf(fwdrID, 0);
      if (index > -1) {
         this.selectedForwardersID.splice(index, 1);
      }
    }
  }
}
