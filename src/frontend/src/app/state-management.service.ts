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

  restartComponent(projectId: string, componentId: string){
    return axios.post(this.apiUrl+`/v1/jobs/${projectId}/component/${componentId}/restart`, {}, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      return response.data;
    });
  }

  stopComponent(projectId: string, componentId: string){
    return axios.post(this.apiUrl+`/v1/jobs/${projectId}/component/${componentId}/stop?purge=false`, {}, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      return response.data;
    });
  }

  deleteComponent(projectId: string, componentId: string){
    return axios.post(this.apiUrl+`/v1/jobs/${projectId}/components/${componentId}/stop?purge=true`, {}, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      return response.data;
    });
  }



}
