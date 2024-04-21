import { Injectable } from '@angular/core';
import axios from 'axios';
import { CookieService } from 'ngx-cookie-service';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class StateManagementService {

  constructor(private cookieService: CookieService) {}
  
  apiUrl = environment.apiUrl;

  startProject(projectId: string){
    return axios.post(this.apiUrl+`/v1/jobs/${projectId}/start`, {}, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      return response.data;
    });
  }

  stopProject(projectId: string){
    return axios.delete(this.apiUrl+`/v1/jobs/${projectId}?purge=false`, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      return response.data;
    });
  }

  deleteProject(projectId: string){
    return axios.delete(this.apiUrl+`/v1/jobs/${projectId}?purge=true`, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      return response.data;
    });
  }

  restartProject(projectId: string){
    return axios.put(this.apiUrl+`/v1/jobs/${projectId}/restart`, {}, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      return response.data;
    });
  }

  restartComponent(projectId: string, componentId: string){
    return axios.put(this.apiUrl+`/v1/jobs/${projectId}/component/${componentId}/restart`, {}, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      console.log(response.data);
      return response.data;
    });
  }

  

}
