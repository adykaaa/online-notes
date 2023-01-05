import { Axios } from "axios";
import { useState } from "react";


function SignUp() {

  const [emailReg, setEmailReg] = useState("")
  const [usernameReg, setUsernameReg] = useState("")
  const [passwordReg, setPasswordReg] = useState("")

  const register = () => {
    Axios.post(`http://localhost:8080/register` , {email: emailReg, username: usernameReg, password: passwordReg})
  }

  return (
  <div class="login-box">
  <h2>Register</h2>
  <form>
    <div class="user-box">
      <input type="text" name="email" required="" onChange={(e) => setEmailReg(e.target.value)}/>
      <label>E-mail address</label>
    </div>
    <div class="user-box">
      <input type="text" name="username" required="" onChange={(e) => setUsernameReg(e.target.value)}/>  
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" required="" onChange={(e) => setPasswordReg(e.target.value)}/>
      <label>Password</label>
    </div>
    <a>
      <button class="register-button" onClick={register}>Register!</button>
    </a>
  </form>
</div>
  );
}

export default SignUp;
