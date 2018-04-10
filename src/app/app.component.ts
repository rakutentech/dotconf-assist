import { Component, AfterViewChecked, AfterViewInit, AfterContentChecked, ViewChild } from '@angular/core';
import { SessionService, SplunkUserService } from "./services/resource.service";
import { MessageService } from "./services/message";
import { CommonComponent } from "./components/common.component";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Router, RouterStateSnapshot } from '@angular/router';
import { User, SplunkUser, RouterOutlet } from './model';
import * as conf from "./configuration";
declare var $;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  providers: [MessageService, SessionService, SplunkUserService]
})

export class AppComponent extends CommonComponent implements AfterViewChecked, AfterViewInit, AfterContentChecked{
  user = new User();
  splunkUser: SplunkUser = new SplunkUser();
  adminRouterOutlets: RouterOutlet[] = [];
  userRouterOutlets: RouterOutlet[] = [];
  userName : string = "guest";
  appName: string = conf.APP_NAME;
  isLoggedIn = false;
  splunkUserCreated: boolean = false;
  lastSessionUpdatedAt: Date = new Date();
  year: number = 0;
  quickStartLink : string = conf.QUICK_START_LINK;
  webToolLink : string = conf.WEB_TOOL_LINK;
  jiraLink : string = conf.JIRA_LINK;
  aboutLink : string = conf.ABOUT_LINK;
  faqLink : string = conf.FAQ_LINK;

  constructor(public messageService: MessageService, public fb: FormBuilder, 
    private splunkUserService: SplunkUserService,
    private sessionService: SessionService, 
    private router: Router){
    super(messageService, fb, sessionService, "session");
    this.form = fb.group({
      "name": this.nameCtl,
      "password": this.passwordCtl,
    });
  }

  ngOnInit(): void {
    this.year = new Date().getFullYear();
    localStorage.setItem("env", "prod");
    this.router.navigate(['/home']);
  }

  ngAfterViewInit(){
    localStorage.setItem("user", this.userName);
    if (this.userName != 'admin') {
      this.getSplunkUser(this.userName, localStorage.getItem("env"));
    }
    $('.ui.dropdown').dropdown('refresh');
    $('.ui.menu a.item')
      .on('click', function() {
        $(this)
          .addClass('active')
          .siblings()
          .removeClass('active')
        ;
      });
  }

  ngAfterViewChecked() {
    var now: Date = new Date();
    if (localStorage.getItem("jwt") !== null && (now.valueOf() - this.lastSessionUpdatedAt.valueOf() > 600000)) { //validate token every 10 min, 600sec
      this.lastSessionUpdatedAt = new Date();
      this.sessionService.getToken().subscribe(
        res => {
          this.sessionService.updateSession().subscribe(
            res => {
              localStorage.setItem("jwt", res.headers.get('x-auth-token'));
            },
            err => {
            });
        },
        err => { //token expired, delete local storage
          this.removeExpiredLocalStorage();
        });
    }
  }

  ngAfterContentChecked(){ //set loggin and admin
    if (localStorage.getItem("jwt") == "expired"){ //set in common.component
      this.removeExpiredLocalStorage();
      return;
    }
    
    if (localStorage.getItem("jwt") === null) { //not logged in
      this.isLoggedIn = false;
    }else{
      $('.ui.dropdown').dropdown('refresh');
      this.isLoggedIn = true;
      if (localStorage.getItem("user") != this.userName && localStorage.getItem("user")!="admin") { // right after login
        this.getSplunkUser(localStorage.getItem("user"), localStorage.getItem("env"));
      }
      if (localStorage.getItem("user")=="admin") {
        this.setRouterOutlets();
      }
      this.userName = localStorage.getItem("user");
      this.user.user_name = this.userName; //in case user name has been set to admin when opened admin login page
    }

  }

  getSplunkUser(user, env: string){
    var callback = (result: any) : void => {
      var oldValue = this.splunkUserCreated;
      if (result == "error") {
        this.splunkUserCreated = false;
      }else{
        this.splunkUserCreated = true;
      }
      this.userRouterOutlets = [];
      this.adminRouterOutlets = [];
      this.setRouterOutlets();
    }
    super._getRecord(this.splunkUserService, user, callback, env, "splunk user");
  }

  setRouterOutlets(){
    if (!this.isLoggedIn) {
      this.userRouterOutlets = [];
      this.adminRouterOutlets = [];
      return;
    }

    if(this.userName == "admin" && this.adminRouterOutlets.length == 0) { //admin
      this.adminRouterOutlets.push(new RouterOutlet("Home", "/home", "home"));
      this.adminRouterOutlets.push(new RouterOutlet("Forwarders", "/forwarders", "disk outline"));
      this.adminRouterOutlets.push(new RouterOutlet("Server Classes", "/server_classes", "server"));
      this.adminRouterOutlets.push(new RouterOutlet("Inputs", "/inputs", "file text outline"));
      this.adminRouterOutlets.push(new RouterOutlet("Apps", "/apps", "folder open outline"));
      this.adminRouterOutlets.push(new RouterOutlet("Deployment", "/deployment", "asterisk"));
      this.adminRouterOutlets.push(new RouterOutlet("Users", "/users", "users"));
      this.adminRouterOutlets.push(new RouterOutlet("Announcements", "/announcements", "announcement"));
      this.adminRouterOutlets.push(new RouterOutlet("Splunk Hosts", "/splunk_hosts", "server"));
      this.adminRouterOutlets.push(new RouterOutlet("Unit Price", "/unit_price", "yen"));
      this.userRouterOutlets = [];
    }

    if (this.userName != "admin" && this.userRouterOutlets.length <= 1){ //non admin
      if (this.userRouterOutlets.length == 0) {
        this.userRouterOutlets.push(new RouterOutlet("Home", "/home", "home"));
      }
      if (this.splunkUserCreated) {
        this.userRouterOutlets.push(new RouterOutlet("Forwarders", "/forwarders", "disk outline"));
        this.userRouterOutlets.push(new RouterOutlet("Server Classes", "/server_classes", "server"));
        this.userRouterOutlets.push(new RouterOutlet("Inputs", "/inputs", "file text outline"));
        this.userRouterOutlets.push(new RouterOutlet("Apps", "/apps", "folder open outline"));
        this.userRouterOutlets.push(new RouterOutlet("Deployment", "/deployment", "asterisk"));
        if (localStorage.getItem("env") == "prod") {
          this.userRouterOutlets.push(new RouterOutlet("Usage (Beta)", "/usage", "yen"));
        }
      }
      this.adminRouterOutlets = [];
    }
  }

  removeExpiredLocalStorage(){
    localStorage.clear();
    this.isLoggedIn = false;
    this.userName = "guest";
    this.router.navigate(['/login']);
    this.adminRouterOutlets = [];
    this.userRouterOutlets = [];
    $('.ui.menu a.item.stg')
          .removeClass('active')
        ;
    $('.ui.menu a.item.dev')
          .removeClass('active')
        ;
    $('.ui.menu a.item.prod')
          .removeClass('active')
          .addClass('active')
        ;
  }

  submitLogout(){
    this.removeExpiredLocalStorage();
  }

  redirectToLogin(){
    this.router.navigate(['/login']);
  }

  switchMode(env){
    localStorage.setItem("env", env);
    if (this.userName != "admin") {
      this.getSplunkUser(this.userName, localStorage.getItem("env"));
    }
    this.router.navigate(['/home'], { queryParams: { 'env': env}});

    if(!this.isLoggedIn) {
      this.router.navigate(['/login']);
    }
  }

  openAdminLoginModal(){
    this.user.user_name = "admin";
    this.isCompleted = true;
    this.form.controls['name'].disable();
  }

  submitAdminLogin(){
    if (!this.form.valid || !this.isCompleted) {
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
        if (err.status == 401)
          this.error = "Incorrect username or password";  //err.status is 401
        else
          this.error = "Sorry, server is scheduled for maintenance";  //err.status isnot 401
        this.isCompleted = true;
      },
      () => {
        this.isCompleted = true;
        document.getElementById("login-modal-close").click();
      });
  }
}