import {
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild,
} from "@angular/core";
import { UserHeaderComponent } from "../user-header/user-header.component";
import { UserSidebarComponent } from "../user-sidebar/user-sidebar.component";
import { FooterComponent } from "../footer/footer.component";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { Chart } from 'chart.js';
import { color } from 'chart.js/helpers';
import { registerables } from 'chart.js';
import { TreemapController, TreemapElement } from 'chartjs-chart-treemap';
import 'chartjs-chart-treemap';
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatTableDataSource } from '@angular/material/table';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { SelectionModel } from "@angular/cdk/collections";
import { AuthService } from '../auth.service';
import { InfoService } from '../info.service';
import { MatTable, MatTableModule } from '@angular/material/table';
import { CommonModule } from "@angular/common";
import { MatIconModule } from "@angular/material/icon";
import { StateManagementService } from "../state-management.service";
import { DialogContentComponent } from '../dialog-content/dialog-content.component';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { RouterModule } from '@angular/router';

Chart.register(...registerables, TreemapController, TreemapElement);

export interface ProjectConstruct {
  project: string;
  component: string;
  type: string;
  status: string;
  created: string;
  isExpanded?: boolean;
}

const ProjectData: ProjectConstruct[] = [
];

@Component({
  selector: "app-user-dashboard",
  standalone: true,
  imports: [
    MatCheckboxModule,
    MatFormFieldModule,
    MatInputModule,
    MatTableModule,
    MatProgressBarModule,
    MatCardModule,
    MatButtonModule,
    UserHeaderComponent,
    UserSidebarComponent,
    FooterComponent,
    CommonModule,
    MatIconModule,
    MatDialogModule, RouterModule
  ],
  templateUrl: "./user-dashboard.component.html",
  styleUrls: ["./user-dashboard.component.css"],
})

export class UserDashboardComponent implements OnInit, OnDestroy {
  userName: string = "";
  cpuQuota: number = 2000;
  memQuota: number = 2000;
  cpuUtil: number = 0;
  memUtil: number = 0;
  cpuAllocated: number = 0;
  memAllocated: number = 0;
  cpuUtilPercentage: number = 0;
  memUtilPercentage: number = 0;
  cpuQuotaPercentage: number = 0;
  memQuotaPercentage: number = 0;

  lengthOfProjects: number = 0;
  dockerProjects: any[] = [];
  vmProjects: any[] = [];
  runningDocker: number = 0;
  stoppedDocker: number = 0;
  runningVM: number = 0;
  stoppedVM: number = 0;

  recentProjects: any[] = [];

  constructor(private authService: AuthService, private StateManagementService: StateManagementService, private elementRef: ElementRef, private InfoService: InfoService, public dialog: MatDialog, private router: Router) { }

  ngAfterViewInit() {
    this.createUtilChart();
  }

  ngOnInit() {
    this.InfoService.getUser().then(((response: any) => {
      this.userName = response.userName;
    }));

    this.getStats();
    this.getAllocatedResources();
    this.updateRealTimeData();
    this.getRecentProjects();
    this.populateTable();
  }

  ngOnDestroy(): void {
    if (this.updateIntervalId) {
      clearInterval(this.updateIntervalId);
    }
  }

  public chart: any;
  public doughnutchart: any;

  createUtilChart() {
    const ctx = this.elementRef.nativeElement.querySelector('#mixedChart');
    this.chart = new Chart(ctx, {
      data: {
        labels: [],
        datasets: [
          {
            label: 'CPU Utilization (%)',
            data: [],
            type: 'line',
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            borderColor: 'rgb(54, 162, 235)',
            pointRadius: 5,
            pointHoverRadius: 7,
          },
          {
            label: 'Memory Utilization (%)',
            data: [],
            type: 'line',
            backgroundColor: 'rgb(40, 167, 69, 0.2)',
            borderColor: 'rgb(40, 167, 69)',
          },
        ],

      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            beginAtZero: true,
            max: 100,
          },
          x: {
          },
        },
      },
    });
  }

  getStats() {
    this.InfoService.getAllProjects().then((response: any[]) => {
      // gets amount of projects
      this.lengthOfProjects = response.length;

      this.dockerProjects = response.filter(project => project.type === 'docker');
      this.vmProjects = response.filter(project => project.type === 'vm');
      this.runningVM = this.vmProjects.filter(project => project.status === 'running').length;
      this.stoppedVM = this.vmProjects.filter(project => project.status === 'dead').length;

      for (let i = 0; i < this.dockerProjects.length; i++) {
        const projectId = this.dockerProjects[i].id;
        // gets the components of the project
        this.InfoService.getProjectComponents(projectId).then((response: any[]) => {
          // if the project has more than one component
          if (response.length > 1) {
            response.forEach((component: any) => {
              if (component.state === 'running') {
                this.runningDocker++;
              }
              else if (component.state === 'dead') {
                this.stoppedDocker++;
              }
            });
          } else {
            if (this.dockerProjects[i].status === 'running') {
              this.runningDocker++;
            }
            else {
              this.stoppedDocker++;
            }
          }
        });
      }
    });
  }

  private updateIntervalId: any;

  updateRealTimeData() {
    this.getCurrentResourceUtilisation();
    this.updateIntervalId = setInterval(() => {
      this.getCurrentResourceUtilisation();

      const time = new Date().toLocaleTimeString(undefined, { hour12: false });
      if (this.chart.data.labels.length >= 30) {
        this.chart.data.labels.shift();
      }
      this.chart.data.labels.push(time);

      if (this.chart.data.datasets[0].data.length >= 30) {
        this.chart.data.datasets[0].data.shift();
      }
      this.chart.data.datasets[0].data.push(this.cpuUtilPercentage);

      if (this.chart.data.datasets[1].data.length >= 30) {
        this.chart.data.datasets[1].data.shift();
      }
      this.chart.data.datasets[1].data.push(this.memUtilPercentage);

      this.chart.update();
    }, 10000);
  }

  getCurrentResourceUtilisation() {
    this.InfoService.getCurrentUtilisation().then((response: any) => {
      this.cpuUtil = response.cpu;
      this.memUtil = response.memory;
      this.cpuUtilPercentage = (response.cpu / this.cpuAllocated) * 100;
      this.memUtilPercentage = (response.memory / this.memAllocated) * 100;
    });
  }

  getAllocatedResources() {
    this.InfoService.getAllocatedResources().then((response: any) => {
      this.cpuAllocated = response.cpu;
      this.memAllocated = response.memory;
      this.cpuQuotaPercentage = (this.cpuAllocated / this.cpuQuota) * 100;
      this.memQuotaPercentage = (this.memAllocated / this.memQuota) * 100;
    });
  }

  dataSource = new MatTableDataSource<ProjectConstruct>(ProjectData);
  selection = new SelectionModel<ProjectConstruct>(true, []);
  displayedColumns: string[] = ['expand', 'project', 'numberOfComponents', 'type', 'status', 'created', 'select'];
  @ViewChild(MatTable) table!: MatTable<ProjectConstruct>;

  populateTable() {
    this.InfoService.getAllProjects().then((projectsResponse: any) => {
      this.dataSource.data = []; // So we can start fresh - Clear previous data

      projectsResponse.forEach((projectSingle: any) => {
        const submitTime = new Date(projectSingle.SubmitTime / 1000000).toLocaleString('en-US', { hour12: false });

        // gets components for each project
        this.InfoService.getProjectComponents(projectSingle.id).then((componentsResponse: any[]) => {
          this.pushProjectData(
            projectSingle.id,
            componentsResponse,
            projectSingle.type,
            // need to see if this component status or project status are they tied ???
            projectSingle.status,
            projectSingle.created
          );
          this.table.renderRows();
        });
      });
    });
  }

  pushProjectData(project: any, component: any, type: any, status: any, created: any) {
    let info = {
      project: project,
      component: component,
      type: type,
      status: status,
      created: this.convertNanosecondsToDate(created)
    };
    this.dataSource.data.push(info);
  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected === numRows;
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.selection.clear();
      return;
    }

    this.selection.select(...this.dataSource.filteredData);
  }

  /** The label for the checkbox on the passed row */
  checkboxLabel(row?: ProjectConstruct): string {
    if (!row) {
      return `${this.isAllSelected() ? 'deselect' : 'select'} all`;
    }
    return `${this.selection.isSelected(row) ? 'deselect' : 'select'} row ${row.project}`;
  }

  // hack to have only one row expanded at a time
  expandedRow: ProjectConstruct | null = null;

  toggleRow(row: ProjectConstruct) {
    if (this.expandedRow) {
      this.expandedRow.isExpanded = false;
    }
    if (this.expandedRow === row) {
      this.expandedRow = null;
    } else {
      row.isExpanded = true;
      this.expandedRow = row;
      this.getComponentInfo(row.project);
    }
    this.table.renderRows();
  }

  isExpansionDetailRow = (index: number, row: any) => row.isExpanded;

  getRecentProjects() {
    this.InfoService.getAllProjects().then((response: any[]) => {
      if (response.length > 0) {
        this.recentProjects = response.slice(0, 3).map(project => ({
          id: this.extractProjectName(project.id),
          Created: this.convertNanosecondsToDate(project.created)
        }));
      } else {
        this.recentProjects = [];
      }
    });
  }

  extractProjectName(projectId: string): string {
    const parts = projectId.split('-');
    return parts[1];
  }

  convertNanosecondsToDate(nanoseconds: number) {
    const submitTime = new Date(nanoseconds / 1000000);
    return submitTime.toLocaleString('en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    });
  }

  componentStats: any[] = [];

  getComponentInfo(projectId: string) {
    this.InfoService.getProjectComponents(projectId).then((response: any[]) => {
      this.componentStats = response;
    });
  }

  // option buttons
  createProject() {
    // just redirects user to the create project page
    this.router.navigate(['/user-createJob']);
  }

  startProjectButton() {
    const dialogRef = this.dialog.open(DialogContentComponent, {
      width: '250px',
      data: { message: 'Are you sure you want to continue?' }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result === true) {
        console.log('Confirmed');
        for (let i = 0; i < this.selection.selected.length; i++) {
          console.log(this.selection.selected[i].project);
          this.StateManagementService.startProject(this.selection.selected[i].project);
        }
      } else {
        console.log('Canceled');
      }
    });
  }

  stopProjectButton() {
    const dialogRef = this.dialog.open(DialogContentComponent, {
      width: '250px',
      data: { message: 'Are you sure you want to continue?' }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result === true) {
        console.log('Confirmed');
        for (let i = 0; i < this.selection.selected.length; i++) {
          console.log(this.selection.selected[i].project);
          this.StateManagementService.stopProject(this.selection.selected[i].project);
        }
      } else {
        console.log('Canceled');
      }
    });
  }

  restartProjectButton() {
    const dialogRef = this.dialog.open(DialogContentComponent, {
      width: '250px',
      data: { message: 'Are you sure you want to continue?' }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result === true) {
        console.log('Confirmed');
        for (let i = 0; i < this.selection.selected.length; i++) {
          console.log(this.selection.selected[i].project);
          this.StateManagementService.restartProject(this.selection.selected[i].project);
        }
      } else {
        console.log('Cancelled');
      }
    });
  }

  deleteProjectButton() {
    const dialogRef = this.dialog.open(DialogContentComponent, {
      width: '250px',
      data: { message: 'Are you sure you want to continue?' }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result === true) {
        for (let i = 0; i < this.selection.selected.length; i++) {
          console.log(this.selection.selected[i].project);
          this.StateManagementService.deleteProject(this.selection.selected[i].project);
        }
      } else {
        console.log('Cancelled');
      }
    });
  }

  restartComponentButton(projectId: string, ComponentName: string) {
    this.StateManagementService.restartComponent(projectId, ComponentName);
    window.location.reload();
  }
}