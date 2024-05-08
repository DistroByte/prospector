import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { HomeComponent } from './home/home.component';
import { UserDashboardComponent } from './user-dashboard/user-dashboard.component';
import { UserCreateJobComponent } from './user-create-job/user-create-job.component';
import { authGuard } from './auth.guard';
import { UserProjectPageComponent } from './user-project-page/user-project-page.component';
import { UserGettingStartedComponent } from './user-getting-started/user-getting-started.component';
import { UserDeployRecipeComponent } from './user-deploy-recipe/user-deploy-recipe.component';

export const routes: Routes = [
    { 'path' : '', component: HomeComponent},
    { 'path' : 'login', component: LoginComponent },
    { 'path' : 'user-dashboard', component: UserDashboardComponent, canActivate: [authGuard]},
    { 'path' : 'user-createJob', component: UserCreateJobComponent, canActivate: [authGuard]},
    { 'path' : 'user-gettingStarted', component: UserGettingStartedComponent, canActivate: [authGuard]},
    { 'path' : 'user-deployRecipe', component: UserDeployRecipeComponent, canActivate: [authGuard]},
    { 'path' : 'user-dashboard/:id', component: UserProjectPageComponent, canActivate: [authGuard]}
];
