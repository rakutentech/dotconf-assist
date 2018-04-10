import { Component, OnInit, ViewChild } from "@angular/core";
import { CommonComponent } from "./common.component";
import { InputService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Router, RouterStateSnapshot } from '@angular/router';
import { FileInput, ScriptInput, UnixAppInput, Message } from '../model';
import * as com from "../common";
import * as _ from 'lodash';
import * as conf from "../configuration";
declare var $;


@Component({
  selector: 'log-input',
  templateUrl: '../../assets/templates/log-input.component.html',
  providers: [InputService, MessageService]
})

export class LogInputComponent extends CommonComponent{
  fileInputs: FileInput[] = [];
  scriptInputs: ScriptInput[] = [];
  unixAppInputs: UnixAppInput[] = [];

  currentFileInput: FileInput = new FileInput();
  oldFileInput: FileInput = new FileInput();
  currentScriptInput: ScriptInput = new ScriptInput();
  oldScriptInput: ScriptInput = new ScriptInput();
  currentUnixAppInput: UnixAppInput = new UnixAppInput();
  oldUnixAppInput: UnixAppInput = new UnixAppInput();

  selectedUnixAppInputs: UnixAppInput[] = [];
  selectedUnixAppInputsStr: string[] = [];

  logSizeOptions : string[] = ["< 100 MB","< 500 MB","< 1000 MB", "> 1000 MB"];
  osOptions : string[] = ["*nix","Windows"];
  exeFileOptions : boolean[] = [true, false];
  retentionOptions : string[] = ["1 Day","2 Days","3 Days", "1 Week", "2 Weeks", "3 Weeks", 
  "1 Month", "3 Months", "6 Months", "1 Year", "3 Years", "6 Years"];

  unixAppInputOptions = [
    { script_name: "cpu.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "df.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "vmstat.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "iostat.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "interfaces.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "hardware.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "lastlog.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "lsof.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "netstat.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "netstat2.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "openPorts.sh", interval: "", data_retention_period: "1 Month"},
    { script_name: "openPortsEnhanced.sh", interval: "", data_retention_period: "1 Month"}
  ];
  blacklistTip = [
    { desc: "To ignore and not monitor only files with the .txt extension:", example: "\\.txt$"},
    { desc: "To ignore and not monitor all files with either the .txt extension OR the .gz extension:", example: "\\.(?:txt|gz)$"},
    { desc: "To ignore all files under log file path within the archive or historical directories and all files ending in *.bak:", example: "archive|historical|\\.bak$"},
    { desc: "To ignore files whose names contain a specific string, e.g: ignore the webserver20090228file.txt and webserver20090229file.txt files:", example: "2009022[8|9]file\\.txt$"}
  ];

  env: string = localStorage.getItem("env");
  user: string = localStorage.getItem("user");
  scriptCode: string = "";

  inputType : string = "file";
  unixAppInputChecked: boolean[] = [];
  crcSaltChecked: boolean = false;

  inputsLoaded : boolean[] = [false, false, false];
  
  fileInputForm: FormGroup;
  scriptInputForm: FormGroup;
  unixAppInputForm: FormGroup;
  @ViewChild("scriptFile") fileEle;

  popupLogFilePathStr: string = "You can use wildcard. (e.g. /usr/local/apache/logs/*)";
  popupSourcetypeStr: string = "Please input sourcetype. Setting a meaningful sourcetype is recommended. If you don't have to set specific sourcetype, please input \"auto\"";
  popupBlacklistStr: string = "The blacklist will be ignored by Splunk so that the log will not be sent to Splunk";
  popupIntervalStr: string = "Please specify how often to execute the script (in seconds)(e.g. 60), or a valid cron schedule. (e.g. 30 18 * * *)";
  popupOptionStr: string = "You can specify options for the script if needed. (e.g. -h localhost -u root)";

  secondEleActive: boolean[][] = [];
  secondEleHovered: boolean[][] = [];

  thirdEleActive: boolean[][] = [];
  thirdEleHovered: boolean[][] = [];

  inputLink : string = conf.INPUT_LINK;
  unixappLink : string = conf.UNIXAPP_LINK;

  constructor(public messageService: MessageService, public fb: FormBuilder, 
    private inputService: InputService,
    private router: Router){
    super(messageService, fb, inputService, "input");
    this.fileInputForm = fb.group({
      "log_file_path": this.logFilePathCtl,
      "fileSourcetype": this.fileSourcetypeCtl,
      "log_file_size": this.logFileSizeCtl,
      "data_retention_period": this.dataRetentionCtl,
      "memo": this.memoCtl,
      "blacklist": this.blacklistCtl
    });
    this.scriptInputForm = fb.group({
      "scriptSourcetype": this.scriptSourcetypeCtl,
      "log_file_size": this.logFileSizeCtl,
      "data_retention_period": this.dataRetentionCtl,
      "interval": this.scriptIntervalCtl,
      "option": this.optionCtl,
      "file": this.fileCtl
    });
    this.unixAppInputForm = fb.group({
      "data_retention_period": this.dataRetentionCtl,
      "interval": this.unixAppIntervalCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.inputsLoaded = [false, false, false];
    if (this.env == "stg") {
      this.retentionOptions = ["1 Day","2 Days","3 Days", "1 Week", "2 Weeks", "3 Weeks", 
      "1 Month", "3 Months", "6 Months", "1 Year"];
    }else if (this.env == "dev") {
      this.retentionOptions = ["1 Day","2 Days","3 Days", "1 Week", "2 Weeks", "3 Weeks", 
      "1 Month", "3 Months", "6 Months"];
    }

    for(var i = 0; i < this.unixAppInputOptions.length; i++){
        this.unixAppInputChecked.push(false);
    }
    this.getFileInputRecords();
    this.getScriptInputRecords();
    this.getUnixAppInputRecords();
  }

  getFileInputRecords(){
    this.inputsLoaded[0] = false;
    var callback = (result: any) : void => {
      this.fileInputs = result;
      this.inputsLoaded[0] = true;
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.fileInputs.length; i++){
        this.fileInputs[i].created_at = com.formatTime(this.fileInputs[i].created_at);
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 7; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.inputService, callback, [this.env, this.user, "file", ""]);
  }

  getScriptInputRecords(){
    this.inputsLoaded[1] = false;
    var callback = (result: any) : void => {
      this.scriptInputs = result;
      this.inputsLoaded[1] = true;
      this.secondEleActive = [];
      this.secondEleHovered = [];
      for (var i = 0; i < this.scriptInputs.length; i++){
        this.scriptInputs[i].created_at = com.formatTime(this.scriptInputs[i].created_at);
        this.secondEleActive.push(new Array<boolean>());
        this.secondEleHovered.push(new Array<boolean>());
        for(var j = 0; j < 8; j++){
          this.secondEleActive[i].push(false);
          this.secondEleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.inputService, callback, [this.env, this.user, "script", ""]);
  }

  getUnixAppInputRecords(){
    this.inputsLoaded[2] = false;
    var callback = (result: any) : void => {
      this.unixAppInputs = result;
      this.inputsLoaded[2] = true;
      this.thirdEleActive = [];
      this.thirdEleHovered = [];
      for (var i = 0; i < this.unixAppInputs.length; i++){
        this.unixAppInputs[i].created_at = com.formatTime(this.unixAppInputs[i].created_at);
        this.thirdEleActive.push(new Array<boolean>());
        this.thirdEleHovered.push(new Array<boolean>());
        for(var j = 0; j < 2; j++){
          this.thirdEleActive[i].push(false);
          this.thirdEleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.inputService, callback, [this.env, this.user, "unixapp", ""]);
  }

  addRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        if (this.inputType == "file"){
          this.getFileInputRecords();
        }
        else if (this.inputType == "script") {
          this.getScriptInputRecords();
        }
        else if (this.inputType == "unixapp"){
          this.getUnixAppInputRecords();
        }
      }
    }
    if (this.inputType == "file") {
      this.currentFileInput.env = this.env;
      this.currentFileInput.user_name = this.user;
      this.currentFileInput.app_id = -1;
      super._addRecord(this.inputService, this.currentFileInput.log_file_path, this.currentFileInput, callback, true, "input", [this.env, this.user, this.inputType]);
    }
    else if (this.inputType == "script") {
      this.currentScriptInput.env = this.env;
      this.currentScriptInput.user_name = this.user;
      this.currentScriptInput.app_id = -1;
      super._addRecord(this.inputService, this.currentScriptInput.script_name, this.currentScriptInput, callback, true, "input", [this.env, this.user, this.inputType]);
    }
    else if (this.inputType == "unixapp") {
      for(var i = 0; i < this.selectedUnixAppInputsStr.length; i++){
        var ipt = new UnixAppInput();
        for(var j = 0; j < this.unixAppInputOptions.length; j++){
          if(this.selectedUnixAppInputsStr[i] == this.unixAppInputOptions[j].script_name){
            ipt.interval = this.unixAppInputOptions[j].interval;
            ipt.script_name = this.unixAppInputOptions[j].script_name;
            ipt.data_retention_period = this.unixAppInputOptions[j].data_retention_period;
            ipt.app_id = -1;
            ipt.env = this.env;
            ipt.user_name = this.user;
            this.selectedUnixAppInputs.push(ipt);
            break;
          }
        }
      }
      super._addRecord(this.inputService, this.selectedUnixAppInputs.length + " inputs", this.selectedUnixAppInputs, callback, true, "input", [this.env, this.user, this.inputType]);
    }
  }
  
  editRecord(r){
    this.submitEdit();
    var callback = (result: any) : void => {
      if(result == "ok") {
        if(this.inputType == "file"){
           for( var i = 0; i < this.fileInputs.length; i++){
             if (this.fileInputs[i].id == this.oldFileInput.id)
             {
               this.fileInputs[i] = _.cloneDeep(this.currentFileInput); //if don't clone, when click add, the record editted last time will be set as empty
               return;
             }
           }
        }
        else if (this.inputType == "script"){
          this.getScriptInputRecords(); //for script, and script_code
          // for( var i = 0; i < this.scriptInputs.length; i++){
          //   if (this.scriptInputs[i].id == this.oldScriptInput.id)
          //   {
          //     this.scriptInputs[i] = _.cloneDeep(this.currentScriptInput); //if don't clone, when click add, the record editted last time will be set as empty
          //     return;
          //   }
          // }
        }
        else{
          for( var i = 0; i < this.unixAppInputs.length; i++){
            if (this.unixAppInputs[i].id == this.oldUnixAppInput.id)
            {
              this.unixAppInputs[i] = _.cloneDeep(this.currentUnixAppInput); //if don't clone, when click add, the record editted last time will be set as empty
              return;
            }
          }
        }
      }
    }

    if(this.inputType == "file"){
      super._editRecord(this.inputService, this.currentFileInput.id, this.oldFileInput, this.currentFileInput, callback, true, this.inputType);
    }else if(this.inputType == "script"){
      super._editRecord(this.inputService, this.currentScriptInput.id, this.oldScriptInput, this.currentScriptInput, callback, true, this.inputType);
    }else if(this.inputType == "unixapp"){
      super._editRecord(this.inputService, this.currentUnixAppInput.id, this.oldUnixAppInput, this.currentUnixAppInput, callback, true, this.inputType);
    }
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        if (this.inputType == "file"){
          this.getFileInputRecords();
        }
        else if (this.inputType == "script") {
          this.getScriptInputRecords();
        }
        else if (this.inputType == "unixapp"){
          this.getUnixAppInputRecords();
        }
      }
    }
    if(this.inputType == "file"){
      super._deleteRecord(this.inputService, this.currentFileInput.id, callback, true, this.inputType);
    }else if(this.inputType == "script"){
      super._deleteRecord(this.inputService, this.currentScriptInput.id, callback, true, this.inputType);
    }else if(this.inputType == "unixapp"){
      super._deleteRecord(this.inputService, this.currentUnixAppInput.id, callback, true, this.inputType);
    }
  }

  openAddModal(){
    //file input
    this.fileInputForm.reset();
    this.currentFileInput.log_file_size = this.logSizeOptions[0];
    this.currentFileInput.data_retention_period = this.retentionOptions[0];
    this.crcSaltChecked = false;

    //script input
    this.scriptInputForm.reset();
    this.currentScriptInput.log_file_size = this.logSizeOptions[0];
    this.currentScriptInput.data_retention_period = this.retentionOptions[0];
    this.currentScriptInput.os = this.osOptions[0];
    this.currentScriptInput.exefile = this.exeFileOptions[0];

    //unix app input
    this.unixAppInputForm.reset();
    this.currentUnixAppInput.data_retention_period = this.retentionOptions[0];
    this.selectedUnixAppInputs = [];
    this.resetUnixAppInputOptions();

    //common
    // super.openAddModal();
    this.inputType = "file";
    this.error = "";
    this.isCompleted = true;
    this.modalResizeTimeOut();
  }

  openDelModal(obj){
    if(obj.log_file_path != null){
      this.inputType = "file";
      this.currentFileInput = _.cloneDeep(obj);
    }else if(obj.script != null){
      this.inputType = "script";
      this.currentScriptInput = _.cloneDeep(obj);
    }else{
      this.inputType = "unixapp";
      this.currentUnixAppInput = _.cloneDeep(obj);
    }
  }

  onChangeLogType(type: string){
    this.inputType = type;
    this.modalResizeTimeOut();
  }

  modalResizeTimeOut(){
    setTimeout(() => {
      $(window).trigger('resize');
    }, 800);
  }

  onChangeCrcSalt(input:FileInput, isChecked: boolean){
    this.currentFileInput.crcsalt = isChecked?"<SOURCE>":"";
  }

  onChangeFile(event){
    var fileList = event.target.files;
    if(fileList.length > 0) {
      let file: File = fileList[0];
      // console.log(file);
      // console.log(file.name);
      this.currentScriptInput.script = file;
      this.currentScriptInput.script_name = file.name;
    }else{
      this.currentScriptInput.script = null;
      this.currentScriptInput.script_name = null;
    }
    // var inputEl = this.fileEle.nativeElement;
    // var fileCount = inputEl.files.length;
    // console.log(inputEl);
    // console.log(fileCount);
    // if(fileCount > 0){
    //   this.fileEmpty = false;
    // }
  }

  openShowScriptModal(input: ScriptInput){
    this.currentScriptInput = input;
    this.currentScriptInput.script = null;
    this.modalResizeTimeOut();
  }

  onChangeUnixAppInput(scriptName: string, isChecked: boolean){
    if(isChecked){
      this.selectedUnixAppInputsStr.push(scriptName);
    }else{ //remove from list
      var index = this.selectedUnixAppInputsStr.indexOf(scriptName, 0);
      if (index > -1) {
         this.selectedUnixAppInputsStr.splice(index, 1);
      }
    }
  }

  resetUnixAppInputOptions(){
    this.selectedUnixAppInputsStr = [];
    for(var i = 0; i < this.unixAppInputOptions.length; i++){
      this.unixAppInputOptions[i].interval = "";
      this.unixAppInputOptions[i].data_retention_period = this.retentionOptions[6];
      this.unixAppInputChecked[i] = false;
    }
  }

  recommendedUnixAppInputConfig(){
    this.selectedUnixAppInputsStr = [];
    for( var i = 0; i < 5; i++){
      this.unixAppInputChecked[i] = true;
      this.selectedUnixAppInputsStr.push(this.unixAppInputOptions[i].script_name);
      if(i == 0)
        this.unixAppInputOptions[i].interval = "30";
      else if(i == 1)
        this.unixAppInputOptions[i].interval = "300";
      else
        this.unixAppInputOptions[i].interval = "60";
    }
  }

  resetUnixAppInputConfig(){
    this.resetUnixAppInputOptions();
  }

  onEdit(row, col, obj){
    this.isEditing = true;
    if(obj.log_file_path != null){
      this.inputType = "file";
      this.currentFileInput = _.cloneDeep(obj);
      this.oldFileInput = _.cloneDeep(obj);
      this.crcSaltChecked = this.currentFileInput.crcsalt == "<SOURCE>" ? true: false;
      this.eleActive[row][col] = true;
      this.eleHovered[row][col] = false;
    }
    else if(obj.script != null){
      this.inputType = "script";
      this.currentScriptInput = _.cloneDeep(obj);
      this.currentScriptInput.script = null;
      if(this.fileEle != null)
        this.fileEle.nativeElement.value = "";
      this.oldScriptInput = _.cloneDeep(obj);
      this.secondEleActive[row][col] = true;
      this.secondEleHovered[row][col] = false;
    }
    else{
      this.inputType = "unixapp";
      this.currentUnixAppInput = _.cloneDeep(obj);
      this.oldUnixAppInput = _.cloneDeep(obj);
      this.thirdEleActive[row][col] = true;
      this.thirdEleHovered[row][col] = false;
    }
  }

  submitEdit(){
    this.isEditing = false;
    if (this.eleActive[this.activeCell[0]] != null 
      && this.eleActive[this.activeCell[0]][this.activeCell[1]] != null) {
      this.eleActive[this.activeCell[0]][this.activeCell[1]] = false;
      this.eleHovered[this.activeCell[0]][this.activeCell[1]] = false;
    }
    if (this.secondEleActive[this.activeCell[0]] != null 
      && this.secondEleActive[this.activeCell[0]][this.activeCell[1]] != null) {
      this.secondEleActive[this.activeCell[0]][this.activeCell[1]] = false;
      this.secondEleHovered[this.activeCell[0]][this.activeCell[1]] = false;
    }
    if (this.thirdEleActive[this.activeCell[0]] != null 
      && this.thirdEleActive[this.activeCell[0]][this.activeCell[1]] != null) {
      this.thirdEleActive[this.activeCell[0]][this.activeCell[1]] = false;
      this.thirdEleHovered[this.activeCell[0]][this.activeCell[1]] = false;
    }
    this.activeCell = [-1, -1];
  }

  cancelEditInput(target){
    this.isEditing = false;
    if (target == "second") {
      this.secondEleActive[this.activeCell[0]][this.activeCell[1]] = false;
    }else if(target == "third"){
      this.thirdEleActive[this.activeCell[0]][this.activeCell[1]] = false;
    }
    this.activeCell = [-1, -1];
  }

  mouseOverInput(row, col:number, target){
    if (this.isEditing)
      return;
    this.activeCell = [row, col];
    if (target == "second") {
      this.secondEleHovered[row][col] = true;
      if (this.activeCell[0] >= 0 && this.activeCell[1] >= 0)
        this.secondEleActive[this.activeCell[0]][this.activeCell[1]] = false;
    }else if(target == "third"){
      this.thirdEleHovered[row][col] = true;
      if (this.activeCell[0] >= 0 && this.activeCell[1] >= 0)
        this.thirdEleActive[this.activeCell[0]][this.activeCell[1]] = false;
    }
  }
   
  mouseLeaveInput(target){
    if (this.isEditing) {
      return;
    }
    if (this.activeCell[0] < 0 || this.activeCell[1] < 0) {
      return
    }
    if (target == "second") {
      this.secondEleActive[this.activeCell[0]][this.activeCell[1]] = false;
      this.secondEleHovered[this.activeCell[0]][this.activeCell[1]] = false;
    }else if(target == "third"){
      this.thirdEleActive[this.activeCell[0]][this.activeCell[1]] = false;
      this.thirdEleHovered[this.activeCell[0]][this.activeCell[1]] = false;
    }
  }
}
