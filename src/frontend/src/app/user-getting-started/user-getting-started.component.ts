import { Component, OnInit } from '@angular/core';
import { UserHeaderComponent } from '../user-header/user-header.component';
import { UserSidebarComponent } from '../user-sidebar/user-sidebar.component';
import {FormBuilder, Validators, FormsModule, ReactiveFormsModule, FormGroup} from '@angular/forms';
import {STEPPER_GLOBAL_OPTIONS} from '@angular/cdk/stepper';
import {MatButtonModule} from '@angular/material/button';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatStepperModule} from '@angular/material/stepper';
import { MatProgressBar } from '@angular/material/progress-bar';
import { CommonModule } from '@angular/common';
import { InfoService } from '../info.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-user-getting-started',
  standalone: true,
  imports: [UserHeaderComponent, UserSidebarComponent, MatStepperModule,
    FormsModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule, MatProgressBar, CommonModule],
  templateUrl: './user-getting-started.component.html',
  styleUrl: './user-getting-started.component.css'
})
export class UserGettingStartedComponent implements OnInit{
  formGroup!: FormGroup;
  loading: boolean = false;
  username: string = '';

  ngOnInit() {
    this.formGroup = this._formBuilder.group({
      projectName: ['', Validators.required],
      componentName: ['', Validators.required],
      componentImage: ['', Validators.required]
    });
    this.getUserName();
  }

  isLinear = false;

  constructor(private _formBuilder: FormBuilder, private InfoService: InfoService, private router: Router) { }

  getUserName() {
    this.InfoService.getUser().then((data) => {
      this.username = data.userName;
    });
  }

  onSubmit() {

    const data = {
      components: [{
        image: this.formGroup.value.componentImage,
        name: this.formGroup.value.componentName,
        network: {
          expose: true,
          port: 80
        },
        resources: {
          cpu: 100,
          memory: 100
        },
        user_config: {
          ssh_key: "",
        },
        volumes: []
      }],
      name: this.formGroup.value.projectName,
      type: "docker",
    };

    this.loading = true;
    this.InfoService.postJob(data);
    setTimeout(() => {
      this.loading = false;
      // send job to backend
      // navigate to dashboard
      this.router.navigate(['/user-dashboard']);
      // open a new tab with the dashboard
      window.open(`https://${this.formGroup.value.componentName}-${this.formGroup.value.projectName}-${this.username}.prospector.ie`, '_blank');
    }, 6000);
  }

}
