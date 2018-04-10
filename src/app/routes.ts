import { Routes, RouterModule } from "@angular/router";
import { HomeComponent } from "./components/home.component";
import { LoginComponent } from "./components/login.component";
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
import { AuthGuard } from './loggedin-router-outlet';

const routes: Routes = [
    { path: '', redirectTo: '/home', pathMatch: 'full'},
    { component: LoginComponent, path: "login"},
    { component: HomeComponent, path: "home"},
    { component: UserComponent, path: "users", canActivate: [AuthGuard]},
    { component: AnnouncementComponent, path: "announcements", canActivate: [AuthGuard]},
    { component: SplunkHostComponent, path: "splunk_hosts", canActivate: [AuthGuard]},
    { component: ForwarderComponent, path: "forwarders", canActivate: [AuthGuard]},
    { component: ServerClassComponent, path: "server_classes", canActivate: [AuthGuard]},
    { component: LogInputComponent, path: "inputs", canActivate: [AuthGuard]},
    { component: LogAppComponent, path: "apps", canActivate: [AuthGuard]},
    { component: DeploymentComponent, path: "deployment", canActivate: [AuthGuard]},
    { component: UnitPriceComponent, path: "unit_price", canActivate: [AuthGuard]},
    { component: UsageComponent, path: "usage", canActivate: [AuthGuard]},
];

export const routing = RouterModule.forRoot(routes, { useHash: true });
