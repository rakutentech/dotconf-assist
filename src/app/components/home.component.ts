import { Component, OnInit} from "@angular/core";
import { Router, ActivatedRoute } from "@angular/router";
import { CommonComponent } from "./common.component";
import { AnnouncementService, UserService, SplunkUserService, SplunkHostService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Announcement, Message, User, SplunkUser, Forwarder } from '../model';
import * as com from "../common";
import * as _ from 'lodash';
import * as conf from "../configuration";

@Component({
  selector: 'home',
  templateUrl: '../../assets/templates/home.component.html',
  providers: [AnnouncementService, MessageService, UserService, SplunkUserService, SplunkHostService]
})

export class HomeComponent extends CommonComponent implements OnInit{
  user: User = new User();
  currentUser: User = new User();
  oldUser: User = new User();
  splunkUser: SplunkUser = new SplunkUser();
  currentSplunkUser: SplunkUser = new SplunkUser();
  oldSplunkUser: SplunkUser = new SplunkUser();
  announcements: Announcement[];
  splunkHosts: Forwarder[];
  appName: string = conf.APP_NAME;
  popupPosiIDStr: string = "Position ID is used for Single Sign-On, all the members with the same position ID will have same permission on Splunk. Use comma to seperate if you want to input multiple IDs. e.g: ra1234,ra1235";
  popupRPaaSStr: string = "If Splunk is used for your application on RPaaS, please write your RPaaS user name";
  env: string = localStorage.getItem("env");
  userName: string = localStorage.getItem("user");
  jwt: string = localStorage.getItem("jwt");
  decodedJwt: string;

  staffLink : string = conf.STAFF_LINK;
  appTeamLink : string = conf.APP_TEAM_LINK;
  serviceIDLink : string = conf.SERVICE_ID_LINK;

  constructor(
    private route:ActivatedRoute, 
    public messageService: MessageService, 
    public fb: FormBuilder, 
    private announcementService: AnnouncementService,
    private userService: UserService,
    private splunkUserService: SplunkUserService,
    private splunkHostService: SplunkHostService,
    private router: Router
    ){
    super(messageService, fb, announcementService, "user");
    this.user = new User();
    this.currentUser = new User();
    this.oldUser = new User();
    this.getAnnouncements();
    this.getUser();
    this.route.queryParams.subscribe(val => { //to refresh page when switching env
      super.ngOnInit();
      this.splunkUser = new SplunkUser();
      this.currentSplunkUser = new SplunkUser();
      this.oldSplunkUser = new SplunkUser();
      this.env = localStorage.getItem("env");
      if (this.userName != "admin") {
        this.getSplunkUser();
        this.getSearchHeads();
      }
    });

    this.form = fb.group({ //splunk user
      "name": this.nameCtl,
      "password": this.splunkPasswordCtl,
      "positionID": this.positionIDCtl,
      // "editPositionID": this.editPositionIDCtl,
      "rpaasName": this.rpaasNameCtl,
      "memo": this.memoCtl,
    });

    this.userForm = fb.group({ // user
      "password": this.passwordCtl,
      "groupName": this.groupCtl,
      "appTeamName": this.appTeamCtl,
      "serviceID": this.serviceIDCtl,
      "email": this.emailCtl,
      "emailForEmg": this.emailForEmgCtl,
    });

    // this.jwt = localStorage.getItem("jwt");
    if (this.jwt != "" && this.jwt != null){
      var base64Url = this.jwt.split('.')[1];
      var base64 = base64Url.replace('-', '+').replace('_', '/');
      this.decodedJwt = JSON.parse(window.atob(base64));
    }
  }

  ngOnInit(): void { 
    this.eleActive = [];
    this.eleHovered = [];
    this.eleActive.push(new Array<boolean>());
    this.eleHovered.push(new Array<boolean>());
    for(var j = 0; j < 9; j++){
      this.eleActive[0].push(false);
      this.eleHovered[0].push(false);
    }
  }

  getAnnouncements(){
    var callback = (result: any[]) : void => {
      this.announcements = result;
      for (var i = 0; i < this.announcements.length; i++){
        this.announcements[i].created_at = com.formatTime(this.announcements[i].created_at);
      }
    }
    super._getRecords(this.announcementService, callback);
  }

  getSearchHeads(){
    var callback = (result: any[]) : void => {
      this.splunkHosts = result;
    }
    super._getRecords(this.splunkHostService, callback, [this.env, "SH"]);
  }

  getUser(){
    var callback = (result: any) : void => {
      if (result != "error") {
        this.user = result;
      }
    }
    super._getRecord(this.userService, localStorage.getItem("user"), callback, "", "user");
  }

  getSplunkUser(){
    var callback = (result: any) : void => {
      this.splunkUser = result;
    }
    super._getRecord(this.splunkUserService, localStorage.getItem("user"), callback, this.env, "splunk user");
  }

  openAddModal(){
    this.currentSplunkUser.user_name = this.user.user_name;
    this.currentSplunkUser.email = this.user.email;
    this.currentSplunkUser.env = localStorage.getItem("env");
    if(this.splunkHosts.length > 0)
      this.currentSplunkUser.search_head = this.splunkHosts[0].name;

    this.form.controls['name'].disable();
    this.memoCtl.reset();
    this.rpaasNameCtl.reset();
    this.positionIDCtl.reset();
    this.splunkPasswordCtl.reset();
    this.error = "";
    this.isCompleted = true;
    // super.openAddModal(); //if call this, form will be reset
  }

  submitCreateSplunkUser(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getSplunkUser();
        // this.router.navigate(['/forwarders']);
        // need to trigger app.component to refresh side view
      }
    }
    super._addRecord(this.splunkUserService, this.currentSplunkUser.user_name, this.currentSplunkUser, callback, true, "Splunk user");
  }

  editRecord(param?:any){
    super.submitEdit();
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.user = _.cloneDeep(this.currentUser);
      }
    }
    super._editRecord(this.userService, this.user.user_name, this.oldUser, this.currentUser, callback, true, param);
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.currentUser = _.cloneDeep(obj);
    this.currentUser.password = "";
    this.oldUser = _.cloneDeep(obj);
  }

  editSplunkUserRecord(param?:any){
    this.isEditing = false;
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.splunkUser = _.cloneDeep(this.currentSplunkUser);
      }
    }
    super._editRecord(this.splunkUserService, this.splunkUser.user_name, this.oldSplunkUser, this.currentSplunkUser, callback, true, param, "Splunk user");
  }

  onEditSplunkUser(row, col, obj){
    super.onEdit(row, col, obj);
    this.currentSplunkUser = _.cloneDeep(obj);
    this.oldSplunkUser = _.cloneDeep(obj);
  }

}
