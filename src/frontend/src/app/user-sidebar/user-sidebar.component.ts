import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink, RouterOutlet } from '@angular/router';
import { UserSidebarService } from '../usersidebar.service';
import { MatListModule } from '@angular/material/list';
import { AuthService } from '../auth.service';

@Component({
  selector: 'app-user-sidebar',
  standalone: true,
  imports: [MatListModule, MatSidenavModule, MatIconModule, MatButtonModule, RouterOutlet, RouterLink],
  templateUrl: './user-sidebar.component.html',
  styleUrl: './user-sidebar.component.css'
})
export class UserSidebarComponent {
  constructor(public sidebarService: UserSidebarService, private authService: AuthService) { }

  toggleSidebar() {
    this.sidebarService.toggle();
  }

  get isOpen() {
    return this.sidebarService.isOpen;
  }

  logOut(): void {
    this.authService.logOut();
  }

}
