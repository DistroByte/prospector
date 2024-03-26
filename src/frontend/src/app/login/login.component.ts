import { CommonModule } from "@angular/common";
import { Component } from "@angular/core";
import { Injectable } from "@angular/core";
import axios from "axios";
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatProgressBarModule} from '@angular/material/progress-bar';

@Component({
  selector: "app-login",
  standalone: true,
  imports: [CommonModule, MatProgressSpinnerModule, MatProgressBarModule],
  providers: [Injectable],
  templateUrl: "./login.component.html",
  styleUrls: ["./login.component.css"],
})
export class LoginComponent {
  errorMessage: string | null = null;
  loading = false;

  constructor() {}

  apiUrl = "/api/login";
  headers = {
    "Content-Type": "application/json",
  };

  onSubmit(username: string, password: string, event: Event): void {
    event.preventDefault();
    this.loading = true;
    setTimeout(() => {
      this.loading = false;
      axios
        .post(
          this.apiUrl,
          { username: username, password: password },
          { headers: this.headers }
        )
        .then((response) => {
          console.log(response.data);
          this.errorMessage = null;
          localStorage.setItem('sessionToken', response.data.token);
        })
        .catch((error) => {
          console.error(error);
          this.errorMessage = "Invalid username or password";
        });
    }, 1500);
  }
}
