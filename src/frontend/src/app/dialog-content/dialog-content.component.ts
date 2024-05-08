import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { Router } from '@angular/router';

interface DialogData {
  message: string;
  reload: boolean;
}

@Component({
  selector: 'app-dialog-content',
  templateUrl: './dialog-content.component.html',
  styleUrls: ['./dialog-content.component.css'],
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatDialogModule, MatProgressBarModule]
})

export class DialogContentComponent {
  constructor(
    public dialogRef: MatDialogRef<DialogContentComponent>,
    @Inject(MAT_DIALOG_DATA) public data: DialogData,
    private router: Router ) {}

  loading = false;

  onConfirm(): void {
    // this.dialogRef.close(true);
    this.loading = true;
    setTimeout(() => {
      this.dialogRef.close(true);
      this.loading = false;
      // refresh page
      if (this.data.reload === true){
        window.location.reload();
      } else {
        // navigate to user-dashboard
        this.router.navigate(['/user-dashboard']);
      }
    }, 1000);
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}