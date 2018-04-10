import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { UsageService, UnitPriceService, UserService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { LogSize, StorageSize, UnitPrice, User,  Message } from '../model';
import * as com from "../common";
import * as _ from 'lodash';

@Component({
  selector: 'usage',
  templateUrl: '../../assets/templates/usage.component.html',
  providers: [UsageService, UnitPriceService, UserService, MessageService]
})

export class UsageComponent extends CommonComponent{
  user: User = new User();
  logSizes: LogSize[];
  price: UnitPrice[];
  storageSizes: StorageSize[];
  orderAttr: string = "-size_mb";
  totalLogSize: number = 0;
  totalStorageSize: number = 0;
  monthOptions : string[] = ["","","","","",""];
  month: string = this.monthOptions[0];

  constructor(public messageService: MessageService, 
    public fb: FormBuilder, 
    private usageService: UsageService,
    private unitPriceService: UnitPriceService,
    private userService: UserService){
    super(messageService, fb, usageService, "usage");
    this.form = fb.group({
      "env": this.envCtl
    });
  }

  ngOnInit(): void {
    var toDoubleDigits = function(num) {
      num += "";
      if (num.length === 1) {
        num = "0" + num;
      }
     return num;     
    };

    for(var i = 0; i < 6; i++){
      var dt = new Date();
      if (new Date().getDate() > 16)
        dt.setMonth(dt.getMonth() - (i)); // 0-5 to include current month
      else dt.setMonth(dt.getMonth() - (i+1)); // 1-6
      this.monthOptions[i] = dt.getFullYear().toString() + "-" + toDoubleDigits(dt.getMonth()+1).toString();
    }
    this.month = this.monthOptions[0];
    super.ngOnInit();
    this.getUnitPrice();
    this.getUser();
  }

  getUser(){
    var callback = (result: any) : void => {
      if (result != "error") {
        this.user = result;
      }
    }
    super._getRecord(this.userService, localStorage.getItem("user"), callback, "", "user");
  }

  loadUsage(){
    this.logSizes = [];
    this.storageSizes = [];
    this.getLogSize();
    this.getStorageSize();
  }

  getUnitPrice(){
    var callback = (result: any[]) : void => {
      this.price = result;
    }
    super._getRecords(this.unitPriceService, callback);
  }

  getLogSize(){
    this.totalLogSize = 0;
    var callback = (result: any[]) : void => {
      this.logSizes = result;
      for(var i = 0; i < this.logSizes.length; i++){
        this.totalLogSize += this.logSizes[i].size_mb;
      }
    }
    super._getRecords(this.usageService, callback, ["log_size", this.month, this.user.service_id]);
  }

  getStorageSize(){
    this.totalStorageSize = 0;
    var callback = (result: any[]) : void => {
      this.storageSizes = result;
      for(var i = 0; i < this.storageSizes.length; i++){
        this.totalStorageSize += this.storageSizes[i].size_mb;
      }
    }
    super._getRecords(this.usageService, callback, ["storage", this.month, this.user.service_id]);
  }
}
