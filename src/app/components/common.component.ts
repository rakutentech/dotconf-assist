import { Component, OnInit } from "@angular/core";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Router, RouterStateSnapshot } from '@angular/router';
import { Message } from '../model';
import * as com from "../common";
import * as _ from 'lodash';

@Component({
  providers: [MessageService]
})

export class CommonComponent implements OnInit{
  //records
  records: any[] = [];
  record: any = null;

  //sorting
  activeTD: string = "";
  orderAttr: string = "+id";
  arrow: string = "sort ascending";
  activeAttr: string = "id";

  //inline editor
  isEditing: boolean = false;
  eleActive: boolean[][] = [];
  eleHovered: boolean[][] = [];
  activeCell: number[] = [-1, -1];
  
  //message
  msg : Message = new Message();
  messages: Array<string> = [];
  duration: number = 4000; //ms

  //table pagination
  itemNumPerPage: number = 10;
  subItemNumPerPage: number = 10; // for sub tables
  curentP: number = 1;
  subCurrentP: number = 1; // for sub tables

  //form
  form: FormGroup;
  userForm: FormGroup;
  
  selectControl: FormControl = new FormControl(false);
  subSelectControl: FormControl = new FormControl(false);
  checkboxControl: FormControl = new FormControl(false);
  formRadioControl: FormControl = new FormControl(false);

  //load and err
  isCompleted: boolean = false;
  loadError: string = "";
  error: string = ""; //used for adding and deleting record

  //validation

  //annoucement
  contentCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);

  //splunk host
  roleCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  envCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);

  //common
  nameCtl = new FormControl("user_", [Validators.required, Validators.minLength(1)]);

  // user
  groupCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  appTeamCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  serviceIDCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  emailCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  passwordCtl = new FormControl("", [Validators.required, Validators.minLength(4)]);
  emailForEmgCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);

  //splunk user
  positionIDCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  // editPositionIDCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  rpaasNameCtl = new FormControl("");
  memoCtl = new FormControl("");
  splunkPasswordCtl = new FormControl("", [Validators.required, Validators.minLength(4)]);

  //file input
  logFilePathCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  sourcetypeCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  fileSourcetypeCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  logFileSizeCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  dataRetentionCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  // memoCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  blacklistCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  //script input
  scriptSourcetypeCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  scriptIntervalCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  unixAppIntervalCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  optionCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  fileCtl = new FormControl("");

  //unit price
  serviceCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);
  storageCtl = new FormControl("", [Validators.required, Validators.minLength(1)]);

  requiredMsg: string = "Required";
  emailMsg: string = "Must contain @";
  passwordMsg: string = "Length must be greater than 4";
  commonRouter: Router;
  constructor(public messageService: MessageService, public fb: FormBuilder, public service: any, public resourceName: string){
    messageService._rx.subscribe((data: any) => {
      if (typeof data === "string") {
        data = {
          text: data
        };
      }
      this.messages.push(data);
      if (this.duration > 0) {
        setTimeout(() => {
          this.messages.shift();
        }, this.duration);
      }
    });
  }

  ngOnInit(): void {
    this.isCompleted = true; //if server is down, isCompleted cannot be set to true in getRecords
  }

  _uploadFile(service, identifier, obj: any, callback : (data: any) => any, notification=false, resource?:any, param?: any[]){
    this.error = "";
    this.isCompleted = false;
    // this.timeOut();
    service.uploadFile(obj, param)
      .subscribe(
      res => {
        // this.msg = com.createMsg("ok", "Added", this.resourceName, identifier);
        // if(resource){
        //   this.msg = com.createMsg("ok", "Added", resource, identifier);
        // }
        // if(notification)
        //   this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        // this.error = "Failed to add " + this.resourceName + ". " + com.getErrMsg(err);
        // if(resource){
        //   this.error = "Failed to add " + resource + ". " + com.getErrMsg(err);
        // }
        this.isCompleted = true;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
        // document.getElementById("add-modal-close").click();
      });
  }

  _resetPassword(service, identifier, obj: any, callback : (data: any) => any, notification=false, resource?:any, param?: any[]){
    this.error = "";
    this.isCompleted = false;
    this.timeOut();
    service.resetPassword(obj, param)
      .subscribe(
      res => {
        this.msg = com.createMsg("ok", "Reset password", this.resourceName, identifier);
        if(resource){
          this.msg = com.createMsg("ok", "Reset password", resource, identifier);
        }
        if(notification)
          this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        this.error = "Failed to reset password. " + com.getErrMsg(err);
        if(resource){
          this.error = "Failed to reset password. " + com.getErrMsg(err);
        }
        this.isCompleted = true;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
        if (document.getElementById("reset-password-modal-close") != null) {
          document.getElementById("reset-password-modal-close").click();
        }
      });
  }

  _addRecord(service, identifier, obj: any, callback : (data: any) => any, notification=true, resource?:any, param?: any[]){
    this.error = "";
    this.isCompleted = false;
    // this.timeOut();
    service.addRecord(obj, param)
      .subscribe(
      res => {
        this.msg = com.createMsg("ok", "Added", this.resourceName, identifier);
        if(resource){
          this.msg = com.createMsg("ok", "Added", resource, identifier);
        }
        if(notification)
          this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        this.error = "Failed to add " + this.resourceName + ". " + com.getErrMsg(err);
        if(resource){
          this.error = "Failed to add " + resource + ". " + com.getErrMsg(err);
        }
        this.isCompleted = true;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
        if (document.getElementById("add-modal-close") != null) {
          document.getElementById("add-modal-close").click();
        }
      });
  }

  _createDeploymentApp(service, identifier, obj: any, callback : (data: any) => any, notification=true, resource?:any, param?: any[]){
    this.error = "";
    this.isCompleted = false;
    this.timeOut();
    service.createDeploymentApp(obj, param)
      .subscribe(
      res => {
        // this.msg = com.createMsg("ok", "Added", this.resourceName, identifier);
        // if(resource){
        //   this.msg = com.createMsg("ok", "Added", resource, identifier);
        // }
        // if(notification)
        //   this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        var reason = com.getErrMsg(err);
        this.msg = com.createMsg("error", "create", this.resourceName, identifier, reason);
        if (resource)
          this.msg = com.createMsg("error", "create", resource, identifier, reason);

        if(notification)
          this.messageService.emitMessage(this.msg);

        this.isCompleted = true;
        this.isEditing = false;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
        document.getElementById("add-modal-close").click();
      });
  }

  _deinstallDeploymentApp(service, identifier, obj: any, callback : (data: any) => any, notification=true, resource?:any, param?: any[]){
    this.error = "";
    this.isCompleted = false;
    this.timeOut();
    service.deinstallDeploymentApp(obj, param)
      .subscribe(
      res => {
        this.msg = com.createMsg("ok", "Deinstalled", this.resourceName, identifier);
        if(resource){
          this.msg = com.createMsg("ok", "Deinstalled", resource, identifier);
        }
        if(notification)
          this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        var reason = com.getErrMsg(err);
        this.msg = com.createMsg("error", "deinstall", this.resourceName, identifier, reason);
        if (resource)
          this.msg = com.createMsg("error", "deinstall", resource, identifier, reason);

        if(notification)
          this.messageService.emitMessage(this.msg);

        this.isCompleted = true;
        this.isEditing = false;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
      });
  }

  _createDeploymentServerClass(service, identifier, obj: any, callback : (data: any) => any, notification=true, resource?:any, param?: any[]){
    this.error = "";
    this.isCompleted = false;
    this.timeOut();
    service.createDeploymentServerClass(obj, param)
      .subscribe(
      res => {
        this.msg = com.createMsg("ok", "Added", this.resourceName, identifier);
        if(resource){
          this.msg = com.createMsg("ok", "Added", resource, identifier);
        }
        if(notification)
          this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        var reason = com.getErrMsg(err);
        this.msg = com.createMsg("error", "create", this.resourceName, identifier, reason);
        if (resource)
          this.msg = com.createMsg("error", "create", resource, identifier, reason);

        if(notification)
          this.messageService.emitMessage(this.msg);

        this.isCompleted = true;
        this.isEditing = false;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
      });
  }

  _getRecords(service: any, callback : (data: any) => any, param?: any[], resource?: string){
    this.isCompleted = false;
    this.loadError = "";
    service.getRecords(param).subscribe(
      res => {
        this.records = res;
        if(res === null){
          this.records = [];//make sure records.lenght == 0, if not set, records won't be empty even if result is null
        }
        callback(this.records);
      },
      err => {
        this.loadError = "Failed to load " + this.resourceName + "s. " + com.getErrMsg(err);
        if(resource){
          this.loadError = "Failed to load " + this.resourceName + "s. " + com.getErrMsg(err);
        }
        this.isCompleted = true;
        this.setExpiredJWT(err);
      },
      () => {
        this.isCompleted = true;
      });
  }

  _getRecord(service: any, identifier: any, callback : (data: any) => any, param?: any, resource?: string){
    service.getRecord(identifier, param).subscribe(
      res => {
        this.record = res;
        callback(this.record);
      },
      err => {
        this.loadError = "Failed to load " + this.resourceName + ". " + com.getErrMsg(err);
        if(resource){
          this.loadError = "Failed to load " + resource + ". " + com.getErrMsg(err);
        }
        this.isCompleted = true;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
      });
  }

  _editRecord(service: any, identifier:any, objOld, objNew: any, callback : (data: any) => any, notification=true, param?: any, resource?: any){
    this.timeOut();
    this.isCompleted = false;
    // com.submitEdit(this.activeTD);
    service.editRecord(objOld, objNew, param)
      .subscribe(
      res => {
        this.msg = com.createMsg("ok", "Updated", this.resourceName, identifier);
        if (resource)
          this.msg = com.createMsg("ok", "Updated", resource, identifier);

        if(notification)
          this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        var reason = com.getErrMsg(err);
        this.msg = com.createMsg("error", "update", this.resourceName, identifier, reason);
        if (resource)
          this.msg = com.createMsg("error", "update", resource, identifier, reason);

        if(notification)
          this.messageService.emitMessage(this.msg);

        this.isCompleted = true;
        this.isEditing = false;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
        this.isEditing = false;
      });
  }

  _deleteRecord(service, identifier: any, callback : (data: any) => any, notification=true, param?: any, resource?: any){
    this.timeOut();
    this.isCompleted = false;
    service.deleteRecord(identifier, param)
      .subscribe(
      res => {
        var rsc = resource? resource: this.resourceName;
        this.msg = com.createMsg("ok", "Deleted", rsc, identifier);
        // if(param){
        //   this.msg = com.createMsg("ok", "Deleted", rsc, identifier + " " + param);
        // }
        if(notification)
          this.messageService.emitMessage(this.msg);
        callback("ok");
      },
      err => {
        var reason = com.getErrMsg(err);
        var rsc = resource? resource: this.resourceName;
        this.msg = com.createMsg("error", "delete", rsc, identifier, reason);
        // if(param){
        //   this.msg = com.createMsg("error", "delete", rsc, identifier + " " + param, reason);
        // }
        if(notification)
          this.messageService.emitMessage(this.msg);
        this.isCompleted = true;
        this.setExpiredJWT(err);
        callback("error");
      },
      () => {
        this.isCompleted = true;
      });
  }

  setExpiredJWT(err){
    if(err.status == 401)
      localStorage.setItem("jwt", "expired");
  }

  sortTable(attr: string) {
    this.activeAttr = attr;
    if (this.orderAttr[0] == "-") { //1,2,3
      this.orderAttr = "+" + attr;
      this.arrow = "sort ascending";
    } else { //3,2,1
      this.orderAttr = "-" + attr;
      this.arrow = "sort descending";
    }
  }

  timeOut(){
    setTimeout(() => {
      this.isCompleted = true;
    }, 10000);
  }

  getItemNum(){
    return (this.curentP - 1) * this.itemNumPerPage;
  }

  getSubItemNum(){
    return (this.subCurrentP - 1) * this.subItemNumPerPage;
  }

  openAddModal(){
    this.error = "";
    this.isCompleted = true;
    this.form.reset();
  }

  submitAddition(){
    this.timeOut();
    this.error = "";
    this.isCompleted = false;
  }

  openDelModal(obj){
  }

  submitDeletion(){
    this.timeOut();
    this.error = "";
    this.isCompleted = false;
  }

  //improved
  onEdit(row, col, obj){
    this.isEditing = true;
    this.eleActive[row][col] = true;
    this.eleHovered[row][col] = false;
  }

  submitEdit(){
    this.isEditing = false;
    this.eleActive[this.activeCell[0]][this.activeCell[1]] = false;
    this.eleHovered[this.activeCell[0]][this.activeCell[1]] = false;
    this.activeCell = [-1, -1];
  }

  cancelEdit(){
    this.isEditing = false;
    this.eleActive[this.activeCell[0]][this.activeCell[1]] = false;
    this.activeCell = [-1, -1];
  }

  mouseOver(row, col:number){
    if (this.isEditing)
      return;
    if (this.activeCell[0] >= 0 && this.activeCell[1] >= 0 && 
      this.eleActive[this.activeCell[0]] != null && 
      this.eleActive[this.activeCell[0]][this.activeCell[1]] != null) { // for second, third tables in one view
      this.eleActive[this.activeCell[0]][this.activeCell[1]] = false;
    }
    this.activeCell = [row, col];
    this.eleHovered[row][col] = true;
  }
   
  mouseLeave(){
    if (this.isEditing) {
      return;
    }
    if (this.activeCell[0] < 0 || this.activeCell[1] < 0) {
      return
    }
    this.eleActive[this.activeCell[0]][this.activeCell[1]] = false;
    this.eleHovered[this.activeCell[0]][this.activeCell[1]] = false;
  }
}
