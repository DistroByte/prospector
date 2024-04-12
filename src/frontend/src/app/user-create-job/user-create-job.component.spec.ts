import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserCreateJobComponent } from './user-create-job.component';

describe('UserCreateJobComponent', () => {
  let component: UserCreateJobComponent;
  let fixture: ComponentFixture<UserCreateJobComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UserCreateJobComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(UserCreateJobComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
