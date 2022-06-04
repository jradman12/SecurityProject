import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { UserServiceService } from 'src/app/services/UserService/user-service.service';

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css']
})
export class LoginPageComponent implements OnInit {

  public form!: FormGroup;
  usernamee!: string;
  constructor( private formBuilder: FormBuilder,
    private _router: Router,
    private _userService: UserServiceService,
    private _snackBar: MatSnackBar) { }

    ngOnInit(): void {
      this.form = this.formBuilder.group({
        username: new FormControl('', [
          Validators.required,
          Validators.pattern('^[a-zA-Z0-9]([._-](?![._-])|[a-zA-Z0-9]){3,18}[a-zA-Z0-9]$'),
        ]), 
        password: new FormControl('', [
          Validators.required,
          Validators.minLength(10),
          Validators.maxLength(30),
          Validators.pattern(
            '^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!"#$@%&()*<>+_|~]).*$'
          )])
      });
    }

    forgotPass() {
      this._userService.sendCode(this.usernamee).subscribe();
      localStorage.setItem('usernamee', this.usernamee);
      this._router.navigate(['/resetPassword']);
  }

      submit():void{
        if (this.form.invalid) return;
        
        const loginObserver = {
          next: (x:any) => {
             this._snackBar.open("Welcome!","",{
              duration : 3000
             });
             
                this._router.navigate(['/user/landing']);
          },
           error: (err:any) => {
             this._snackBar.open("Username or password are incorrect.Try again,please.","",{
               duration : 3000
              }); 

           }};
        
        this._userService.login(this.form.getRawValue()).subscribe(loginObserver);
       }
   
    }


