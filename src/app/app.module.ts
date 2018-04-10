import { BrowserModule } from '@angular/platform-browser';
import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { Ng2PaginationModule } from 'ng2-pagination';

import { AppComponent } from './app.component';
import { LoginComponent } from './components/login.component';
import { NgSemanticModule } from 'ng-semantic/ng-semantic';
// import { NgSemanticModule } from "ng-semantic";
import { HomeComponent } from "./components/home.component";
import { UserComponent } from "./components/user.component";
import { AnnouncementComponent } from "./components/announcement.component";
import { SplunkHostComponent } from "./components/splunk-host.component";
import { ForwarderComponent } from "./components/forwarder.component";
import { ServerClassComponent } from "./components/server-class.component";
import { LogInputComponent } from "./components/log-input.component";
import { LogAppComponent } from "./components/log-app.component";
import { DeploymentComponent } from "./components/deployment.component";
import { UnitPriceComponent } from "./components/unit-price.component";
import { UsageComponent } from "./components/usage.component";
import { routing } from "./routes";
import { SearchPipe, SortPipe } from './pipes';
import { AuthGuard } from './loggedin-router-outlet';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    LoginComponent,
    UserComponent,
    AnnouncementComponent,
    SplunkHostComponent,
    ForwarderComponent,
    LogInputComponent,
    LogAppComponent,
    DeploymentComponent,
    ServerClassComponent,
    UnitPriceComponent,
    UsageComponent,
    SearchPipe,
    SortPipe
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    NgSemanticModule,
    Ng2PaginationModule,
    ReactiveFormsModule, //to solve "Can't bind to 'formGroup' since it isn't a known property of 'form'"
    routing
  ],
  providers: [AuthGuard],
  bootstrap: [AppComponent],
  schemas: [ CUSTOM_ELEMENTS_SCHEMA ]
})
export class AppModule { }