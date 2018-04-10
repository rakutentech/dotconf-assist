import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { UserService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { User, Message } from '../model';
import * as com from "../common";
import * as _ from 'lodash';
declare var $;

@Component({
  selector: 'user',
  templateUrl: '../../assets/templates/user.component.html',
  providers: [UserService, MessageService]
})

export class UserComponent extends CommonComponent{
  admin: string = ""; // used for select ele
  users: User[];
  currentUser: User = new User();
  oldUser: User = new User();
  constructor(public messageService: MessageService, public fb: FormBuilder, private userService: UserService){
    super(messageService, fb, userService, "user");
    this.form = fb.group({
      "userName": this.nameCtl,
      "groupName": this.groupCtl,
      "appTeamName": this.appTeamCtl,
      "serviceID": this.serviceIDCtl,
      "email": this.emailCtl,
      "emailForEmg": this.emailForEmgCtl,
      "password": this.passwordCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.getRecords();
  }

  getRecords(){
    var callback = (result: any[]) : void => {
      this.users = result;
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.users.length; i++){
        this.users[i].last_login_at = com.formatTime(this.users[i].last_login_at);
        this.users[i].created_at = com.formatTime(this.users[i].created_at);
        this.users[i].updated_at = com.formatTime(this.users[i].updated_at);
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 5; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.userService, callback);
  }

  addRecord(){
    this.currentUser.admin = false;
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._addRecord(this.userService, this.currentUser.user_name, this.currentUser, callback);
  }
  
  editRecord(param?:any){
    super.submitEdit();
    this.currentUser.admin = this.admin == "true"? true: false;
    var callback = (result: any) : void => {
      if(result == "ok") {
        for( var i = 0; i < this.users.length; i++){
          if (this.users[i].user_name == this.oldUser.user_name)
          {
            // this.users[r + this.itemNumPerPage * (this.curentP -1)] = _.cloneDeep(this.currentUser); //if don't clone, when click add user, the record editted last time will be set as empty
            this.users[i] = _.cloneDeep(this.currentUser); //if don't clone, when click, the record editted last time will be set as empty
          }
        }
      }
    }
    super._editRecord(this.userService, this.currentUser.user_name, this.oldUser, this.currentUser, callback, true, param);
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._deleteRecord(this.userService, this.currentUser.user_name, callback);
  }

  openAddModal(){
    super.openAddModal();
    this.checkboxControl = new FormControl(false);
  }

  openDelModal(obj){
    this.currentUser = _.cloneDeep(obj);
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.admin = obj.admin? "true": "false";
    this.currentUser = _.cloneDeep(obj);
    this.oldUser = _.cloneDeep(obj);
  }
}
