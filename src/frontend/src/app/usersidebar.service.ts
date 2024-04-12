import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class UserSidebarService {

  isOpen = false;

  toggle(): void {
    this.isOpen = !this.isOpen;
  }
}