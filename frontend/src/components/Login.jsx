import { Link } from "react-router-dom";
import './login.css';
import { useState } from "react";
import { Axios } from "axios";

function Login() {
  
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")

  const login = () => {
    Axios.post("http://localhost:8080/login" , { username: username, password: password })
  }

  return (
  <div class="login-box">
  <h2>Login to OnlineNoteZ!</h2>
  <form>
    <div class="user-box">
      <input type="text" name="username" required="" onChange={(e) => setUsername(e.target.value)}/>
      <label>Username</label>
    </div>
    <div class="user-box">
      <input type="password" name="password" required="" onChange={(e) => setPassword(e.target.value)}/>
      <label>Password</label>
    </div>
    <a>
      <span></span>
      <span></span>
      <span></span>
      <span></span>
      <button class="submit-button" onClick={login}>Submit</button>
    </a>
  </form>
    <div class="signup">
      New here? <Link to="/register">Sign Up!</Link>
    </div>
</div>
  );
}

export default Login;

