import { Component } from '@angular/core';

import { Router, RouterLink, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink, RouterOutlet],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css'
})
export class HeaderComponent {

  constructor(private router: Router) {}

  navigateToLogin() {
    console.log('navigateToLogin called');
    this.router.navigate(['login']);
  }

}
