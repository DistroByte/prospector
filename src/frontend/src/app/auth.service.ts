import { Injectable } from "@angular/core";
import { Router } from "@angular/router";
import axios from "axios";
import { CookieService } from "ngx-cookie-service";
import { environment } from '../environments/environment';

@Injectable({
  providedIn: "root",
})
export class AuthService {
  constructor(private router: Router, private cookieService: CookieService) {}

  // login api url uncomment the below line and comment the next line
  apiUrl = environment.apiUrl;
  headers = {
    "Content-Type": "application/json",
  };

  // logging in implementation
  login(username: string, password: string): void {
    axios
      .post(
        this.apiUrl+`/login`,
        { username: username, password: password },
        { headers: this.headers }
      )
      .then((response) => {
        this.cookieService.set("sessionToken", response.data.token);
      });
  }

  // logging out implementation
  logOut(): void {
    this.cookieService.delete("sessionToken");
    if (this.isLoggedIn() === false) {
      this.router.navigate(["/login"]);
    }
  }

  // checking if user is logged in implementation
  isLoggedIn(): boolean {
    if (this.cookieService.get("sessionToken") === "") {
      return false;
    } else {
      return true;
    }
  }

  // TODO checking if admin is logged in implementation

}
