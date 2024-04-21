import { Injectable } from '@angular/core';
import axios from 'axios';
import { CookieService } from 'ngx-cookie-service';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class InfoService {

  constructor(private cookieService: CookieService) {}
  
  apiUrl = environment.apiUrl;

  getUser() {
    return axios.get(this.apiUrl+`/v1/user`, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
    .then(response => {
      // console.log(response.data);
      return response.data;
    });
  }

  // placeholder for now to post JOB
  postJob(data: any) {
    return axios.post(this.apiUrl + `/v1/jobs`, data, {
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
      }
    })
      .then(response => {

        return response.data;
      });
  }

  // get all projects created
getAllProjects(){
  return axios.get(this.apiUrl+`/v1/jobs`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  })
  .then(response => {
    // console.log(response.data);
    return response.data;
  });
}

// get all long projects (/v1/jobs?long=true)
getAllLongProjects(){
  return axios.get(this.apiUrl+`/v1/jobs?long=true`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  })
  .then(response => {
    // console.log(response.data);
    return response.data;
  });
}

// get project components
getProjectComponents(projectId: string){
  return axios.get(this.apiUrl+`/v1/jobs/${projectId}/components`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  })
  .then(response => {
    // console.log(response.data);
    return response.data;
  });
}


// return project when given id
getProjectById(id: number): Promise<any> {
  return axios.get(this.apiUrl+`/jobs/${id}`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  })
  .then(response => {
    return response.data;
  });
}

//get Utilization
getCurrentUtilisation(){
  return axios.get(this.apiUrl+`/v1/resources`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    return response.data;
  });
}

getAllocatedResources(){
  return axios.get(this.apiUrl+`/v1/resources/allocated`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    // convert to percentage
    return response.data
  });
}

getProjectAllocatedResources(projectId: string){
  return axios.get(this.apiUrl+`/v1/resources/${projectId}/allocated`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    return response.data;
  });
}

getComponentAllocatedResources(projectId: string, componentId: string){
  return axios.get(this.apiUrl+`/v1/resources/${projectId}/${componentId}/allocated`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    return response.data;
  });
}

getProjectUtilisation(projectId: string){
  return axios.get(this.apiUrl+`/v1/resources/${projectId}`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    return response.data;
  });
}

getComponentUtilisation(projectId: string, componentId: string){
  return axios.get(this.apiUrl+`/v1/resources/${projectId}/${componentId}`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    return response.data;
  });
}

getComponentLogs(projectId: string, componentId: string, type: string){
  return fetch(this.apiUrl+`/v1/jobs/${projectId}/logs?task=${componentId}&type=${type}`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`,
      'Connection': `Upgrade`,
      'Upgrade': `websocket`
    }
  }).then(response => {
    console.log(response);
    return response;
  });
}

getProjectDefinition(projectId: string){
  return axios.get(this.apiUrl+`/v1/jobs/${projectId}/definition`, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    return response.data;
  });
}

updateProjectDefinition(projectId: string, data: any){
  return axios.put(this.apiUrl+`/v1/jobs/${projectId}/definition`, data, {
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${this.cookieService.get("sessionToken")}`
    }
  }).then(response => {
    // console.log(response.data);
    return response.data;
  });
}

}