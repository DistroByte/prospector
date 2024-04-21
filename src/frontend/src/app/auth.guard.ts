import { CanActivateFn, Router, ActivatedRouteSnapshot, RouterStateSnapshot} from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from './auth.service';

export const authGuard: CanActivateFn = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  const router: Router = inject(Router);
  const authService: AuthService = inject(AuthService);
  const protectedRoutes: string[] = ['/user-dashboard', '/user-createJob', '/user-dashboard/:id'];
  // checks if the user session is unactive and redirects to login page else allows the user to access the page
  if (protectedRoutes.includes(state.url) && authService.isLoggedIn() === false){
    router.navigate(['/login']);
    return false;
  } else {
    return true;
  }
};
