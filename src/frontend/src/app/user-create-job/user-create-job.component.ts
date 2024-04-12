import { Component, OnInit } from '@angular/core';
import { MatSliderModule } from '@angular/material/slider';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FooterComponent } from '../footer/footer.component';
import { HeaderComponent } from '../header/header.component';
import { UserHeaderComponent } from '../user-header/user-header.component';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { CommonModule } from '@angular/common';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatButtonModule } from '@angular/material/button';
import { Router, RouterLink, RouterOutlet } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { UserSidebarComponent } from '../user-sidebar/user-sidebar.component';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTooltipModule } from '@angular/material/tooltip';
import { InfoService } from '../info.service';

interface image {
  value: string;
  viewValue: string;
}

@Component({
  selector: 'app-user-create-job',
  standalone: true,
  imports: [MatTooltipModule, MatExpansionModule, MatIconModule, RouterLink, RouterOutlet, MatProgressBarModule, CommonModule, MatSliderModule, MatSlideToggleModule, MatFormFieldModule, MatInputModule, MatSelectModule, FooterComponent, UserHeaderComponent, UserSidebarComponent, MatButtonToggleModule, FormsModule, ReactiveFormsModule, HeaderComponent, MatButtonModule],
  templateUrl: './user-create-job.component.html',
  styleUrl: './user-create-job.component.css'
})

export class UserCreateJobComponent {
  projectName: string;
  instanceType: string;
  formSubmitted: boolean;
  selectedValue: string;

  components: any[] = [];

  images: image[] = [
    { value: 'steak-0', viewValue: 'Ubuntu' },
    { value: 'pizza-1', viewValue: 'Fedora' },
    { value: 'tacos-2', viewValue: 'Debian' },
  ];

  constructor(private InfoService: InfoService, private router: Router) {
    this.projectName = '';
    this.instanceType = '';
    this.formSubmitted = false;
    this.selectedValue = '';
  }

  onSubmit() {

    const data = {
      components: this.components.map(component => ({
        image: component.Image,
        name: this.projectName,
        network: {
          expose: component.Network.Expose,
          mac: "",
          port: parseInt(component.Network.Port)
        },
        resources: {
          cpu: parseInt(component.Resources.cpuValue),
          memory: parseInt(component.Resources.ramValue)
        },
        user_config: {
          ssh_key: "",
          user: ""
        }
      })),
      name: this.projectName,
      type: this.instanceType,
      user: ""
    };

    console.log('Form submitted with data', data);
    this.InfoService.postJob(data);

    this.formSubmitted = true;
    setTimeout(() => {
      this.router.navigate(['/user-dashboard']);
    }, 2000);
  }

  step = 0;

  setStep(index: number) {
    this.step = index;
  }

  nextStep() {
    this.step++;
  }

  prevStep() {
    this.step--;
  }

  componentAdded = false;

  addComponent() {
    this.componentAdded = true;
    console.log('Adding component');
    this.components.push({
      Name: '',
      Image: '',
      Type: '',
      Network: {
        Port: 0,
        Expose: false,
      },
      Resources: {
        cpuValue: 0,
        ramValue: 0
      }
    });
  }

  removeComponent(index: number) {
    console.log('Removing component');
    this.components.splice(index, 1);
  }

  // addPort(component: any) {
  //   component.Network.Ports.push({
  //     Name: '',
  //     Port: 0,
  //   });
  // }

  onToggleChange(event: any, index: number) {
    if (event.checked) {
      // The toggle is checked
      console.log('Toggle is on for component', index);
    } else {
      // The toggle is not checked
      console.log('Toggle is off for component', index);
    }
  }

  removePort(component: any, index: number) {
    component.Network.Ports.splice(index, 1);
    // If there are no ports, then we should disable the Expose toggle AND set it to false
    if (component.Network.Ports.length === 0) {
      component.Network.Expose = false;
    }
  }

  formatCPULabel(value: number) {
    return value + 'hz';
  }

  formatRAMLabel(value: number) {
    return value + 'MB';
  }

  // manual fix for the form validation
  isFormValid() {
    if (!this.projectName || !this.instanceType) {
      return false;
    }

    for (let component of this.components) {
      if (this.instanceType === 'vm') {
        // Validation rules for 'Virtual Machine'
        if (!component.Name || !component.Image || component.Resources.cpuValue === 0 || component.Resources.ramValue === 0) {
          return false;
        }
      } else if (this.instanceType === 'docker') {
        // Validation rules for 'Container'
        if (!component.Name || component.Resources.cpuValue === 0 || component.Resources.ramValue === 0) {
          return false;
        }
      }
    }

    return true;
  }
  // manual fix for when toggling between container or vm the form data is reset
  resetForm() {
    this.componentAdded = false;
    this.components = [];
  }

}
