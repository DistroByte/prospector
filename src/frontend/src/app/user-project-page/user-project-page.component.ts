import { Component, OnInit, input, output } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { UserHeaderComponent } from '../user-header/user-header.component';
import { UserSidebarComponent } from '../user-sidebar/user-sidebar.component';
import { FooterComponent } from '../footer/footer.component';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from "@angular/material/card";
import { MatButtonToggleModule } from "@angular/material/button-toggle";
import { MatFormFieldControl, MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { InfoService } from '../info.service';
import { CommonModule } from '@angular/common';
import { MatChipsModule } from '@angular/material/chips';
import {
  ElementRef,
  OnDestroy,
  ViewChild,
} from "@angular/core";
import { Chart } from 'chart.js';
import { MatDialog } from '@angular/material/dialog';
import { MatIcon } from '@angular/material/icon';
import { StateManagementService } from '../state-management.service';
import { DialogContentComponent } from '../dialog-content/dialog-content.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatSliderModule } from '@angular/material/slider';
import { MatFormField } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

@Component({
  selector: 'app-user-project-page',
  standalone: true,
  imports: [CommonModule, UserHeaderComponent,
    UserSidebarComponent,
    FooterComponent,
    MatTabsModule, MatCardModule, MatButtonToggleModule, MatFormFieldModule, MatSelectModule, MatButtonModule, MatChipsModule, MatIcon
    , FormsModule, ReactiveFormsModule, MatFormFieldModule, MatSliderModule, MatFormField, MatInputModule, MatSlideToggleModule],
  templateUrl: './user-project-page.component.html',
  styleUrl: './user-project-page.component.css'
})

export class UserProjectPageComponent implements OnInit, OnDestroy {
  id: string = '';
  components: any = [];
  selectedComponent: any;
  logs: string[] = [];

  projectDefinition: any = {};
  componentToBeReplaced: any = {};

  constructor(private route: ActivatedRoute, private InfoService: InfoService, private elementRef: ElementRef, public dialog: MatDialog, private router: Router, private StateManagementService: StateManagementService) { }

  ngAfterViewInit() {
  }

  ngOnInit() {
    this.id = this.route.snapshot.paramMap.get('id') ?? '';
    this.getComponents();
    this.getProjectDefintion();
  }

  ngOnDestroy(): void {
  }

  getComponents() {
    this.InfoService.getProjectComponents(this.id).then((data) => {
      this.components = data;
    });
  }

  selectComponent(component: any) {
    this.selectedComponent = component;
    this.findComponentInProjectDefinition();
  }

  getComponentLogs(projectId: string, componentId: string, type: string) {
    this.logs = [];
    this.InfoService.getComponentLogs(projectId, componentId, type).then(async (data) => {
      let stream: any;

      let readStream = async () => {
        const { value, done } = await stream.read();

        if (done) {
          console.log("Stream is done");
          return;
        }

        console.log("Reading stdout");
        const text = new TextDecoder().decode(value);
        const lines = text.split('\n');
        this.logs.push(...lines);
        readStream();
      }

      if (data.status === 200 && data.body) {
        stream = await data.body.getReader();
        readStream();
      }
    });
  }

  restartComponentButton() {
    this.StateManagementService.restartComponent(this.id, this.selectedComponent.id).then((data) => {
    });
  }

  getProjectDefintion() {
    this.InfoService.getProjectDefinition(this.id).then((data) => {
      this.projectDefinition = data;
    });
  }

  // lets find the component that is selected in the project definition
  findComponentInProjectDefinition() {
    for (let i = 0; i < this.projectDefinition.components.length; i++) {
      if (this.projectDefinition.components[i].name === this.selectedComponent.name) {
        this.componentToBeReplaced = this.projectDefinition.components[i];
        // modify the volumes to be a string
        console.log('Component to be replaced', this.componentToBeReplaced);
        this.componentToBeReplaced.volumes = this.setVolumeString(this.componentToBeReplaced.volumes);
        console.log('Component to be replaced', this.componentToBeReplaced);
      }
    }
  }

  onSubmit() {
    this.componentToBeReplaced.volumes = this.componentToBeReplaced.volumes.split(',').map((value: string) => {
      return value.trim();
    });
    this.componentToBeReplaced.resources.cpu = parseInt(this.componentToBeReplaced.resources.cpu);
    this.componentToBeReplaced.resources.memory = parseInt(this.componentToBeReplaced.resources.memory);
    this.componentToBeReplaced.network.port = parseInt(this.componentToBeReplaced.network.port);
    this.InfoService.updateProjectDefinition(this.id, this.projectDefinition).then((data) => {
      console.log('Project definition updated');
    });
  }

  outputTest() {
    console.log(this.componentToBeReplaced);
  }

  formatCPULabel(value: number) {
    return value + 'hz';
  }

  formatRAMLabel(value: number) {
    return value + 'MB';
  }

  setVolumeString(volumes: string[]) {
    return volumes.join(', ');
  }

}