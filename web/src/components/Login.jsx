import { Link } from "react-router-dom";
import './login.css';
import { useState } from "react";
import axios from "axios";

function Login() {
  
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")

  const login = () => {
    axios.post("http://localhost:8080/login" , { username: username, password: password })
  }

  return (
  <div class="login-box">
  <h2>Login to OnlineNoteZ!</h2>
  <div class="inputs">
    <div class="user-box">
      <input type="text" name="username" onChange={(e) => setUsername(e.target.value)}/>
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" onChange={(e) => setPassword(e.target.value)}/>
      <label>Password</label>
    </div>
    <a>
      <span></span>
      <span></span>
      <span></span>
      <span></span>
      <button class="submit-button" onClick={login}>Submit</button>
    </a>
  </div>
    <div class="signup">
      New here? <Link to="/register">Sign Up!</Link>
    </div>
</div>
  );
}

export default Login;

