import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { SplunkHostService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { SplunkHost, Message } from '../model';
import * as com from "../common";
import * as _ from 'lodash';

@Component({
  selector: 'splunkHost',
  templateUrl: '../../assets/templates/splunk-host.component.html',
  providers: [SplunkHostService, MessageService]
})

export class SplunkHostComponent extends CommonComponent{
  admin: string = ""; // used for select ele
  splunkHosts: SplunkHost[];
  currentSplunkHost: SplunkHost = new SplunkHost();
  oldSplunkHost: SplunkHost = new SplunkHost();
  constructor(public messageService: MessageService, public fb: FormBuilder, private splunkHostService: SplunkHostService){
    super(messageService, fb, splunkHostService, "splunk host");
    this.form = fb.group({
      "name": this.nameCtl,
      "role": this.roleCtl,
      "env": this.envCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.getRecords();
  }

  getRecords(){
    var callback = (result: any[]) : void => {
      this.splunkHosts = result;
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.splunkHosts.length; i++){
        this.splunkHosts[i].created_at = com.formatTime(this.splunkHosts[i].created_at);
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 3; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.splunkHostService, callback, ["", ""]);
  }

  addRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._addRecord(this.splunkHostService, this.currentSplunkHost.name, this.currentSplunkHost, callback);
  }
  
  editRecord(r){
    super.submitEdit();
    var callback = (result: any) : void => {
      if(result == "ok") {
        for( var i = 0; i < this.splunkHosts.length; i++){
          if (this.splunkHosts[i].name == this.oldSplunkHost.name)
          {
            this.splunkHosts[i] = _.cloneDeep(this.currentSplunkHost); //if don't clone, when click add, the record editted last time will be set as empty
          }
        }
      }
    }
    super._editRecord(this.splunkHostService, this.currentSplunkHost.name, this.oldSplunkHost, this.currentSplunkHost, callback);
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._deleteRecord(this.splunkHostService, this.currentSplunkHost.name, callback);
  }

  openAddModal(){
    super.openAddModal();
    this.currentSplunkHost.role = "SH";
    this.currentSplunkHost.env = "prod";
  }

  openDelModal(obj){
    this.currentSplunkHost = _.cloneDeep(obj);
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.admin = obj.admin? "true": "false";
    this.currentSplunkHost = _.cloneDeep(obj);
    this.oldSplunkHost = _.cloneDeep(obj);
  }
}
