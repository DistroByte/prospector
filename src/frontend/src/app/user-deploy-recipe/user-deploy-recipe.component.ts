import { Component } from '@angular/core';
import { UserHeaderComponent } from '../user-header/user-header.component';
import { UserSidebarComponent } from '../user-sidebar/user-sidebar.component';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { RouterModule } from '@angular/router';
import { InfoService } from '../info.service';

@Component({
  selector: 'app-user-deploy-recipe',
  standalone: true,
  imports: [UserHeaderComponent, UserSidebarComponent, MatCardModule, MatDividerModule, MatButtonModule, MatProgressBarModule, CommonModule, MatIconModule, RouterModule],
  templateUrl: './user-deploy-recipe.component.html',
  styleUrl: './user-deploy-recipe.component.css'
})
export class UserDeployRecipeComponent {

  constructor(private router: Router, private InfoService: InfoService) { }

  isLoading: { [key: string]: boolean } = {};

  isButtonClicked = false;

  submitRecipe(image: string, type: string, name: string) {
    this.isButtonClicked = true;

    const data = {
      components: [{
        image: image,
        name: name,
        network: {
          expose: true,
          port: 8080
        },
        resources: {
          cpu: 50,
          memory: 50
        },
        user_config: {
          ssh_key: "",
        },
        volumes: []
      }],
      name: "Recipe_" + name,
      type: type,
    };
    this.InfoService.postJob(data);
    console.log(data);
    this.isLoading[name] = true;

    setTimeout(() => {
      this.isLoading[image] = false;
      this.router.navigate(['/user-dashboard']);
    }, 2000);
  }

}
