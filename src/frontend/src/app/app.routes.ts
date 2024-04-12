import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { HomeComponent } from './home/home.component';
import { UserDashboardComponent } from './user-dashboard/user-dashboard.component';
import { UserCreateJobComponent } from './user-create-job/user-create-job.component';
import { authGuard } from './auth.guard';

export const routes: Routes = [
    { 'path' : '', component: HomeComponent},
    { 'path' : 'login', component: LoginComponent },
    { 'path' : 'user-dashboard', component: UserDashboardComponent, canActivate: [authGuard]},
    { 'path' : 'user-createJob', component: UserCreateJobComponent},
];
