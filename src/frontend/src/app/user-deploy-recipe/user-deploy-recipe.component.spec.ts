import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserDeployRecipeComponent } from './user-deploy-recipe.component';

describe('UserDeployRecipeComponent', () => {
  let component: UserDeployRecipeComponent;
  let fixture: ComponentFixture<UserDeployRecipeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UserDeployRecipeComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(UserDeployRecipeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
