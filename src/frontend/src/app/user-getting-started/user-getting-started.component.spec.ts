import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserGettingStartedComponent } from './user-getting-started.component';

describe('UserGettingStartedComponent', () => {
  let component: UserGettingStartedComponent;
  let fixture: ComponentFixture<UserGettingStartedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UserGettingStartedComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(UserGettingStartedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
