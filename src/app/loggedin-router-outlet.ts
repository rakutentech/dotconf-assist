import {Directive, Attribute, ElementRef, Component, ViewContainerRef, ComponentFactoryResolver, ResolvedReflectiveProvider, Inject, Injector, Injectable} from '@angular/core';
import {Router, RouterOutlet, ActivatedRoute, CanActivate, ActivatedRouteSnapshot,
  RouterStateSnapshot} from '@angular/router';

// import {Router, RouterOutlet, RouterOutletMap, ActivatedRoute, CanActivate, ActivatedRouteSnapshot,
//   RouterStateSnapshot} from '@angular/router';
import * as conf from "./configuration";

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private router: Router) {}

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
    let url: string = state.url;
    return this.checkLogin(url);
  }

  checkLogin(url: string): boolean {
    // console.log(localStorage.getItem("mode"));
    if (localStorage.getItem("jwt") == null && localStorage.getItem("mode") == "annotation") {
      this.router.navigate(['/login']);
      return false;
    }
    return true;
    // if (this.authService.isLoggedIn) { return true; }
    // Store the attempted URL for redirecting
    // this.authService.redirectUrl = url;
  }
}