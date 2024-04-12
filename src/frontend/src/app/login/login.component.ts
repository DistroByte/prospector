import { CommonModule } from "@angular/common";
import { Component } from "@angular/core";
import { Injectable } from "@angular/core";
import axios from "axios";
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import { Router } from "@angular/router";
import {
  MatSnackBar,
  MatSnackBarHorizontalPosition,
  MatSnackBarVerticalPosition,
} from '@angular/material/snack-bar';
import { CookieService } from 'ngx-cookie-service';
import { authGuard } from "../auth.guard";
import { AuthService } from "../auth.service";

@Component({
  selector: "app-login",
  standalone: true,
  imports: [CommonModule, MatProgressSpinnerModule, MatProgressBarModule],
  providers: [Injectable],
  templateUrl: "./login.component.html",
  styleUrls: ["./login.component.css"],
})
export class LoginComponent {
  static showPleaseLogin = false;
  errorMessage: string | null = null;
  loading = false;
  Username: string | null = null;

  horizontalPosition: MatSnackBarHorizontalPosition = 'center';
  verticalPosition: MatSnackBarVerticalPosition = 'top';
  
  constructor(private router: Router, private _snackBar: MatSnackBar, private cookieService: CookieService, public authService: AuthService) {}

  ngOnInit(): void{
    if (LoginComponent.showPleaseLogin === true) {
      this._snackBar.open('Unauthorised access, Please login', 'Close', {
        horizontalPosition: this.horizontalPosition,
        verticalPosition: this.verticalPosition,
      });
    }
  }

  onSubmit(username: string, password: string, event: Event): void {
    event.preventDefault();
    this.loading = true;
    this.authService.login(username, password);
    setTimeout(() => {
      this.loading = false;
      if (this.authService.isLoggedIn() === true) {
        this.errorMessage = null;
        this.router.navigate(["/user-dashboard"]);
      } else {
        this.errorMessage = "Invalid Username or Password";
      }
    }, 1000);
  }
}
