import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { UnitPriceService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { UnitPrice, Message } from '../model';
import * as com from "../common";
import * as _ from 'lodash';

@Component({
  selector: 'unitPrice',
  templateUrl: '../../assets/templates/unit-price.component.html',
  providers: [UnitPriceService, MessageService]
})

export class UnitPriceComponent extends CommonComponent{
  unitPrices: UnitPrice[];
  currentUnitPrice: UnitPrice = new UnitPrice();
  oldUnitPrice: UnitPrice = new UnitPrice();
  constructor(public messageService: MessageService, public fb: FormBuilder, private unitPriceService: UnitPriceService){
    super(messageService, fb, unitPriceService, "unit price");
    this.form = fb.group({
      "service_price": this.serviceCtl,
      "storage_price": this.storageCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.getRecords();
  }

  getRecords(){
    var callback = (result: any[]) : void => {
      this.unitPrices = result;
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.unitPrices.length; i++){
        this.unitPrices[i].created_at = com.formatTime(this.unitPrices[i].created_at);
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 2; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.unitPriceService, callback);
  }

  addRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._addRecord(this.unitPriceService, "service price and storage price", this.currentUnitPrice, callback);
  }
  
  editRecord(r){
    super.submitEdit();
    var callback = (result: any) : void => {
      if(result == "ok") {
        for( var i = 0; i < this.unitPrices.length; i++){
          if (this.unitPrices[i].id == this.oldUnitPrice.id)
          {
            this.unitPrices[i] = _.cloneDeep(this.currentUnitPrice); //if don't clone, when click add, the record editted last time will be set as empty
          }
        }
      }
    }
    super._editRecord(this.unitPriceService, this.currentUnitPrice.id, this.oldUnitPrice, this.currentUnitPrice, callback);
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._deleteRecord(this.unitPriceService, this.currentUnitPrice.id, callback);
  }

  openAddModal(){
    super.openAddModal();
  }

  openDelModal(obj){
    this.currentUnitPrice = _.cloneDeep(obj);
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.currentUnitPrice = _.cloneDeep(obj);
    this.oldUnitPrice = _.cloneDeep(obj);
  }
}
