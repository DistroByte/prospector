import { TestBed } from '@angular/core/testing';

import { UsersidebarService } from './usersidebar.service';

describe('UsersidebarService', () => {
  let service: UsersidebarService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UsersidebarService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
