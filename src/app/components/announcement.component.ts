import { Component, OnInit } from "@angular/core";
import { CommonComponent } from "./common.component";
import { AnnouncementService } from "../services/resource.service";
import { MessageService } from "../services/message";
import { FormControl, Validators, FormGroup, FormBuilder } from "@angular/forms";
import { Announcement, Message } from '../model';
import * as com from "../common";
import * as _ from 'lodash';

@Component({
  selector: 'announcement',
  templateUrl: '../../assets/templates/announcement.component.html',
  providers: [AnnouncementService, MessageService]
})

export class AnnouncementComponent extends CommonComponent{
  admin: string = ""; // used for select ele
  announcements: Announcement[];
  currentAnnouncement: Announcement = new Announcement();
  oldAnnouncement: Announcement = new Announcement();
  orderAttr: string = "created_at";
  constructor(public messageService: MessageService, public fb: FormBuilder, private announcementService: AnnouncementService){
    super(messageService, fb, announcementService, "announcement");
    this.form = fb.group({
      "content": this.contentCtl
    });
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.getRecords();
  }

  getRecords(){
    var callback = (result: any[]) : void => {
      this.announcements = result;
      this.eleActive = [];
      this.eleHovered = [];
      for (var i = 0; i < this.announcements.length; i++){
        this.announcements[i].created_at = com.formatTime(this.announcements[i].created_at);
        this.eleActive.push(new Array<boolean>());
        this.eleHovered.push(new Array<boolean>());
        for(var j = 0; j < 1; j++){
          this.eleActive[i].push(false);
          this.eleHovered[i].push(false);
        }
      }
    }
    super._getRecords(this.announcementService, callback);
  }

  addRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._addRecord(this.announcementService, "announcement", this.currentAnnouncement, callback);
  }
  
  editRecord(r){
    super.submitEdit();
    var callback = (result: any) : void => {
      if(result == "ok") {
        for( var i = 0; i < this.announcements.length; i++){
          if (this.announcements[i].id == this.oldAnnouncement.id)
          {
            this.announcements[i] = _.cloneDeep(this.currentAnnouncement); //if don't clone, when click add, the record editted last time will be set as empty
          }
        }
      }
    }
    super._editRecord(this.announcementService, this.currentAnnouncement.id, this.oldAnnouncement, this.currentAnnouncement, callback);
  }

  deleteRecord(){
    var callback = (result: any) : void => {
      if(result == "ok") {
        this.getRecords();
      }
    }
    super._deleteRecord(this.announcementService, this.currentAnnouncement.id, callback);
  }

  openAddModal(){
    super.openAddModal();
  }

  openDelModal(obj){
    this.currentAnnouncement = _.cloneDeep(obj);
  }

  onEdit(row, col, obj){
    super.onEdit(row, col, obj);
    this.admin = obj.admin? "true": "false";
    this.currentAnnouncement = _.cloneDeep(obj);
    this.oldAnnouncement = _.cloneDeep(obj);
  }
}
