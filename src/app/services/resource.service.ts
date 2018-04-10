import {Inject, Injectable} from '@angular/core';
import {User, SplunkUser, Announcement, SplunkHost, Forwarder, ServerClass, 
  FileInput, ScriptInput, UnixAppInput, Deployment, 
  ServerClassForwarder, App, DeploymentApp, 
  DeploymentServerClass, UnitPrice, StorageSize, LogSize} from '../model';
import {Http, Headers, Request, RequestOptions, Response, RequestMethod, ResponseContentType} from '@angular/http';
import * as conf from "../configuration";
import 'rxjs/add/operator/map';
import {Observable} from 'rxjs/Rx';
import * as FileSaver from 'file-saver';

const userUrl = "/v1/users";
const splunkUserUrl = "/v1/splunk_users";
const announcementUrl = "/v1/announcements";
const forwarderUrl = "/v1/forwarders";
const serverClassUrl = "/v1/server_classes";
const inputUrl = "/v1/inputs";
const appUrl = "/v1/apps";
const deploymentUrl = "/v1/deployment";
const splunkHostUrl = "/v1/splunk_hosts";
const unitPriceUrl = "/v1/unit_price";
const usageUrl = "/v1/usage";
const loginUrl = "/v1/login";
const logoutUrl = "/v1/logout";
const updateSessionUrl = "/v1/update_session";
const tokenUrl = "/v1/tokens";

function getDataHeader(isArray: boolean): string {
  if (isArray)
    return '[{';
  else
    return '{';
}

function getDataPair(key: string, value: any, isNumber: boolean): string {
  if (isNumber)
    return '"' + key + '":' + value + '';
  else
    return '"' + key + '":"' + value + '"';
}

function getDataTail(isArray: boolean): string {
  if (isArray)
    return '}]';
  else
    return '}';
}

function getRequestOption(method: RequestMethod, url: string, data: any): RequestOptions {
  var requestHeaders = new Headers();
  requestHeaders.append('x-auth-token', localStorage.getItem("jwt"));
  var options = new RequestOptions({
    method: method,
    url: conf.API_HOST + url,
    headers: requestHeaders,
    body: data
  });
  return options;
}

function getBlobRequestOption(method: RequestMethod, url: string, data: any): RequestOptions {
  var requestHeaders = new Headers();
  requestHeaders.append('x-auth-token', localStorage.getItem("jwt"));
  var options = new RequestOptions({
    method: method,
    url: conf.API_HOST + url,
    headers: requestHeaders,
    body: data,
    responseType: ResponseContentType.Blob
  });
  return options;
}

@Injectable()
export class SessionService {
  constructor(
    private http:Http
    ) { 
  }
  login(userName:string, pwd:string, host:string){
    var data = getDataHeader(false) + getDataPair("user_name", userName, false) + "," + getDataPair("password", pwd, false) + getDataTail(false);
    var options = new RequestOptions({
      method: RequestMethod.Post,
      url: host + loginUrl,
      body: data
    });
    return this.http.request(new Request(options));
  }
  updateSession() {
    var options = getRequestOption(RequestMethod.Post, updateSessionUrl, "");
    return this.http.request(new Request(options));
  }

  getToken() {
    var options = getRequestOption(RequestMethod.Get, tokenUrl, "");
    return this.http.request(new Request(options));
  }
}

@Injectable()
export class UserService {
  constructor(
    private http:Http
    ) { }
  getRecords(){
    var options = getRequestOption(RequestMethod.Get, userUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  getRecord(userName: string){
    var options = getRequestOption(RequestMethod.Get, userUrl + "/" + userName, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  resetPassword(user: User){
    var data = getDataHeader(false) 
      + getDataPair("user_name", user.user_name, false) + ","
      + getDataPair("email", user.email, false) + ","
      + getDataPair("password", user.password, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, userUrl + "/reset_password", data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  addRecord(user: User){
    var data = getDataHeader(false) 
      + getDataPair("user_name", user.user_name, false) + ","
      + getDataPair("group_name", user.group_name, false) + ","
      + getDataPair("app_team_name", user.app_team_name, false) + ","
      + getDataPair("service_id", user.service_id, false) + ","
      + getDataPair("email", user.email, false) + ","
      + getDataPair("email_for_emergency", user.email_for_emergency, false) + ","
      + getDataPair("admin", false, true) + ","
      + getDataPair("password", user.password, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, userUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(oldUser, newUser: User, action?: string) {
    var data = getDataHeader(false)
      + getDataPair("user_name", newUser.user_name, false) + ","
      + getDataPair("group_name", newUser.group_name, false) + ","
      + getDataPair("app_team_name", newUser.app_team_name, false) + ","
      + getDataPair("service_id", newUser.service_id, false) + ","
      + getDataPair("email", newUser.email, false) + ","
      + getDataPair("email_for_emergency", newUser.email_for_emergency, false) + ","
      + getDataPair("admin", false, true) + ","
      + getDataPair("status", newUser.status, false) + ","
      + getDataPair("password", newUser.password, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Put, userUrl + "/" + oldUser.user_name, data);
    if (action != null && action.startsWith("change_")){
      options = getRequestOption(RequestMethod.Put, userUrl + "/" + oldUser.user_name + "?action=" + action, data);
    }
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(user: string) {
    var options = getRequestOption(RequestMethod.Delete, userUrl + "/" + user, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class SplunkUserService {
  constructor(
    private http:Http
    ) { }
  getRecords(){
    var options = getRequestOption(RequestMethod.Get, splunkUserUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  getRecord(userName: string, env: string){
    var options = getRequestOption(RequestMethod.Get, splunkUserUrl + "/" + userName + "?env=" + env, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(user: SplunkUser){
    var data = getDataHeader(false) 
      + getDataPair("user_name", user.user_name, false) + ","
      + getDataPair("password", user.password, false) + ","
      + getDataPair("rpaas_user_name", user.rpaas_user_name, false) + ","
      + getDataPair("email", user.email, false) + ","
      + getDataPair("env", user.env, false) + ","
      + getDataPair("memo", user.memo, false) + ","
      + getDataPair("position_ids", user.position_ids, false) + ","
      + getDataPair("search_head", user.search_head, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, splunkUserUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
  editRecord(oldUser, newUser: SplunkUser) {
    var data = getDataHeader(false) 
      // + getDataPair("user_name", newUser.user_name, false) + ","
      // + getDataPair("password", newUser.password, false) + ","
      // + getDataPair("search_head", newUser.search_head, false) + ","
      // + getDataPair("env", newUser.env, false) + ","
      + getDataPair("rpaas_user_name", newUser.rpaas_user_name, false) + ","
      + getDataPair("memo", newUser.memo, false) + ","
      + getDataPair("position_ids", newUser.position_ids, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Put, splunkUserUrl + "/" + oldUser.user_name, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
  deleteRecord(user: string) {
    var options = getRequestOption(RequestMethod.Delete, splunkUserUrl + "/" + user, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class AnnouncementService {
  constructor(
    private http:Http
    ) { }

  getRecords(){
    var options = getRequestOption(RequestMethod.Get, announcementUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(announcement: Announcement){
    var re = /"/gi;
    var content = announcement.content.replace(re, "\\\"");
    var data = getDataHeader(false)
      + getDataPair("content", content, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, announcementUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(oldAnnouncement, newAnnouncement: Announcement) {
    var re = /"/gi;
    var content = newAnnouncement.content.replace(re, "\\\"");
    var data = getDataHeader(false)
      + getDataPair("content", content, false)
      + getDataTail(false);

    var options = getRequestOption(RequestMethod.Put, announcementUrl + "/" + newAnnouncement.id, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(id: number) {
    var options = getRequestOption(RequestMethod.Delete, announcementUrl + "/" + id, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class SplunkHostService {
  constructor(
    private http:Http
    ) { }

  getRecords(envAndRole: string[]){
    var options = getRequestOption(RequestMethod.Get, splunkHostUrl + "?env=" + envAndRole[0] + "&role=" + envAndRole[1], "");
    // var options = getRequestOption(RequestMethod.Get, splunkHostUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(splunkHost: SplunkHost){
    var data = getDataHeader(false)
      + getDataPair("name", splunkHost.name, false) + "," 
      + getDataPair("role", splunkHost.role, false) + "," 
      + getDataPair("env", splunkHost.env, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, splunkHostUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(oldSplunkHost, newSplunkHost: SplunkHost) {
    var data = getDataHeader(false)
      + getDataPair("name", newSplunkHost.name, false) + "," 
      + getDataPair("role", newSplunkHost.role, false) + "," 
      + getDataPair("env", newSplunkHost.env, false)
      + getDataTail(false);

    var options = getRequestOption(RequestMethod.Put, splunkHostUrl + "/" + oldSplunkHost.name, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(name: string) {
    var options = getRequestOption(RequestMethod.Delete, splunkHostUrl + "/" + name, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class ForwarderService {
  constructor(
    private http:Http
    ) { }

  getRecords(envUserFrom: string[]){
    var options = getRequestOption(RequestMethod.Get, forwarderUrl + "?env=" + envUserFrom[0] + "&user=" + envUserFrom[1] + "&from=" + envUserFrom[2], "");
    // var options = getRequestOption(RequestMethod.Get, forwarderUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(forwarders: Forwarder[]){
    var data = "[";
    for (var i = 0; i < forwarders.length; i++){
      if (i != 0) {
        data += ',';
      }
      data += getDataHeader(false);
      data += getDataPair("name", forwarders[i].name, false) + ",";
      data += getDataPair("user_name", forwarders[i].user_name, false) + ",";
      data += getDataPair("share", forwarders[i].share, false) + ","; //share is set as "private" when adding forwarder
      data += getDataPair("env", forwarders[i].env, false);
      data += getDataTail(false);
    }
    data += "]";
    var options = getRequestOption(RequestMethod.Post, forwarderUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(oldForwarder, newForwarder: Forwarder) {
    var data = getDataHeader(false)
      + getDataPair("name", newForwarder.name, false) + "," 
      + getDataPair("user_name", newForwarder.user_name, false) + "," 
      + getDataPair("share", newForwarder.share, false) + "," 
      + getDataPair("env", newForwarder.env, false)
      + getDataTail(false);

    var options = getRequestOption(RequestMethod.Put, forwarderUrl + "/" + oldForwarder.name + "?user=" + oldForwarder.user_name, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(name, user: string) {
    var options = getRequestOption(RequestMethod.Delete, forwarderUrl + "/" + name + "?user=" + user, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class ServerClassService {
  constructor(
    private http:Http
    ) { }

  getRecords(envUser: string[]){
    var options = getRequestOption(RequestMethod.Get, serverClassUrl + "?env=" + envUser[0] + "&user=" + envUser[1], "");
    // var options = getRequestOption(RequestMethod.Get, serverClassUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(serverClass: ServerClass){
    var forwarders = '"forwarders":[';
    for (var i = 0; i < serverClass.forwarder_ids.length; i++){
      if (i != 0) {
        forwarders += ',';
      }
      forwarders += getDataHeader(false);
      forwarders += getDataPair("id", serverClass.forwarder_ids[i], true);
      forwarders += getDataTail(false);
    }
    forwarders += "]";

    var data = getDataHeader(false);
    data += getDataPair("name", serverClass.name, false) + ",";
    data += getDataPair("user_name", serverClass.user_name, false) + ",";
    data += getDataPair("env", serverClass.env, false) + ",";
    data += forwarders;
    data += getDataTail(false);

    var options = getRequestOption(RequestMethod.Post, serverClassUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(oldServerClass, newServerClass: ServerClass) {
    var forwarders = '"forwarders":[';
    for (var i = 0; i < newServerClass.forwarder_ids.length; i++){
      if (i != 0) {
        forwarders += ',';
      }
      forwarders += getDataHeader(false);
      forwarders += getDataPair("id", newServerClass.forwarder_ids[i], true);
      forwarders += getDataTail(false);
    }
    forwarders += "]";

    var data = getDataHeader(false);
    data += getDataPair("name", newServerClass.name, false) + ",";
    data += getDataPair("user_name", newServerClass.user_name, false) + ",";
    data += getDataPair("env", newServerClass.env, false) + ",";
    data += forwarders;
    data += getDataTail(false);

    var options = getRequestOption(RequestMethod.Put, serverClassUrl + "/" + oldServerClass.name + "?user=" + oldServerClass.user_name, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(name, user: string) {
    var options = getRequestOption(RequestMethod.Delete, serverClassUrl + "/" + name + "?user=" + user, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class InputService {
  constructor(
    private http:Http
    ) { }

  getRecords(params: string[]){
    var options = getRequestOption(RequestMethod.Get, inputUrl + "?env=" + params[0] + "&user=" + params[1] + "&type=" + params[2] + "&app_id=" + params[3], "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(input: any, params: string[]){
    if(params[2] == "file"){
      var re = /\\/gi;
      var blacklist = input.blacklist? input.blacklist:"";
      var logFilePath = input.log_file_path;
      blacklist = blacklist.replace(re, "\\\\");
      logFilePath = logFilePath.replace(re, "\\\\");
      var data = getDataHeader(false)
        + getDataPair("log_file_path", logFilePath, false) + "," 
        + getDataPair("sourcetype", input.sourcetype, false) + "," 
        + getDataPair("log_file_size", input.log_file_size, false) + "," 
        + getDataPair("data_retention_period", input.data_retention_period, false) + "," 
        + getDataPair("memo", input.memo, false) + "," 
        + getDataPair("env", input.env, false) + "," 
        + getDataPair("user_name", input.user_name, false) + "," 
        + getDataPair("app_id", input.app_id, true) + "," 
        + getDataPair("blacklist", blacklist, false) + "," 
        + getDataPair("crcsalt", input.crcsalt, false)
        + getDataTail(false);
      var options = getRequestOption(RequestMethod.Post, inputUrl + "?env=" + params[0] + "&user=" + params[1] + "&type=" + params[2], data);
      return this.http.request(new Request(options)).map((res: Response) => {
      });
    }else if(params[2] == "script"){
      var formData = new FormData();
      formData.append('script', input.script, input.script_name);
      formData.append('sourcetype', input.sourcetype);
      formData.append('log_file_size', input.log_file_size);
      formData.append('data_retention_period', input.data_retention_period);
      formData.append('env', input.env);
      formData.append('user_name', input.user_name);
      formData.append('os', input.os);
      // formData.append('app_id', input.app_id);
      formData.append('interval', input.interval);
      formData.append('exefile', input.exefile);
      formData.append('script_name', input.script_name);
      formData.append('option', input.option);
      var options = getBlobRequestOption(RequestMethod.Post, inputUrl + "?env=" + params[0] + "&user=" + params[1] + "&type=" + params[2], formData);
      return this.http.request(new Request(options)).map((res: Response) => {
      });
    }else if(params[2] == "unixapp"){
      var data = "[";
      for (var i = 0; i < input.length; i++){
        if (i != 0) {
          data += ',';
        }
        data += getDataHeader(false);
        data += getDataPair("script_name", input[i].script_name, false) + ",";
        data += getDataPair("env", input[i].env, false) + ",";
        data += getDataPair("user_name", input[i].user_name, false) + ",";
        data += getDataPair("app_id", input[i].app_id, true) + ",";
        data += getDataPair("interval", input[i].interval, false) + ",";
        data += getDataPair("data_retention_period", input[i].data_retention_period, false);
        data += getDataTail(false);
      }
      data += "]";
      var options = getRequestOption(RequestMethod.Post, inputUrl + "?env=" + params[0] + "&user=" + params[1] + "&type=" + params[2], data);
      return this.http.request(new Request(options)).map((res: Response) => {
      });
    }
  }

  editRecord(oldInput, newInput: any, inputType: string) {
    if(inputType == "file"){
      var re = /\\/gi;
      var blacklist = newInput.blacklist? newInput.blacklist:"";
      var logFilePath = newInput.log_file_path;
      blacklist = blacklist.replace(re, "\\\\");
      logFilePath = logFilePath.replace(re, "\\\\");
      var data = getDataHeader(false)
        + getDataPair("log_file_path", logFilePath, false) + "," 
        + getDataPair("sourcetype", newInput.sourcetype, false) + "," 
        + getDataPair("log_file_size", newInput.log_file_size, false) + "," 
        + getDataPair("data_retention_period", newInput.data_retention_period, false) + "," 
        + getDataPair("memo", newInput.memo, false) + "," 
        // + getDataPair("splunk_user_id", newInput.splunk_user_id, true) + "," 
        + getDataPair("app_id", newInput.app_id, true) + "," 
        + getDataPair("blacklist", blacklist, false) + "," 
        + getDataPair("crcsalt", newInput.crcsalt, false)
        + getDataTail(false);
      var options = getRequestOption(RequestMethod.Put, inputUrl + "/" + oldInput.id + "?type=" + inputType, data);
      return this.http.request(new Request(options)).map((res: Response) => {
      });
    }else if(inputType == "script"){
      var formData = new FormData();
      if(newInput.script != null)
        formData.append('script', newInput.script, newInput.script_name);
      formData.append('sourcetype', newInput.sourcetype);
      formData.append('log_file_size', newInput.log_file_size);
      formData.append('data_retention_period', newInput.data_retention_period);
      formData.append('os', newInput.os);
      formData.append('interval', newInput.interval);
      formData.append('exefile', newInput.exefile);
      formData.append('script_name', newInput.script_name);
      formData.append('option', newInput.option);
      var options = getBlobRequestOption(RequestMethod.Put, inputUrl + "/" + oldInput.id + "?type=" + inputType, formData);
      return this.http.request(new Request(options)).map((res: Response) => {
      });
    }else if(inputType == "unixapp"){
      var data = getDataHeader(false)
        + getDataPair("data_retention_period", newInput.data_retention_period, false) + "," 
        + getDataPair("interval", newInput.interval, false)
        + getDataTail(false);
      var options = getRequestOption(RequestMethod.Put, inputUrl + "/" + oldInput.id + "?type=" + inputType, data);
      return this.http.request(new Request(options)).map((res: Response) => {
      });
    }
  }

  deleteRecord(id: number, inputType: string) {
    var options = getRequestOption(RequestMethod.Delete, inputUrl + "/" + id + "?type=" + inputType, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class AppService {
  constructor(
    private http:Http
    ) { }

  getRecords(envUser: string[]){
    var options = getRequestOption(RequestMethod.Get, appUrl + "?env=" + envUser[0] + "&user=" + envUser[1], "");
    // var options = getRequestOption(RequestMethod.Get, appUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(app: App){
    var fileInputs = "";
    var scriptInputs = "";
    var unixAppInputs = "";
    if( !app.unix_app){
      fileInputs = '"file_input_ids":[';
      for (var i = 0; i < app.file_input_ids.length; i++){
        if (i != 0) {
          fileInputs += ',';
        }
        fileInputs += getDataHeader(false);
        fileInputs += getDataPair("id", app.file_input_ids[i], true);
        fileInputs += getDataTail(false);
      }
      fileInputs += "]";

      scriptInputs = '"script_input_ids":[';
      for (var i = 0; i < app.script_input_ids.length; i++){
        if (i != 0) {
          scriptInputs += ',';
        }
        scriptInputs += getDataHeader(false);
        scriptInputs += getDataPair("id", app.script_input_ids[i], true);
        scriptInputs += getDataTail(false);
      }
      scriptInputs += "]";
    }

    else if (app.unix_app){
      unixAppInputs = '"unix_app_input_ids":[';
      for (var i = 0; i < app.unix_app_input_ids.length; i++){
        if (i != 0) {
          unixAppInputs += ',';
        }
        unixAppInputs += getDataHeader(false);
        unixAppInputs += getDataPair("id", app.unix_app_input_ids[i], true);
        unixAppInputs += getDataTail(false);
      }
      unixAppInputs += "]";
    }

    var data = getDataHeader(false);
    data += getDataPair("name", app.name, false) + ",";
    data += getDataPair("env", app.env, false) + ",";
    data += getDataPair("user_name", app.user_name, false) + ",";
    data += getDataPair("unix_app", app.unix_app, true) + ",";
    if( !app.unix_app){
      data += fileInputs + ",";
      data += scriptInputs;
    }else if (app.unix_app){
      data += unixAppInputs;
    }
    data += getDataTail(false);

    var options = getRequestOption(RequestMethod.Post, appUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(oldApp, newApp: App) {
   var fileInputs = "";
    var scriptInputs = "";
    var unixAppInputs = "";
    if( !oldApp.unix_app){
      fileInputs = '"file_input_ids":[';
      for (var i = 0; i < newApp.file_input_ids.length; i++){
        if (i != 0) {
          fileInputs += ',';
        }
        fileInputs += getDataHeader(false);
        fileInputs += getDataPair("id", newApp.file_input_ids[i], true);
        fileInputs += getDataTail(false);
      }
      fileInputs += "]";

      scriptInputs = '"script_input_ids":[';
      for (var i = 0; i < newApp.script_input_ids.length; i++){
        if (i != 0) {
          scriptInputs += ',';
        }
        scriptInputs += getDataHeader(false);
        scriptInputs += getDataPair("id", newApp.script_input_ids[i], true);
        scriptInputs += getDataTail(false);
      }
      scriptInputs += "]";
    }

    else{
      unixAppInputs = '"unix_app_input_ids":[';
      for (var i = 0; i < newApp.unix_app_input_ids.length; i++){
        if (i != 0) {
          unixAppInputs += ',';
        }
        unixAppInputs += getDataHeader(false);
        unixAppInputs += getDataPair("id", newApp.unix_app_input_ids[i], true);
        unixAppInputs += getDataTail(false);
      }
      unixAppInputs += "]";
    }

    var data = getDataHeader(false);
    data += getDataPair("name", newApp.name, false) + ",";
    if( !oldApp.unix_app){
      data += fileInputs + ",";
      data += scriptInputs;
    }else{
      data += unixAppInputs;
    }
   
    data += getDataTail(false);

    var options = getRequestOption(RequestMethod.Put, appUrl + "/" + oldApp.id, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(id: number) {
    var options = getRequestOption(RequestMethod.Delete, appUrl + "/" + id, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class DeploymentService {
  constructor(
    private http:Http
    ) { }

  getRecords(envUser: string[]){
    var options = getRequestOption(RequestMethod.Get, deploymentUrl + "?env=" + envUser[0] + "&user=" + envUser[1], "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(data: any, action: string){
  }

  createDeploymentApp(deploymentApp: DeploymentApp){
    var script = '"script_ids":[';
    for (var i = 0; i < deploymentApp.script_ids.length; i++) {
      if (i != 0) {
        script += ',';
      }
     script += getDataHeader(false);
     script += getDataPair("id", deploymentApp.script_ids[i], true);
     script += getDataTail(false);
    }
    script += ']';

    var data = getDataHeader(false)
      + getDataPair("app_id", deploymentApp.app_id, true) + "," 
      + getDataPair("env", deploymentApp.env, false) + "," 
      + getDataPair("app_type", deploymentApp.app_type, false) + "," 
      + getDataPair("folder_name", deploymentApp.folder_name, false) + "," 
      + getDataPair("inputs_conf", deploymentApp.inputs_conf, false) + ","
      + script
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, deploymentUrl + "/create_app", data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  createDeploymentServerClass(deploymentServerClasses: DeploymentServerClass[]){
    var data = "[";
    for (var i = 0; i < deploymentServerClasses.length; i++){

      var forwarders = '"forwarder_names":[';
      for (var j = 0; j < deploymentServerClasses[i].forwarder_names.length; j++) {
        if (j != 0) {
          forwarders += ',';
        }
       forwarders += getDataHeader(false);
       forwarders += getDataPair("name", deploymentServerClasses[i].forwarder_names[j], false);
       forwarders += getDataTail(false);
      }
      forwarders += ']';

      if (i != 0) {
        data += ',';
      }
      data += getDataHeader(false);
      data += getDataPair("server_class_name", deploymentServerClasses[i].server_class_name, false) + ",";
      data += getDataPair("app_name", deploymentServerClasses[i].app_name, false) + ",";
      data += getDataPair("app_id", deploymentServerClasses[i].app_id, true) + ",";
      data += getDataPair("user", deploymentServerClasses[i].user, false) + ",";
      data += getDataPair("env", deploymentServerClasses[i].env, false) + ",";
      data += forwarders;
      data += getDataTail(false);
    }
    data += "]";
    var options = getRequestOption(RequestMethod.Post, deploymentUrl + "/create_server_class", data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deinstallDeploymentApp(deploymentApp: DeploymentApp){
    var data = getDataHeader(false)
      + getDataPair("app_id", deploymentApp.app_id, true) + "," 
      + getDataPair("env", deploymentApp.env, false) + "," 
      + getDataPair("app_type", deploymentApp.app_type, false) + "," 
      + getDataPair("folder_name", deploymentApp.folder_name, false)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, deploymentUrl + "/deinstall_app", data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(appID: number, scIDs: number[]) {
    var data = "[";
    for (var i = 0; i < scIDs.length; i++) {
      if (i != 0) {
        data += ',';
      }
      data += getDataHeader(false);
      data += getDataPair("app_id", appID, true) + ",";
      data += getDataPair("server_class_id", scIDs[i], true) 
      + getDataTail(false);
    }
    data += ']';
    var options = getRequestOption(RequestMethod.Put, deploymentUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(id: number) {
    var options = getRequestOption(RequestMethod.Delete, deploymentUrl + "/" + id, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class UnitPriceService {
  constructor(
    private http:Http
    ) { }

  getRecords(){
    var options = getRequestOption(RequestMethod.Get, unitPriceUrl, "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }

  addRecord(unitPrice: UnitPrice){
    var data = getDataHeader(false)
      + getDataPair("service_price", unitPrice.service_price, true) + "," 
      + getDataPair("storage_price", unitPrice.storage_price, true)
      + getDataTail(false);
    var options = getRequestOption(RequestMethod.Post, unitPriceUrl, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  editRecord(oldUnitPrice, newUnitPrice: UnitPrice) {
    var data = getDataHeader(false)
      + getDataPair("service_price", newUnitPrice.service_price, true) + "," 
      + getDataPair("storage_price", newUnitPrice.storage_price, true)
      + getDataTail(false);

    var options = getRequestOption(RequestMethod.Put, unitPriceUrl + "/" + oldUnitPrice.id, data);
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }

  deleteRecord(id: number) {
    var options = getRequestOption(RequestMethod.Delete, unitPriceUrl + "/" + id, "");
    return this.http.request(new Request(options)).map((res: Response) => {
    });
  }
}

@Injectable()
export class UsageService {
  constructor(
    private http:Http
    ) { }

  getRecords(typeMonthSID: string[]){
    var options = getRequestOption(RequestMethod.Get, usageUrl + "?type="+typeMonthSID[0] + "&month=" + typeMonthSID[1] + "&serviceid=" + typeMonthSID[2], "");
    return this.http.request(new Request(options))
      .map((res:Response) => res.json())
      ;
  }
}