import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserProjectPageComponent } from './user-project-page.component';

describe('UserProjectPageComponent', () => {
  let component: UserProjectPageComponent;
  let fixture: ComponentFixture<UserProjectPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UserProjectPageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(UserProjectPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
