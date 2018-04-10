import { Component, AfterViewChecked, AfterViewInit, AfterContentChecked } from '@angular/core';
import { SessionService, UserService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { CommonComponent } from "../components/common.component";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Router, RouterStateSnapshot } from '@angular/router';
import { User} from '../model';
import * as conf from "../configuration";
declare var $;

@Component({
  selector: 'login',
  templateUrl: '../../assets/templates/login.component.html',
  providers: [MessageService, SessionService, UserService]
})

export class LoginComponent extends CommonComponent implements AfterViewInit {
  user = new User();
  popupUserStr: string = "User name will be used as Splunk account as well, any of group, service or project name is preferred. Personal name is not accepted";
  popupAppStr: string = "If you don't have app team name, input sys team";
  popupSIDStr: string = "Service ID is used for cost allcoation";
  popupEmgStr: string = "Emergency contact when Splunk is down or forwarder is down";

  appTeamLink : string = conf.APP_TEAM_LINK;
  serviceIDLink : string = conf.SERVICE_ID_LINK;

  constructor(public messageService: MessageService, 
    public fb: FormBuilder, 
    private sessionService: SessionService,  
    private userService: UserService,  
    private router: Router){
    super(messageService, fb, sessionService, "session");
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

  ngAfterViewInit(){
    $('.ui.dropdown').dropdown('refresh');
    setTimeout(() => {
      this.user.user_name = "user_";
    });
  }

  submitLogin(){
    if (!this.nameCtl.valid || !this.passwordCtl.valid || !this.isCompleted) {
      return;
    }
    super.timeOut();
    this.error = "";
    this.isCompleted = false;
    this.sessionService.login(this.user.user_name, this.user.password, conf.API_HOST)
      .subscribe(
      res => { //login seccessful
        localStorage.setItem("jwt", res.headers.get('x-auth-token'));
        localStorage.setItem("user", this.user.user_name);
        localStorage.setItem("env", "prod");
        this.router.navigate(['/home']); //can't add other statement after this line, otherwise it won't redirect sometimes
      },
      err => {
        if (err.status == 401){
          if( err._body.includes("has not been approved"))
            this.error = "Your account has not been approved yet";
          else
          this.error = "Incorrect username or password";  //err.status is 401
        }
        else
          this.error = "Sorry, server is scheduled for maintenance";  //err.status isnot 401
        this.isCompleted = true;
      },
      () => {
        this.isCompleted = true;
        // document.getElementById("login-modal-close").click();
      });
  }

  onUsernameKeyUp(){
    if(!this.user.user_name.startsWith('user_')){
      this.user.user_name = "user_";
    }
  }

  submitSignup(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.router.navigate(['/home']); //can't add other statement after this line, otherwise it won't redirect sometimes
      }
    }
    super._addRecord(this.userService, this.user.user_name, this.user, callback, true, "user");
  }

  submitResetPassword(){
    var callback = (result: any) : void => {
      if(result == "ok") {
      }
    }
    super._resetPassword(this.userService, this.user.user_name, this.user, callback, true, "user");
  }
}