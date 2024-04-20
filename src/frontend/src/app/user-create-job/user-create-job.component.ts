import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
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

  ngOnInit() {
    this.step = 0;
  }

  onSubmit() {

    const data = {
      components: this.components.map(component => ({
        image: component.Image,
        name: component.Name,
        network: {
          expose: component.Network.Expose,
          port: parseInt(component.Network.Port)
        },
        resources: {
          cpu: parseInt(component.Resources.cpuValue),
          memory: parseInt(component.Resources.ramValue)
        }
      })),
      name: this.projectName,
      type: this.instanceType,
    };

    console.log('Form submitted with data', data);
    console.log(this.components)
    this.InfoService.postJob(data);

    this.formSubmitted = true;
    setTimeout(() => {
      this.router.navigate(['/user-dashboard']);
    }, 2000);
  }

  step = 1;

  setStep(index: number) {
    this.step = index;
  }

  nextStep() {
    this.step++;
  }

  lastStep() {
    this.step;
  }

  componentAdded = false;

  addComponent() {
    console.log(this.step)
    this.componentAdded = true;
    this.components.push({
      Name: '',
      Image: '',
      Network: {
        Port: 0,
        Expose: false,
      },
      Resources: {
        cpuValue: 20,
        ramValue: 20
      }
    });
  }

  removeComponent(index: number) {
    console.log('Removing component');
    this.components.splice(index, 1);
  }

  onToggleChange(event: any, index: number) {
    if (event.checked) {
      // The toggle is checked
      console.log('Toggle is on for component', index);
    } else {
      // The toggle is not checked
      console.log('Toggle is off for component', index);
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
        // Source : https://regex101.com/r/hP8bK1/1
        let dockerRegex = new RegExp("^(?:(?=[^:\/]{4,253})(?!-)[a-zA-Z0-9-]{1,63}(?<!-)(?:\.(?!-)[a-zA-Z0-9-]{1,63}(?<!-))*(?::[0-9]{1,5})?/)?((?![._-])(?:[a-z0-9._-]*)(?<![._-])(?:/(?![._-])[a-z0-9._-]*(?<![._-]))*)(?::(?![.-])[a-zA-Z0-9_.-]{1,128})?$");
        if (!component.Name || component.Resources.cpuValue === 0 || component.Resources.ramValue === 0 || (!dockerRegex.test(component.Image) && component.Image)) {          return false;
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
