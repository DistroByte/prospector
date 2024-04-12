import {Component} from '@angular/core';
import {MatButtonModule} from '@angular/material/button';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatIconModule} from '@angular/material/icon';
import { RouterOutlet } from '@angular/router';
import { UserSidebarService } from '../usersidebar.service';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatMenuModule} from '@angular/material/menu';
import { UserSidebarComponent } from '../user-sidebar/user-sidebar.component';
import { AuthService } from '../auth.service';
import { InfoService } from '../info.service';

@Component({
  selector: 'app-user-header',
  standalone: true,
  imports: [MatMenuModule, MatGridListModule, MatToolbarModule,MatSidenavModule, MatIconModule,MatButtonModule,RouterOutlet, UserSidebarComponent],
  templateUrl: './user-header.component.html',
  styleUrl: './user-header.component.css'
})
export class UserHeaderComponent {
  // showFiller = false;
  // public isMenuOpen: boolean = false;
  constructor(public userSidebarService: UserSidebarService, public authService: AuthService, public infoService: InfoService) { }

  userName: string = "";
  ngOnInit(): void {
  }

  toggleSidebar() {
    this.userSidebarService.toggle();
  }

  // logs out the user
  logOut(): void {
    this.authService.logOut();
  }
}
